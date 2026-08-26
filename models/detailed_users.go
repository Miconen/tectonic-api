package models

import "tectonic-api/database"

type DetailedUserRows struct {
	Bases              []database.GetDetailedUserBasesRow
	RSNs               []database.GetDetailedUserRsnsRow
	Records            []database.GetDetailedUserRecordsRow
	Events             []database.GetDetailedUserEventsRow
	Achievements       []database.GetDetailedUserAchievementsRow
	CombatAchievements []database.GetDetailedUserCombatAchievementsRow
}

func DetailedUsersFromRows(rows DetailedUserRows) []DetailedUser {
	users := make([]DetailedUser, 0, len(rows.Bases))
	userIndexes := make(map[string]int, len(rows.Bases))

	for _, row := range rows.Bases {
		user := DetailedUser{
			UserId:             row.UserID,
			GuildId:            row.GuildID,
			Points:             int(row.Points),
			Rank:               row.UserRank,
			RSNs:               []UserRsn{},
			Records:            []UserRecord{},
			Events:             []UserEvent{},
			Achievements:       []UserAchievement{},
			CombatAchievements: []UserCombatAchievement{},
		}

		user.Tier = &UserTier{
			Name:         row.TierName,
			Icon:         row.TierIcon,
			RoleID:       row.TierRoleID,
			MinPoints:    row.TierMinPoints,
			DisplayOrder: row.TierDisplayOrder,
		}

		userIndexes[row.UserID] = len(users)
		users = append(users, user)
	}

	for _, row := range rows.RSNs {
		index, ok := userIndexes[row.UserID]
		if !ok {
			continue
		}

		users[index].RSNs = append(users[index].RSNs, UserRsn{
			RSN:   row.Rsn,
			WomId: row.WomID,
		})
	}

	for _, row := range rows.Achievements {
		index, ok := userIndexes[row.UserID]
		if !ok {
			continue
		}

		users[index].Achievements = append(
			users[index].Achievements,
			UserAchievement{
				Name:        row.Name,
				Thumbnail:   row.Thumbnail,
				DiscordIcon: row.DiscordIcon,
				Order:       row.AchievementOrder,
			},
		)
	}

	for _, row := range rows.CombatAchievements {
		index, ok := userIndexes[row.UserID]
		if !ok {
			continue
		}

		users[index].CombatAchievements = append(
			users[index].CombatAchievements,
			UserCombatAchievement{Name: row.Name},
		)
	}

	for _, row := range rows.Events {
		index, ok := userIndexes[row.UserID]
		if !ok {
			continue
		}

		users[index].Events = append(users[index].Events, UserEvent{
			Name:           row.Name,
			WomID:          row.EventID,
			GuildID:        row.GuildID,
			Placement:      row.Placement,
			PositionCutoff: row.PositionCutoff,
			Solo:           row.Solo,
		})
	}

	// Record ID is only unique within a guild, but every result here belongs
	// to the same requested guild. Keep indexes per owner because the same
	// record can appear on multiple requested users.
	recordIndexes := make(map[string]map[int32]int)

	for _, row := range rows.Records {
		userIndex, ok := userIndexes[row.OwnerUserID]
		if !ok {
			continue
		}

		ownerRecords, ok := recordIndexes[row.OwnerUserID]
		if !ok {
			ownerRecords = make(map[int32]int)
			recordIndexes[row.OwnerUserID] = ownerRecords
		}

		recordIndex, exists := ownerRecords[row.RecordID]
		if !exists {
			recordIndex = len(users[userIndex].Records)
			ownerRecords[row.RecordID] = recordIndex

			users[userIndex].Records = append(
				users[userIndex].Records,
				UserRecord{
					Id:          row.RecordID,
					BossName:    row.BossName,
					DisplayName: row.DisplayName,
					Category:    row.Category,
					Solo:        row.Solo,
					ValueType:   row.ValueType,
					Date:        row.Date.Time,
					Value:       row.Value,
					Teammates:   []RecordTeammate{},
				},
			)
		}

		users[userIndex].Records[recordIndex].Teammates = append(
			users[userIndex].Records[recordIndex].Teammates,
			RecordTeammate{
				UserID:  row.TeammateUserID,
				GuildID: row.TeammateGuildID,
			},
		)
	}

	return users
}
