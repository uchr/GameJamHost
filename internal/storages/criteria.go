package storages

import (
	"context"

	"GameJamPlatform/internal/models/gamejams"
)

func (st *storage) createCriteria(ctx context.Context, criteria gamejams.Criteria) error {
	_, err := st.db.Exec(ctx, "INSERT INTO criteria (title, description, jam_id) VALUES ($1, $2, $3)", criteria.Title, criteria.Description, criteria.JamID)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) getCriteria(ctx context.Context, jamID int) ([]gamejams.Criteria, error) {
	rows, err := st.db.Query(ctx, "SELECT criteria_id, title, description, jam_id FROM criteria WHERE jam_id = $1", jamID)
	if err != nil {
		return nil, err
	}

	var criteria []gamejams.Criteria
	for rows.Next() {
		var c gamejams.Criteria
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.JamID)
		if err != nil {
			return nil, err
		}
		criteria = append(criteria, c)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return criteria, nil
}

func (st *storage) deleteCriteria(ctx context.Context, jamID int) error {
	_, err := st.db.Exec(ctx, "DELETE FROM criteria WHERE jam_id = $1", jamID)
	if err != nil {
		return err
	}
	return nil
}
