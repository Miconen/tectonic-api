package database

// TODO: maybe the query is too convoluted to normalize, so might be better
// to search alternatives.
func NormalizeDetailedUsersRow(rows []GetDetailedUsersRow) []DetailedUser {
	users := make([]DetailedUser, 0)

	var (
		currentUser     *DetailedUser = nil
		currentTime     *DetailedTime = nil
		currentTeammate *UserData     = nil
	)

	for _, row := range rows {
		user := newDetailedUser(row)
		time := newDetailedTime(row)
		teammate := newTeammate(row)

		teammate_rsn := Rsn{
			Rsn:     row.TeamUserRsn.String,
			WomID:   row.TeamUserWomID.String,
			UserID:  row.TeamUserID.String,
			GuildID: row.TeamGuildID.String,
		}

		if currentUser != nil {
			currentUser = user
		} else if row.UserID != currentUser.User.UserID {
			newUser := newDetailedUser(row)
			users = append(users, *currentUser)
			currentUser = newUser
		}

		if currentTime != nil {
			currentTime = time
		} else if row.RunID.Int32 != currentTime.Time.RunID {
			currentUser.Times = append(currentUser.Times, *currentTime)
			currentTime = time
		}

		if currentTeammate != nil {
			currentTeammate = newTeammate(row)
		} else if row.TeamUserID.String != currentTeammate.UserID {
			currentTime.Teammates = append(currentTime.Teammates, *currentTeammate)
			currentTeammate = teammate
		}
		currentTeammate.RSNs = append(currentTeammate.RSNs, teammate_rsn)
	}

	return users
}

func newDetailedUser(row GetDetailedUsersRow) *DetailedUser {
	return &DetailedUser{
		User: UserData{
			UserID:  row.UserID,
			GuildID: row.UserGuildID,
			Points:  row.UserPoints,
			RSNs:    make([]Rsn, 0),
		},
		Times: make([]DetailedTime, 0),
	}
}

func newDetailedTime(row GetDetailedUsersRow) *DetailedTime {
	return &DetailedTime{
		Time: Time{
			Time:     row.TimeValue.Int32,
			BossName: row.BossName.String,
			RunID:    row.RunID.Int32,
			Date:     row.RunDate,
		},
		Teammates: make([]UserData, 0),
	}
}

func newTeammate(row GetDetailedUsersRow) *UserData {
	return &UserData{
		UserID:  row.TeamUserID.String,
		GuildID: row.TeamGuildID.String,
		Points:  row.TeamUserPoints.Int32,
		RSNs:    make([]Rsn, 0),
	}
}

