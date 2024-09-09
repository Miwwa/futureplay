package futureplay_assignment

import (
	"context"
	"time"
)

type CompetitionId string

type Competition struct {
	Id        CompetitionId       `json:"id"`
	CreatedAt time.Time           `json:"created_at"`
	EndsAt    time.Time           `json:"ends_at"`
	Records   []CompetitionRecord `json:"leaderboard"`
}

type CompetitionRecord struct {
	PlayerId PlayerId `json:"player_id"`
	Score    int      `json:"score"`
}

type CompetitionStorage interface {
	Create(ctx context.Context, players []PlayerId, endsAt time.Time) (Competition, error)
	IncrementScore(ctx context.Context, competitionId CompetitionId, playerId PlayerId, score int) error
	GetCompetitionById(ctx context.Context, competitionId CompetitionId) (Competition, error)
	GetLastCompetitionByPlayerId(ctx context.Context, playerId PlayerId) (Competition, error)
}

// MongoCompetitionStorage implement CompetitionStorage for production storage in database
// use Competition struct as bson document
type MongoCompetitionStorage struct {
}
