package storages

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"

	"GameJamPlatform/internal/models/gamejams"
)

func (st *storage) CreateGame(ctx context.Context, game gamejams.Game) error {
	_, err := st.db.Exec(ctx, "INSERT INTO games (game_jam_id, title, content, cover_image, screenshot_images, url, build, is_banned) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		game.GameJamID, game.Title, game.Content, game.CoverImageURL, pq.StringArray(game.ScreenshotURLs), game.URL, game.Build, game.IsBanned)

	if err != nil {
		return err
	}
	return nil
}

func (st *storage) GetGame(ctx context.Context, jamID int, gameURL string) (*gamejams.Game, error) {
	row := st.db.QueryRow(ctx, "SELECT game_id, game_jam_id, title, content, cover_image, screenshot_images, url, build, is_banned FROM games WHERE game_jam_id = $1 AND url = $2", jamID, gameURL)
	var game gamejams.Game
	err := row.Scan(&game.ID, &game.GameJamID, &game.Title, &game.Content, &game.CoverImageURL, (*pq.StringArray)(&game.ScreenshotURLs), &game.URL, &game.Build, &game.IsBanned)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &game, nil
}

func (st *storage) UpdateGame(ctx context.Context, game gamejams.Game) error {
	_, err := st.db.Exec(ctx, "UPDATE games SET title = $1, content = $2, cover_image = $3, screenshot_images = $4, url = $5, build = $6, is_banned = $7 WHERE game_id = $8 AND game_jam_id = $9",
		game.Title, game.Content, game.CoverImageURL, pq.StringArray(game.ScreenshotURLs), game.URL, game.Build, game.IsBanned, game.ID, game.GameJamID)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) GetGames(ctx context.Context, jamID int) ([]gamejams.Game, error) {
	rows, err := st.db.Query(ctx, "SELECT game_id, game_jam_id, title, content, cover_image, screenshot_images, url, build, is_banned FROM games WHERE game_jam_id = $1", jamID)
	if err != nil {
		return nil, err
	}

	var games []gamejams.Game
	for rows.Next() {
		var game gamejams.Game
		err = rows.Scan(&game.ID, &game.GameJamID, &game.Title, &game.Content, &game.CoverImageURL, (*pq.StringArray)(&game.ScreenshotURLs), &game.URL, &game.Build, &game.IsBanned)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (st *storage) DeleteGame(ctx context.Context, gameID int) error {
	_, err := st.db.Exec(ctx, "DELETE FROM games WHERE game_id = $1", gameID)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) BanGame(ctx context.Context, jamID int, gameURL string) error {
	_, err := st.db.Exec(ctx, "UPDATE games SET is_banned = true WHERE game_jam_id = $1 AND url = $2", jamID, gameURL)
	if err != nil {
		return err
	}
	return nil
}
