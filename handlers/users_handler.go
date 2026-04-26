package handlers

import (
	"context"
	"strconv"
	"strings"

	"tectonic-api/database"
	"tectonic-api/models"
)

// Unchanged private helper
func (s *Server) getDetailedUsers(ctx context.Context, userIDs []string, guildID string) ([]models.DetailedUser, *database.ErrorInfo) {
	if len(userIDs) == 0 {
		return []models.DetailedUser{}, nil
	}

	detailedUsers := make([]models.DetailedUser, 0, len(userIDs))

	for _, userID := range userIDs {
		userRows, err := database.WrapQuery(s.queries.GetUsersById, ctx, database.GetUsersByIdParams{
			GuildID: guildID,
			UserIds: []string{userID},
		})
		if err != nil {
			return nil, err
		}
		if len(userRows) != 1 {
			continue
		}

		user := userRows[0]

		rsnsRows, err := database.WrapQuery(s.queries.GetUserRsns, ctx, database.GetUserRsnsParams{
			UserID: userID, GuildID: guildID,
		})
		if err != nil {
			return nil, err
		}

		recordsRows, err := database.WrapQuery(s.queries.GetUserRecords, ctx, database.GetUserRecordsParams{
			UserID: userID, GuildID: guildID,
		})
		if err != nil {
			return nil, err
		}

		achievementsRows, err := database.WrapQuery(s.queries.GetUserAchievements, ctx, userID)
		if err != nil {
			return nil, err
		}

		eventsRows, err := database.WrapQuery(s.queries.GetUserEvents, ctx, database.GetUserEventsParams{
			UserID: userID, GuildID: guildID,
		})
		if err != nil {
			return nil, err
		}

		caRows, err := database.WrapQuery(s.queries.GetUserCombatAchievements, ctx, database.GetUserCombatAchievementsParams{
			UserID: userID, GuildID: guildID,
		})
		if err != nil {
			return nil, err
		}

		// Get user rank (leaderboard position)
		var userRank int64
		rank, rankErr := database.WrapQuery(s.queries.GetUserRank, ctx, database.GetUserRankParams{
			GuildID: guildID,
			UserID:  userID,
		})
		if rankErr == nil {
			userRank = rank
		}

		// Get user tier (based on points)
		var userTier *models.UserTier
		tier, tierErr := database.WrapQuery(s.queries.GetUserTier, ctx, database.GetUserTierParams{
			GuildID: guildID,
			Points:  user.Points,
		})
		if tierErr == nil {
			t := models.UserTier{
				Name:         tier.Name,
				MinPoints:    tier.MinPoints,
				DisplayOrder: tier.DisplayOrder,
			}
			if tier.Icon.Valid {
				t.Icon = &tier.Icon.String
			}
			if tier.RoleID.Valid {
				t.RoleID = &tier.RoleID.String
			}
			userTier = &t
		}

		detailedUsers = append(detailedUsers, models.DetailedUser{
			UserId:             user.UserID,
			GuildId:            user.GuildID,
			Points:             int(user.Points),
			Rank:               userRank,
			Tier:               userTier,
			RSNs:               models.UserRsnsFromRows(rsnsRows),
			Records:            models.UserRecordsFromRows(recordsRows),
			Events:             models.UserEventFromRows(eventsRows),
			Achievements:       models.UserAchievementsFromRows(achievementsRows),
			CombatAchievements: models.UserCombatAchievementsFromRows(caRows),
		})
	}

	return detailedUsers, nil
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
		return nil, models.NewTectonicError(models.ERROR_RSN_NOT_FOUND)
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
	rows, err := s.queries.DeleteUserById(ctx, database.DeleteUserByIdParams{
		GuildID: input.GuildID,
		UserID:  input.UserID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}
	return nil, nil
}

type RemoveUserByRsnInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	RSN     string `path:"rsn" doc:"RuneScape Name"`
}

func (s *Server) RemoveUserByRsn(ctx context.Context, input *RemoveUserByRsnInput) (*struct{}, error) {
	rows, err := s.queries.DeleteUserByRsn(ctx, database.DeleteUserByRsnParams{
		GuildID: input.GuildID,
		Rsn:     input.RSN,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}
	return nil, nil
}

type RemoveUserByWomInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	WomID   string `path:"wom_id" doc:"WOM ID"`
}

func (s *Server) RemoveUserByWom(ctx context.Context, input *RemoveUserByWomInput) (*struct{}, error) {
	rows, err := s.queries.DeleteUserByWom(ctx, database.DeleteUserByWomParams{
		GuildID: input.GuildID,
		WomID:   input.WomID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}
	return nil, nil
}
