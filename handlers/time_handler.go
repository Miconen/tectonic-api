package handlers

import (
	"context"
	"fmt"
	"time"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetGuildTimesInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetGuildTimesOutput struct {
	Body models.Guild
}

func (s *Server) GetGuildTimes(ctx context.Context, input *GetGuildTimesInput) (*GetGuildTimesOutput, error) {
	row, err := s.queries.GetDetailedGuild(ctx, input.GuildID)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	guild := models.GuildResponseFromDetailedRow(row)
	return &GetGuildTimesOutput{Body: guild}, nil
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

type AddTeammateByRunIdInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	RunID   int    `path:"run_id" doc:"Run ID"`
	Body    models.InputTeammate
}

func (s *Server) AddTeammateByRunId(ctx context.Context, input *AddTeammateByRunIdInput) (*struct{}, error) {
	params := database.AddToTeamByIdParams{
		GuildID: string(input.Body.GuildID),
		UserID:  string(input.Body.UserID),
		RunID:   int32(input.RunID),
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

type RemoveTeammateByRunIdInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	RunID   int    `path:"run_id" doc:"Run ID"`
	Body    models.InputTeammate
}

func (s *Server) RemoveTeammateByRunId(ctx context.Context, input *RemoveTeammateByRunIdInput) (*struct{}, error) {
	params := database.RemoveFromTeamByIdParams{
		GuildID: string(input.Body.GuildID),
		UserID:  string(input.Body.UserID),
		RunID:   int32(input.RunID),
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

// Times

type CreateTimeInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    models.InputTime
}
type CreateTimeOutput struct {
	Body models.TimeResponse
}

func (s *Server) CreateTime(ctx context.Context, input *CreateTimeInput) (*CreateTimeOutput, error) {
	res := models.TimeResponse{
		BossName: input.Body.BossName,
		Time:     input.Body.Time,
	}

	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	pb, err := q.CheckPb(ctx, database.CheckPbParams{
		Boss:    input.Body.BossName,
		GuildID: input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		if ei.Recoverable && ei.Code == "P0002" {
			return nil, models.NewTectonicError(models.ERROR_GUILD_BOSS_NOT_FOUND)
		}
		return nil, s.dbError(*ei)
	}

	if pb.Time.Valid {
		if int(pb.Time.Int32) <= input.Body.Time {
			return &CreateTimeOutput{Body: res}, nil
		}
	}

	runID, err := q.CreateTime(ctx, database.CreateTimeParams{
		Time:     int32(input.Body.Time),
		BossName: input.Body.BossName,
		Date:     pgtype.Timestamp{Time: time.Now(), Valid: true},
		GuildID:  input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	_, err = q.UpdatePb(ctx, database.UpdatePbParams{
		RunID:   pgtype.Int4{Int32: runID, Valid: true},
		GuildID: input.GuildID,
		Boss:    input.Body.BossName,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	err = q.CreateTeam(ctx, database.CreateTeamParams{
		RunID:   runID,
		UserIds: models.SnowflakesToStrings(input.Body.UserIDs),
		GuildID: input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}

	res.RunID = int(runID)
	res.OldTime = int(pb.Time.Int32)
	return &CreateTimeOutput{Body: res}, nil
}

type RemoveTimeInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	TimeID  int    `path:"time_id" doc:"Time/Run ID"`
}

func (s *Server) RemoveTime(ctx context.Context, input *RemoveTimeInput) (*struct{}, error) {
	params := database.DeleteTimeParams{
		GuildID: input.GuildID,
		RunID:   int32(input.TimeID),
	}
	deleted, err := s.queries.DeleteTime(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if deleted == 0 {
		return nil, models.NewTectonicError(models.ERROR_TIME_NOT_FOUND)
	}
	return nil, nil
}

type ClearClanPbInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Boss    string `path:"boss" doc:"Boss name"`
}

func (s *Server) ClearClanPb(ctx context.Context, input *ClearClanPbInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	removed, ei := s.deletePbRecord(ctx, q, input.GuildID, input.Boss)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	if removed == 0 {
		return nil, models.NewTectonicError(models.ERROR_TIME_NOT_FOUND)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	return nil, nil
}

type RevertClanPbInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Boss    string `path:"boss" doc:"Boss name"`
}

func (s *Server) RevertClanPb(ctx context.Context, input *RevertClanPbInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	removed, ei := s.deletePbRecord(ctx, q, input.GuildID, input.Boss)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	if removed == 0 {
		return nil, models.NewTectonicError(models.ERROR_TIME_NOT_FOUND)
	}

	updateParams := database.RevertGuildBossPbParams{
		GuildID:  input.GuildID,
		BossName: input.Boss,
	}
	fmt.Println(updateParams)
	revertEi := database.WrapExec(q.RevertGuildBossPb, ctx, updateParams)
	if revertEi != nil {
		return nil, s.dbError(*revertEi)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	return nil, nil
}

// Unchanged private helper
func (s *Server) deletePbRecord(ctx context.Context, q *database.Queries, guildID, bossName string) (int64, *database.ErrorInfo) {
	removed, err := q.DeletePb(ctx, database.DeletePbParams{
		GuildID:  guildID,
		BossName: bossName,
	})
	if err != nil {
		if ei := database.ClassifyError(err); ei != nil {
			return 0, ei
		}
	}
	return removed, nil
}
