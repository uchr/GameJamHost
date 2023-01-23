package gamejams

import "time"

func (jam GameJam) GetState() (JamState, time.Time) {
	now := time.Now().UTC()

	if now.Before(jam.StartDate) {
		return JamStateNotStarted, jam.StartDate
	} else if now.Before(jam.EndDate) {
		return JamStateStarted, jam.EndDate
	} else if now.Before(jam.VotingEndDate) {
		return JamStateEnded, jam.VotingEndDate
	}

	return JamStateVotingEnded, jam.VotingEndDate
}
