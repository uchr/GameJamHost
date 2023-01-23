package gamejammanager

import (
	"context"
	"sort"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/voting"
)

func (jm *gameJamManager) VoteGame(ctx context.Context, vote []gamejams.Vote) error {
	for _, v := range vote {
		err := jm.repo.AddVote(ctx, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (jm *gameJamManager) GetGameResult(ctx context.Context, jamURL string, gameURL string) (*voting.GameResult, error) {
	jam, err := jm.repo.GetJamByURL(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	game, err := jm.repo.GetGame(ctx, jam.ID, gameURL)
	if err != nil {
		return nil, err
	}

	result := &voting.GameResult{
		Scores:   make(map[int]float32),
		Criteria: make(map[int]*gamejams.Criteria),
	}

	for _, c := range jam.Criteria {
		votes, err := jm.repo.GetVote(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		var sum float32
		for _, v := range votes {
			if v.GameID == game.ID {
				sum += float32(v.Value)
			}
		}

		result.Scores[c.ID] = sum / float32(len(votes))
		result.Criteria[c.ID] = &c
	}

	return result, nil
}

func (jm *gameJamManager) GetJamResult(ctx context.Context, jamURL string) (*voting.JamResult, error) {
	jam, err := jm.repo.GetJamByURL(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	games, err := jm.repo.GetGames(ctx, jam.ID)
	if err != nil {
		return nil, err
	}

	gameResults := make(map[int]*voting.GameResult)
	for _, g := range games {
		result, err := jm.GetGameResult(ctx, jamURL, g.URL)
		if err != nil {
			return nil, err
		}

		gameResults[g.ID] = result
	}

	criteriaResults := make([]voting.CriteriaResult, 0, len(jam.Criteria))
	for _, c := range jam.Criteria {
		criteriaResult := voting.CriteriaResult{
			Criteria: c,
			Scores:   make(map[int]float32),
			Games:    make([]*gamejams.Game, 0),
		}

		for i, g := range games {
			criteriaResult.Scores[g.ID] = gameResults[g.ID].Scores[c.ID]
			criteriaResult.Games = append(criteriaResult.Games, &games[i])
		}

		sort.Slice(criteriaResult.Games, func(i, j int) bool {
			return criteriaResult.Scores[criteriaResult.Games[i].ID] > criteriaResult.Scores[criteriaResult.Games[j].ID]
		})

		criteriaResults = append(criteriaResults, criteriaResult)
	}

	return &voting.JamResult{
		CriteriaResults: criteriaResults,
	}, nil
}
