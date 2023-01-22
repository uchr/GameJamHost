package storages

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"GameJamPlatform/internal/models/gamejams"
)

func (st *storage) CreateJam(ctx context.Context, jam gamejams.GameJam) error {
	row := st.db.QueryRow(ctx, "INSERT INTO game_jams (user_id, title, url, content, cover_image, start_date, end_date, voting_end_date, hide_results, hide_submissions) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING game_jam_id",
		jam.UserID, jam.Title, jam.URL, jam.Content, jam.CoverImageURL, jam.StartDate, jam.EndDate, jam.VotingEndDate, jam.HideResults, jam.HideSubmissions)

	err := row.Scan(&jam.ID)
	if err != nil {
		return err
	}

	for _, c := range jam.Criteria {
		c.JamID = jam.ID
		err = st.createCriteria(ctx, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *storage) GetJamByID(ctx context.Context, jamID int) (*gamejams.GameJam, error) {
	row := st.db.QueryRow(ctx, "SELECT game_jam_id, user_id, title, url, content, cover_image, start_date, end_date, voting_end_date, hide_results, hide_submissions FROM game_jams WHERE game_jam_id = $1", jamID)

	var jam gamejams.GameJam
	err := row.Scan(&jam.ID, &jam.UserID, &jam.Title, &jam.URL, &jam.Content, &jam.CoverImageURL, &jam.StartDate, &jam.EndDate, &jam.VotingEndDate, &jam.HideResults, &jam.HideSubmissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	jam.Criteria, err = st.getCriteria(ctx, jamID)
	if err != nil {
		return nil, err
	}

	return &jam, nil
}

func (st *storage) GetJamByURL(ctx context.Context, jamURL string) (*gamejams.GameJam, error) {
	row := st.db.QueryRow(ctx, "SELECT game_jam_id, user_id, title, url, content, cover_image, start_date, end_date, voting_end_date, hide_results, hide_submissions FROM game_jams WHERE url = $1", jamURL)

	var jam gamejams.GameJam
	err := row.Scan(&jam.ID, &jam.UserID, &jam.Title, &jam.URL, &jam.Content, &jam.CoverImageURL, &jam.StartDate, &jam.EndDate, &jam.VotingEndDate, &jam.HideResults, &jam.HideSubmissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	jam.Criteria, err = st.getCriteria(ctx, jam.ID)
	if err != nil {
		return nil, err
	}

	return &jam, nil
}

func (st *storage) UpdateJam(ctx context.Context, jam gamejams.GameJam) error {
	_, err := st.db.Exec(ctx, "UPDATE game_jams SET title = $1, url = $2, content = $3, cover_image = $4, start_date = $5, end_date = $6, voting_end_date = $7, hide_results = $8, hide_submissions = $9 WHERE game_jam_id = $10",
		jam.Title, jam.URL, jam.Content, jam.CoverImageURL, jam.StartDate, jam.EndDate, jam.VotingEndDate, jam.HideResults, jam.HideSubmissions, jam.ID)
	if err != nil {
		return err
	}

	err = st.deleteCriteria(ctx, jam.ID)
	if err != nil {
		return err
	}

	for _, c := range jam.Criteria {
		c.JamID = jam.ID
		err = st.createCriteria(ctx, c)
		if err != nil {
			return err
		}
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

		criteria, err := st.getCriteria(ctx, jam.ID)
		if err != nil {
			return nil, err
		}
		jam.Criteria = criteria

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

		criteria, err := st.getCriteria(ctx, jam.ID)
		if err != nil {
			return nil, err
		}
		jam.Criteria = criteria

		gameJams = append(gameJams, jam)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gameJams, nil
}

func (st *storage) DeleteJam(ctx context.Context, jamID int) error {
	_, err := st.db.Exec(ctx, "DELETE FROM criteria WHERE jam_id = $1", jamID)
	if err != nil {
		return err
	}

	_, err = st.db.Exec(ctx, "DELETE FROM game_jams WHERE game_jam_id = $1", jamID)
	if err != nil {
		return err
	}
	return nil
}
