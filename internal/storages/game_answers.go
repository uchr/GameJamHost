package storages

import (
	"context"

	"GameJamPlatform/internal/models/gamejams"
)

func (st *storage) createAnswer(ctx context.Context, answer gamejams.GameAnswer) error {
	_, err := st.db.Exec(ctx, "INSERT INTO game_answers (game_id, question_id, answer) VALUES ($1, $2, $3)", answer.GameID, answer.QuestionID, answer.Answer)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) getAnswers(ctx context.Context, gameID int) ([]gamejams.GameAnswer, error) {
	rows, err := st.db.Query(ctx, "SELECT answer_id, game_id, question_id, answer FROM game_answers WHERE game_id = $1", gameID)
	if err != nil {
		return nil, err
	}

	var answers []gamejams.GameAnswer
	for rows.Next() {
		var a gamejams.GameAnswer
		err := rows.Scan(&a.ID, &a.GameID, &a.QuestionID, &a.Answer)
		if err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return answers, nil
}

func (st *storage) deleteAnswers(ctx context.Context, gameID int) error {
	_, err := st.db.Exec(ctx, "DELETE FROM game_answers WHERE game_id = $1", gameID)
	if err != nil {
		return err
	}
	return nil
}
