package voting

import "GameJamPlatform/internal/models/gamejams"

type GameResult struct {
	Scores   map[int]float32            // CriteriaID -> Score
	Criteria map[int]*gamejams.Criteria // CriteriaID -> Criteria
}

type CriteriaResult struct {
	Criteria gamejams.Criteria
	Games    []*gamejams.Game // Sorted by score
	Scores   map[int]float32  // GameID -> Score
}

type JamResult struct {
	CriteriaResults []CriteriaResult
}
