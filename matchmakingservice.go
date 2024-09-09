package futureplay_assignment

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

var (
	PlayerAlreadyInMatchmakingError = errors.New("player already in matchmaking")
)

// interface

// MatchmakingOptions represents the options for matchmaking process.
// MaxPlayerInMatch how many players can participate in leaderboard
// MatchInterval how quickly the matchmaker attempts to form matches
// EntryTimeout after this time new competition will be started even we don't found 'MaxPlayerInMatch' participants
// MaxLevelDiff looking for players with similar levels
type MatchmakingOptions struct {
	MaxPlayerInMatch int           `json:"max_player_in_match"`
	MatchInterval    time.Duration `json:"match_interval"`
	EntryTimeout     time.Duration `json:"entry_timeout"`
	MaxLevelDiff     int           `json:"max_level_diff"`
}

type MatchmakingResult struct {
	IsMatchFound bool
	Players      []PlayerData
}

type MatchFoundHandler = func(MatchmakingResult)
type StopFunc = func()

type MatchmakingService interface {
	// Join player to matchmaking queue, return new match if enough similar players are found
	Join(player PlayerData) (MatchmakingResult, error)
	// Start matching players every MatchInterval time interval
	// callback is called when match found for expiring player
	// return stop function
	Start(callback MatchFoundHandler) StopFunc
}

func DefaultMatchmakingOptions() MatchmakingOptions {
	// MaxPlayerInMatch = 10 - task requirement
	// MaxLevelDiff = 2 - for example
	// MatchInterval = 10s, EntryTimeout = 20s
	// we are trying to form a new competition every 10s, player entry expires in 20s, so next tick will be no longer than 30s
	// options can be adjusted to optimize time in queue
	return MatchmakingOptions{
		MaxPlayerInMatch: 10,
		MatchInterval:    10 * time.Second,
		EntryTimeout:     20 * time.Second,
		MaxLevelDiff:     2,
	}
}

func matchNotFound() MatchmakingResult {
	return MatchmakingResult{IsMatchFound: false}
}

func matchFound(players []PlayerData) MatchmakingResult {
	return MatchmakingResult{IsMatchFound: true, Players: players}
}

// impl

type MatchmakingServiceImpl struct {
	options MatchmakingOptions
	// active players index
	playersPool map[PlayerId]PlayerData
	// active players array sorted by join time, oldest are first
	queue []entry

	callback MatchFoundHandler
}

type entry struct {
	createdAt time.Time
	expiresAt time.Time
	player    PlayerData
}

func NewMatchmakingServiceImpl(options MatchmakingOptions) *MatchmakingServiceImpl {
	return &MatchmakingServiceImpl{
		options:     options,
		playersPool: make(map[PlayerId]PlayerData),
		queue:       make([]entry, 0),
	}
}

// Start matching players every MatchInterval time interval
func (m *MatchmakingServiceImpl) Start(callback MatchFoundHandler) func() {
	m.callback = callback

	timer := time.NewTicker(m.options.MatchInterval)
	stopTimer := make(chan bool)

	go func() {
		for {
			select {
			case <-stopTimer:
				return
			case <-timer.C:
				fmt.Println(time.Now().Format(time.TimeOnly), "Tick")
				m.processQueue()
			}
		}
	}()

	return func() {
		timer.Stop()
		stopTimer <- true
	}
}

func (m *MatchmakingServiceImpl) Join(player PlayerData) (MatchmakingResult, error) {
	_, ok := m.playersPool[player.Id]
	if ok {
		return MatchmakingResult{}, PlayerAlreadyInMatchmakingError
	}

	m.playersPool[player.Id] = player
	m.queue = append(m.queue, entry{
		createdAt: time.Now(),
		expiresAt: time.Now().Add(m.options.EntryTimeout),
		player:    player,
	})

	if len(m.playersPool) < m.options.MaxPlayerInMatch {
		return matchNotFound(), nil
	}

	match := m.findMatchForPlayer(player)

	isMaxPlayersMatch := len(match) == m.options.MaxPlayerInMatch
	if isMaxPlayersMatch {
		// found enough players for match
		m.removePlayers(match)
		return matchFound(match), nil
	}

	// we try to match players later, when someone else joins or N seconds passes
	return matchNotFound(), nil
}

func (m *MatchmakingServiceImpl) processQueue() {
	now := time.Now()
	for {
		if len(m.queue) == 0 {
			break
		}

		oldestEntry := m.queue[0]
		isExpired := now.After(oldestEntry.expiresAt)
		if !isExpired {
			break
		}

		// start a new competition after N seconds for player even we don't have enough players
		match := m.findMatchForPlayer(oldestEntry.player)
		m.removePlayers(match)

		m.callback(matchFound(match))
	}
}

func (m *MatchmakingServiceImpl) findMatchForPlayer(player PlayerData) []PlayerData {
	match := make([]PlayerData, 0, m.options.MaxPlayerInMatch)
	match = append(match, player)

	for _, other := range m.playersPool {
		if player.Id == other.Id {
			continue
		}

		isLevelMatches := abs(player.Level-other.Level) <= m.options.MaxLevelDiff
		if isLevelMatches {
			match = append(match, other)
			if len(match) == m.options.MaxPlayerInMatch {
				return match
			}
		}
	}

	return match
}

func (m *MatchmakingServiceImpl) removePlayers(players []PlayerData) {
	for _, player := range players {
		delete(m.playersPool, player.Id)
	}
	m.queue = slices.DeleteFunc(m.queue, func(e entry) bool {
		return slices.Contains(players, e.player)
	})
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return 1
}
