package handlers

import (
	"context"
	"strconv"
	"strings"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
)

func (s *Server) getDetailedUsers(
	ctx context.Context,
	userIDs []string,
	guildID string,
) ([]models.DetailedUser, *database.ErrorInfo) {
	if len(userIDs) == 0 {
		return []models.DetailedUser{}, nil
	}

	params := database.GetDetailedUserBasesParams{
		GuildID: guildID,
		UserIds: userIDs,
	}

	bases, ei := database.WrapQuery(
		s.queries.GetDetailedUserBases,
		ctx,
		params,
	)
	if ei != nil {
		return nil, ei
	}

	// Do not query child collections when none of the requested users exist.
	if len(bases) == 0 {
		return []models.DetailedUser{}, nil
	}

	// Use IDs returned by the base query so child queries only process users
	// that actually exist in this guild.
	existingUserIDs := make([]string, len(bases))
	for i, user := range bases {
		existingUserIDs[i] = user.UserID
	}

	rsns, ei := database.WrapQuery(
		s.queries.GetDetailedUserRsns,
		ctx,
		database.GetDetailedUserRsnsParams{
			GuildID: guildID,
			UserIds: existingUserIDs,
		},
	)
	if ei != nil {
		return nil, ei
	}

	records, ei := database.WrapQuery(
		s.queries.GetDetailedUserRecords,
		ctx,
		database.GetDetailedUserRecordsParams{
			GuildID: guildID,
			UserIds: existingUserIDs,
		},
	)
	if ei != nil {
		return nil, ei
	}

	events, ei := database.WrapQuery(
		s.queries.GetDetailedUserEvents,
		ctx,
		database.GetDetailedUserEventsParams{
			GuildID: guildID,
			UserIds: existingUserIDs,
		},
	)
	if ei != nil {
		return nil, ei
	}

	achievements, ei := database.WrapQuery(
		s.queries.GetDetailedUserAchievements,
		ctx,
		database.GetDetailedUserAchievementsParams{
			GuildID: guildID,
			UserIds: existingUserIDs,
		},
	)
	if ei != nil {
		return nil, ei
	}

	combatAchievements, ei := database.WrapQuery(
		s.queries.GetDetailedUserCombatAchievements,
		ctx,
		database.GetDetailedUserCombatAchievementsParams{
			GuildID: guildID,
			UserIds: existingUserIDs,
		},
	)
	if ei != nil {
		return nil, ei
	}

	return models.DetailedUsersFromRows(models.DetailedUserRows{
		Bases:              bases,
		RSNs:               rsns,
		Records:            records,
		Events:             events,
		Achievements:       achievements,
		CombatAchievements: combatAchievements,
	}), nil
}

// Handlers

type GetUsersByIDInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserIDs string `path:"user_ids" doc:"Comma-separated User Snowflake IDs"`
}
type GetUsersByIDOutput struct {
	Body []models.DetailedUser
}

func (s *Server) GetUsersById(ctx context.Context, input *GetUsersByIDInput) (*GetUsersByIDOutput, error) {
	users, ei := s.getDetailedUsers(ctx, strings.Split(input.UserIDs, ","), input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetUsersByIDOutput{Body: users}, nil
}

type GetUsersByRsnInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	RSNs    string `path:"rsns" doc:"Comma-separated RuneScape Names"`
}
type GetUsersByRsnOutput struct {
	Body []models.DetailedUser
}

func (s *Server) GetUsersByRsn(ctx context.Context, input *GetUsersByRsnInput) (*GetUsersByRsnOutput, error) {
	userIDs, ei := database.WrapQuery(s.queries.GetGuildUserByRsn, ctx, database.GetGuildUserByRsnParams{
		GuildID: input.GuildID,
		Rsns:    strings.Split(input.RSNs, ","),
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	users, ei := s.getDetailedUsers(ctx, userIDs, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetUsersByRsnOutput{Body: users}, nil
}

type GetUsersByWomInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	WomIDs  string `path:"wom_ids" doc:"Comma-separated WOM IDs"`
}
type GetUsersByWomOutput struct {
	Body []models.DetailedUser
}

func (s *Server) GetUsersByWom(ctx context.Context, input *GetUsersByWomInput) (*GetUsersByWomOutput, error) {
	userIDs, ei := database.WrapQuery(s.queries.GetGuildUserByWom, ctx, database.GetGuildUserByWomParams{
		GuildID: input.GuildID,
		WomIds:  strings.Split(input.WomIDs, ","),
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	users, ei := s.getDetailedUsers(ctx, userIDs, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetUsersByWomOutput{Body: users}, nil
}

type GetUserAchievementsInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID  string `path:"user_id" doc:"User Snowflake ID"`
}
type GetUserAchievementsOutput struct {
	Body []database.GetUserAchievementsRow
}

func (s *Server) GetUserAchievements(ctx context.Context, input *GetUserAchievementsInput) (*GetUserAchievementsOutput, error) {
	achievements, ei := database.WrapQuery(s.queries.GetUserAchievements, ctx, input.UserID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetUserAchievementsOutput{Body: achievements}, nil
}

type GetUserEventsInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID  string `path:"user_id" doc:"User Snowflake ID"`
}
type GetUserEventsOutput struct {
	Body []database.GetUserEventsRow
}

func (s *Server) GetUserEvents(ctx context.Context, input *GetUserEventsInput) (*GetUserEventsOutput, error) {
	events, ei := database.WrapQuery(s.queries.GetUserEvents, ctx, database.GetUserEventsParams{
		UserID:  input.UserID,
		GuildID: input.GuildID,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetUserEventsOutput{Body: events}, nil
}

type GetUserRecordsInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID  string `path:"user_id" doc:"User Snowflake ID"`
}
type GetUserRecordsOutput struct {
	Body []models.UserRecord
}

func (s *Server) GetUserRecords(ctx context.Context, input *GetUserRecordsInput) (*GetUserRecordsOutput, error) {
	rows, ei := database.WrapQuery(s.queries.GetUserRecords, ctx, database.GetUserRecordsParams{
		UserID:  input.UserID,
		GuildID: input.GuildID,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetUserRecordsOutput{Body: models.UserRecordsFromRows(rows)}, nil
}

type CreateUserInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    models.CreateUserBody
}

type CreateUserOutput struct {
	Body database.CreateUserRow
}

func (s *Server) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	wom, err := s.womClient.GetWom(input.Body.RSN)
	if err != nil {
		return nil, s.womError(err, models.ERROR_WOM_USER_NOT_FOUND)
	}

	params := database.CreateUserParams{
		GuildID: input.GuildID,
		WomID:   strconv.Itoa(wom.Id),
		Rsn:     wom.DisplayName,
		UserID:  string(input.Body.UserID),
	}

	user, err := s.queries.CreateUser(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	return &CreateUserOutput{Body: user}, nil
}

type RemoveUserByIDInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID  string `path:"user_id" doc:"User Snowflake ID"`
}

func (s *Server) RemoveUserById(ctx context.Context, input *RemoveUserByIDInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	// Purge every record (time) the user is part of; cascades to all teammates.
	if _, err := q.DeleteRecordsByUserId(ctx, database.DeleteRecordsByUserIdParams{
		GuildID: input.GuildID,
		UserID:  input.UserID,
	}); err != nil {
		if ei := database.ClassifyError(err); ei != nil {
			return nil, s.dbError(*ei)
		}
	}

	rows, err := q.DeleteUserById(ctx, database.DeleteUserByIdParams{
		GuildID: input.GuildID,
		UserID:  input.UserID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	return nil, nil
}

type RemoveUserByRsnInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	RSN     string `path:"rsn" doc:"RuneScape Name"`
}

func (s *Server) RemoveUserByRsn(ctx context.Context, input *RemoveUserByRsnInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	if _, err := q.DeleteRecordsByRsn(ctx, database.DeleteRecordsByRsnParams{
		GuildID: input.GuildID,
		Rsn:     input.RSN,
	}); err != nil {
		if ei := database.ClassifyError(err); ei != nil {
			return nil, s.dbError(*ei)
		}
	}

	rows, err := q.DeleteUserByRsn(ctx, database.DeleteUserByRsnParams{
		GuildID: input.GuildID,
		Rsn:     input.RSN,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	return nil, nil
}

type RemoveUserByWomInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	WomID   string `path:"wom_id" doc:"WOM ID"`
}

func (s *Server) RemoveUserByWom(ctx context.Context, input *RemoveUserByWomInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	if _, err := q.DeleteRecordsByWom(ctx, database.DeleteRecordsByWomParams{
		GuildID: input.GuildID,
		WomID:   input.WomID,
	}); err != nil {
		if ei := database.ClassifyError(err); ei != nil {
			return nil, s.dbError(*ei)
		}
	}

	rows, err := q.DeleteUserByWom(ctx, database.DeleteUserByWomParams{
		GuildID: input.GuildID,
		WomID:   input.WomID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	return nil, nil
}

type GetBasicUsersOutput struct {
	Body []database.User
}

func (s *Server) GetBasicUsers(ctx context.Context, input *GetUsersByIDInput) (*GetBasicUsersOutput, error) {
	users, ei := database.WrapQuery(
		s.queries.GetUsersById,
		ctx,
		database.GetUsersByIdParams{
			GuildID: input.GuildID,
			UserIds: strings.Split(input.UserIDs, ","),
		},
	)
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	return &GetBasicUsersOutput{Body: users}, nil
}
