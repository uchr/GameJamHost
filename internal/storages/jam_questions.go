package storages

import (
	"context"

	"GameJamPlatform/internal/models/gamejams"
)

func (st *storage) createQuestion(ctx context.Context, question gamejams.JamQuestion) error {
	_, err := st.db.Exec(ctx, "INSERT INTO jam_questions (jam_id, title, description, hidden_criteria) VALUES ($1, $2, $3, $4)", question.JamID, question.Title, question.Description, question.HiddenCriteria)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) getQuestions(ctx context.Context, jamID int) ([]gamejams.JamQuestion, error) {
	rows, err := st.db.Query(ctx, "SELECT question_id, jam_id, title, description, hidden_criteria FROM jam_questions WHERE jam_id = $1", jamID)
	if err != nil {
		return nil, err
	}

	var questions []gamejams.JamQuestion
	for rows.Next() {
		var q gamejams.JamQuestion
		err := rows.Scan(&q.ID, &q.JamID, &q.Title, &q.Description, &q.HiddenCriteria)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return questions, nil
}

func (st *storage) deleteQuestions(ctx context.Context, jamID int) error {
	_, err := st.db.Exec(ctx, "DELETE FROM jam_questions WHERE jam_id = $1", jamID)
	if err != nil {
		return err
	}
	return nil
}
