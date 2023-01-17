package storages

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"GameJamPlatform/internal/models/gamejams"
)

func (st *storage) CreateJam(ctx context.Context, jam gamejams.GameJam) error {
	_, err := st.db.Exec(ctx, "INSERT INTO game_jams (user_id, title, url, content, cover_image, start_date, end_date, voting_end_date, hide_results, hide_submissions) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		jam.UserID, jam.Title, jam.URL, jam.Content, jam.CoverImageURL, jam.StartDate, jam.EndDate, jam.VotingEndDate, jam.HideResults, jam.HideSubmissions)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, ErrNotFound
		}
		return -1, err
	}

	return gameJamID, nil
}

func (st *storage) GetJam(ctx context.Context, jamID int) (*gamejams.GameJam, error) {
	row := st.db.QueryRow(ctx, "SELECT game_jam_id, user_id, title, url, content, cover_image, start_date, end_date, voting_end_date, hide_results, hide_submissions FROM game_jams WHERE game_jam_id = $1", jamID)

	var gameJam gamejams.GameJam
	err := row.Scan(&gameJam.ID, &gameJam.UserID, &gameJam.Title, &gameJam.URL, &gameJam.Content, &gameJam.CoverImageURL, &gameJam.StartDate, &gameJam.EndDate, &gameJam.VotingEndDate, &gameJam.HideResults, &gameJam.HideSubmissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &gameJam, nil
}

func (st *storage) UpdateJam(ctx context.Context, jam gamejams.GameJam) error {
	_, err := st.db.Exec(ctx, "UPDATE game_jams SET title = $1, url = $2, content = $3, cover_image = $4, start_date = $5, end_date = $6, voting_end_date = $7, hide_results = $8, hide_submissions = $9 WHERE game_jam_id = $9",
		jam.Title, jam.URL, jam.Content, jam.CoverImageURL, jam.StartDate, jam.EndDate, jam.VotingEndDate, jam.HideResults, jam.HideSubmissions, jam.ID)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) GetJams(ctx context.Context) ([]gamejams.GameJam, error) {
	rows, err := st.db.Query(ctx, "SELECT game_jam_id, user_id, title, url, content, cover_image, start_date, end_date, voting_end_date, hide_results, hide_submissions FROM game_jams")
	if err != nil {
		return nil, err
	}

	var gameJams []gamejams.GameJam
	for rows.Next() {
		var jam gamejams.GameJam
		err = rows.Scan(&jam.ID, &jam.UserID, &jam.Title, &jam.URL, &jam.Content, &jam.CoverImageURL, &jam.StartDate, &jam.EndDate, &jam.VotingEndDate, &jam.HideResults, &jam.HideSubmissions)
		if err != nil {
			return nil, err
		}
		gameJams = append(gameJams, jam)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gameJams, nil
}

func (st *storage) GetJamsByUserID(ctx context.Context, userID int) ([]gamejams.GameJam, error) {
	rows, err := st.db.Query(ctx, "SELECT game_jam_id, title, url, content, cover_image, start_date, end_date, voting_end_date, hide_results, hide_submissions FROM game_jams WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	var gameJams []gamejams.GameJam
	for rows.Next() {
		var jam gamejams.GameJam
		err = rows.Scan(&jam.ID, &jam.Title, &jam.URL, &jam.Content, &jam.CoverImageURL, &jam.StartDate, &jam.EndDate, &jam.VotingEndDate, &jam.HideResults, &jam.HideSubmissions)
		if err != nil {
			return nil, err
		}
		gameJams = append(gameJams, jam)
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
