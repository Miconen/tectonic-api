package handlers

import (
	"context"
	"sort"
	"time"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetGuildRecordsInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetGuildRecordsOutput struct {
	Body models.GuildResponse
}

func (s *Server) GetGuildRecords(ctx context.Context, input *GetGuildRecordsInput) (*GetGuildRecordsOutput, error) {
	row, err := s.queries.GetDetailedGuild(ctx, input.GuildID)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	guild := models.GuildResponseFromDetailedRow(row)
	return &GetGuildRecordsOutput{Body: guild}, nil
}

// Teams

type AddTeammateByBossInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Boss    string `path:"boss" doc:"Boss name"`
	Body    models.InputTeammate
}

func (s *Server) AddTeammateByBoss(ctx context.Context, input *AddTeammateByBossInput) (*struct{}, error) {
	params := database.AddToTeamByBossParams{
		GuildID:  string(input.Body.GuildID),
		BossName: input.Boss,
		UserID:   string(input.Body.UserID),
	}
	ei := database.WrapExec(s.queries.AddToTeamByBoss, ctx, params)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

type AddTeammateByRecordIdInput struct {
	GuildID  string `path:"guild_id" doc:"Guild Snowflake ID"`
	RecordID int    `path:"record_id" doc:"Record ID"`
	Body     models.InputTeammate
}

func (s *Server) AddTeammateByRecordId(ctx context.Context, input *AddTeammateByRecordIdInput) (*struct{}, error) {
	params := database.AddToTeamByIdParams{
		GuildID:  string(input.Body.GuildID),
		UserID:   string(input.Body.UserID),
		RecordID: int32(input.RecordID),
	}
	ei := database.WrapExec(s.queries.AddToTeamById, ctx, params)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

type RemoveTeammateByBossInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Boss    string `path:"boss" doc:"Boss name"`
	Body    models.InputTeammate
}

func (s *Server) RemoveTeammateByBoss(ctx context.Context, input *RemoveTeammateByBossInput) (*struct{}, error) {
	params := database.RemoveFromTeamByBossParams{
		GuildID:  string(input.Body.GuildID),
		UserID:   string(input.Body.UserID),
		BossName: input.Boss,
	}
	rows, ei := database.WrapQuery(s.queries.RemoveFromTeamByBoss, ctx, params)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_TEAM_NOT_FOUND)
	}
	return nil, nil
}

type RemoveTeammateByRecordIdInput struct {
	GuildID  string `path:"guild_id" doc:"Guild Snowflake ID"`
	RecordID int    `path:"record_id" doc:"Record ID"`
	Body     models.InputTeammate
}

func (s *Server) RemoveTeammateByRecordId(ctx context.Context, input *RemoveTeammateByRecordIdInput) (*struct{}, error) {
	params := database.RemoveFromTeamByIdParams{
		GuildID:  string(input.Body.GuildID),
		UserID:   string(input.Body.UserID),
		RecordID: int32(input.RecordID),
	}
	rows, ei := database.WrapQuery(s.queries.RemoveFromTeamById, ctx, params)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_TEAM_NOT_FOUND)
	}
	return nil, nil
}

// Records

type CreateRecordInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    models.InputRecord
}
type CreateRecordOutput struct {
	Body models.RecordResponse
}

func (s *Server) CreateRecord(ctx context.Context, input *CreateRecordInput) (*CreateRecordOutput, error) {
	res := models.RecordResponse{
		BossName: input.Body.BossName,
		Value:    input.Body.Value,
	}

	// Verify boss exists and get its info
	bossInfo, err := s.queries.GetBossInfo(ctx, input.Body.BossName)
	if ei := database.ClassifyError(err); ei != nil {
		if ei.Recoverable && ei.Code == "P0002" {
			return nil, models.NewTectonicError(models.ERROR_GUILD_BOSS_NOT_FOUND)
		}
		return nil, s.dbError(*ei)
	}

	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	// Always insert the record
	recordID, err := q.CreateRecord(ctx, database.CreateRecordParams{
		Value:    int32(input.Body.Value),
		BossName: input.Body.BossName,
		Date:     pgtype.Timestamp{Time: time.Now(), Valid: true},
		GuildID:  input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	// Create team entries
	err = q.CreateTeam(ctx, database.CreateTeamParams{
		RecordID: recordID,
		UserIds:  models.SnowflakesToStrings(input.Body.UserIDs),
		GuildID:  input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	// Get all records for this boss to determine position
	allRecords, err := q.GetBossRecords(ctx, database.GetBossRecordsParams{
		GuildID:  input.GuildID,
		BossName: input.Body.BossName,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}

	// Compute position in Go
	position := computeRecordPosition(recordID, allRecords, bossInfo)
	res.RecordID = int(recordID)
	if position > 0 {
		res.Position = &position
	}

	// Find the previous value at this position (if any) for the response
	oldValue := findOldValueAtPosition(recordID, allRecords, bossInfo, position)
	res.OldValue = oldValue

	return &CreateRecordOutput{Body: res}, nil
}

// computeRecordPosition determines where a record ranks among all records for a boss.
// For solo bosses, deduplicates by user (best per user).
// Returns the 1-based position, or 0 if the record is not in the ranking
// (e.g., the user has a better record for a solo boss).
func computeRecordPosition(recordID int32, rows []database.GetBossRecordsRow, bossInfo database.GetBossInfoRow) int {
	type recordEntry struct {
		RecordID int32
		Value    int32
		UserID   string // first user (for solo dedup)
	}

	// Group rows by record_id
	recordMap := make(map[int32]*recordEntry)
	var recordOrder []int32
	for _, row := range rows {
		if _, ok := recordMap[row.RecordID]; !ok {
			recordMap[row.RecordID] = &recordEntry{
				RecordID: row.RecordID,
				Value:    row.Value,
				UserID:   row.UserID,
			}
			recordOrder = append(recordOrder, row.RecordID)
		}
	}

	var eligible []recordEntry

	if bossInfo.Solo {
		// For solo bosses: pick each user's best record
		bestByUser := make(map[string]*recordEntry)
		for _, rid := range recordOrder {
			entry := recordMap[rid]
			existing, ok := bestByUser[entry.UserID]
			if !ok || isBetter(entry.Value, existing.Value, bossInfo.HigherIsBetter) {
				bestByUser[entry.UserID] = entry
			}
		}
		for _, entry := range bestByUser {
			eligible = append(eligible, *entry)
		}
	} else {
		// For team bosses: all records are eligible
		for _, rid := range recordOrder {
			eligible = append(eligible, *recordMap[rid])
		}
	}

	// Sort by value (respecting higher_is_better)
	sort.Slice(eligible, func(i, j int) bool {
		return isBetter(eligible[i].Value, eligible[j].Value, bossInfo.HigherIsBetter)
	})

	// Find position of the target record
	for i, entry := range eligible {
		if entry.RecordID == recordID {
			return i + 1 // 1-based
		}
	}
	return 0
}

// findOldValueAtPosition finds what value was previously at the given position
// before the new record was inserted.
func findOldValueAtPosition(newRecordID int32, rows []database.GetBossRecordsRow, bossInfo database.GetBossInfoRow, position int) int {
	if position <= 0 {
		return 0
	}

	type recordEntry struct {
		RecordID int32
		Value    int32
		UserID   string
	}

	// Build entries excluding the new record
	recordMap := make(map[int32]*recordEntry)
	var recordOrder []int32
	for _, row := range rows {
		if row.RecordID == newRecordID {
			continue
		}
		if _, ok := recordMap[row.RecordID]; !ok {
			recordMap[row.RecordID] = &recordEntry{
				RecordID: row.RecordID,
				Value:    row.Value,
				UserID:   row.UserID,
			}
			recordOrder = append(recordOrder, row.RecordID)
		}
	}

	var eligible []recordEntry
	if bossInfo.Solo {
		bestByUser := make(map[string]*recordEntry)
		for _, rid := range recordOrder {
			entry := recordMap[rid]
			existing, ok := bestByUser[entry.UserID]
			if !ok || isBetter(entry.Value, existing.Value, bossInfo.HigherIsBetter) {
				bestByUser[entry.UserID] = entry
			}
		}
		for _, entry := range bestByUser {
			eligible = append(eligible, *entry)
		}
	} else {
		for _, rid := range recordOrder {
			eligible = append(eligible, *recordMap[rid])
		}
	}

	sort.Slice(eligible, func(i, j int) bool {
		return isBetter(eligible[i].Value, eligible[j].Value, bossInfo.HigherIsBetter)
	})

	// Return the value that was at this position before
	idx := position - 1
	if idx < len(eligible) {
		return int(eligible[idx].Value)
	}
	return 0
}

func isBetter(a, b int32, higherIsBetter bool) bool {
	if higherIsBetter {
		return a > b
	}
	return a < b
}

type RemoveRecordInput struct {
	GuildID  string `path:"guild_id" doc:"Guild Snowflake ID"`
	RecordID int    `path:"record_id" doc:"Record ID"`
}

func (s *Server) RemoveRecord(ctx context.Context, input *RemoveRecordInput) (*struct{}, error) {
	params := database.DeleteRecordParams{
		GuildID:  input.GuildID,
		RecordID: int32(input.RecordID),
	}
	deleted, err := s.queries.DeleteRecord(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if deleted == 0 {
		return nil, models.NewTectonicError(models.ERROR_RECORD_NOT_FOUND)
	}
	return nil, nil
}

type ClearBossRecordsInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Boss    string `path:"boss" doc:"Boss name"`
}

func (s *Server) ClearBossRecords(ctx context.Context, input *ClearBossRecordsInput) (*struct{}, error) {
	removed, err := s.queries.DeleteBossRecords(ctx, database.DeleteBossRecordsParams{
		GuildID:  input.GuildID,
		BossName: input.Boss,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if removed == 0 {
		return nil, models.NewTectonicError(models.ERROR_RECORD_NOT_FOUND)
	}
	return nil, nil
}

type RevertTopRecordInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Boss    string `path:"boss" doc:"Boss name"`
}

func (s *Server) RevertTopRecord(ctx context.Context, input *RevertTopRecordInput) (*struct{}, error) {
	removed, err := s.queries.DeleteTopRecord(ctx, database.DeleteTopRecordParams{
		GuildID:  input.GuildID,
		BossName: input.Boss,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if removed == 0 {
		return nil, models.NewTectonicError(models.ERROR_RECORD_NOT_FOUND)
	}
	return nil, nil
}
