package storages

import (
	"context"

	"GameJamPlatform/internal/gamejam"
)

func (st *storage) CreateJam(ctx context.Context, jam gamejam.GameJam) error {
	_, err := st.db.Exec(ctx, "INSERT INTO game_jams (name, url, content, start_date, end_date, voting_end_date, hide_results, hide_submissions) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		jam.Name, jam.URL, jam.Content, jam.StartDate, jam.EndDate, jam.VotingEndDate, jam.HideResults, jam.HideSubmissions)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) GetJamID(ctx context.Context, jamURL string) (int, error) {
	row := st.db.QueryRow(ctx, "SELECT game_jam_id FROM game_jams WHERE url = $1", jamURL)

	var gameJamID int
	err := row.Scan(&gameJamID)
	if err != nil {
		return -1, err
	}

	return gameJamID, nil
}

func (st *storage) GetJam(ctx context.Context, jamID int) (*gamejam.GameJam, error) {
	row := st.db.QueryRow(ctx, "SELECT game_jam_id, name, url, content, start_date, end_date, voting_end_date, hide_results, hide_submissions FROM game_jams WHERE game_jam_id = $1", jamID)

	var gameJam gamejam.GameJam
	err := row.Scan(&gameJam.ID, &gameJam.Name, &gameJam.URL, &gameJam.Content, &gameJam.StartDate, &gameJam.EndDate, &gameJam.VotingEndDate, &gameJam.HideResults, &gameJam.HideSubmissions)
	if err != nil {
		return nil, err
	}

	return &gameJam, nil
}

func (st *storage) UpdateJam(ctx context.Context, jam gamejam.GameJam) error {
	_, err := st.db.Exec(ctx, "UPDATE game_jams SET name = $1, url = $2, content = $3, start_date = $4, end_date = $5, voting_end_date = $6, hide_results = $7, hide_submissions = $8 WHERE game_jam_id = $9",
		jam.Name, jam.URL, jam.Content, jam.StartDate, jam.EndDate, jam.VotingEndDate, jam.HideResults, jam.HideSubmissions, jam.ID)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) GetJams(ctx context.Context) ([]gamejam.GameJam, error) {
	rows, err := st.db.Query(ctx, "SELECT game_jam_id, name, url, content, start_date, end_date, voting_end_date, hide_results, hide_submissions FROM game_jams")
	if err != nil {
		return nil, err
	}

	var gameJams []gamejam.GameJam
	for rows.Next() {
		var gameJam gamejam.GameJam
		err = rows.Scan(&gameJam.ID, &gameJam.Name, &gameJam.URL, &gameJam.Content, &gameJam.StartDate, &gameJam.EndDate, &gameJam.VotingEndDate, &gameJam.HideResults, &gameJam.HideSubmissions)
		if err != nil {
			return nil, err
		}
		gameJams = append(gameJams, gameJam)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gameJams, nil
}

func (st *storage) DeleteJam(ctx context.Context, gameJamID int) error {
	_, err := st.db.Exec(ctx, "DELETE FROM game_jams WHERE game_jam_id = $1", gameJamID)
	if err != nil {
		return err
	}
	return nil
}
