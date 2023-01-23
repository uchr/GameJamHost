package storages

import (
	"context"

	"GameJamPlatform/internal/models/gamejams"
)

func (st *storage) AddVote(ctx context.Context, vote gamejams.Vote) error {
	_, err := st.db.Exec(ctx, "INSERT INTO votes (game_id, user_id, criteria_id, value) VALUES ($1, $2, $3, $4)", vote.GameID, vote.UserID, vote.CriteriaID, vote.Value)
	if err != nil {
		return err
	}

	return nil
}

func (st *storage) GetVote(ctx context.Context, criteriaID int) ([]gamejams.Vote, error) {
	rows, err := st.db.Query(ctx, "SELECT game_id, user_id, criteria_id, value FROM votes WHERE criteria_id = $1", criteriaID)
	if err != nil {
		return nil, err
	}

	var votes []gamejams.Vote

	for rows.Next() {
		var vote gamejams.Vote

		err = rows.Scan(&vote.GameID, &vote.UserID, &vote.CriteriaID, &vote.Value)
		if err != nil {
			return nil, err
		}

		votes = append(votes, vote)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return votes, nil
}
