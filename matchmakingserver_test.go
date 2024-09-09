package futureplay_assignment

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMatchmakingServiceImpl_Join(t *testing.T) {
	options := MatchmakingOptions{
		MaxPlayerInMatch: 3,
		MatchInterval:    0,
		EntryTimeout:     0,
		MaxLevelDiff:     0,
	}

	playerData := []PlayerData{
		{
			Id:          "1",
			Level:       1,
			CountryCode: "US",
		},
		{
			Id:          "2",
			Level:       1,
			CountryCode: "US",
		},
		{
			Id:          "3",
			Level:       1,
			CountryCode: "US",
		},
		{
			Id:          "4",
			Level:       1,
			CountryCode: "US",
		},
		{
			Id:          "5",
			Level:       1,
			CountryCode: "US",
		},
		{
			Id:          "6",
			Level:       1,
			CountryCode: "US",
		},
	}

	t.Run("Join form a group", func(t *testing.T) {
		service := NewMatchmakingServiceImpl(options)

		// expecting form a group every 3 players joined
		match, err := service.Join(playerData[0])
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(playerData[1])
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(playerData[2])
		assert.Nil(t, err)
		assert.True(t, match.IsMatchFound)

		match, err = service.Join(playerData[3])
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(playerData[4])
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(playerData[5])
		assert.Nil(t, err)
		assert.True(t, match.IsMatchFound)
	})

	t.Run("Join the same player twice", func(t *testing.T) {
		service := NewMatchmakingServiceImpl(options)

		match, err := service.Join(playerData[0])
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(playerData[0])
		assert.EqualError(t, err, PlayerAlreadyInMatchmakingError.Error())
	})

	t.Run("Join players with different levels", func(t *testing.T) {
		service := NewMatchmakingServiceImpl(MatchmakingOptions{
			MaxPlayerInMatch: 3,
			MatchInterval:    0,
			EntryTimeout:     0,
			MaxLevelDiff:     1,
		})

		// mixing players with different levels
		// expecting form two groups with different levels

		match, err := service.Join(PlayerData{
			Id:    "1",
			Level: 1,
		})
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(PlayerData{
			Id:    "2",
			Level: 1,
		})
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(PlayerData{
			Id:    "3",
			Level: 3,
		})
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(PlayerData{
			Id:    "4",
			Level: 1,
		})
		assert.Nil(t, err)
		assert.True(t, match.IsMatchFound)

		match, err = service.Join(PlayerData{
			Id:    "5",
			Level: 4,
		})
		assert.Nil(t, err)
		assert.False(t, match.IsMatchFound)

		match, err = service.Join(PlayerData{
			Id:    "6",
			Level: 4,
		})
		assert.Nil(t, err)
		assert.True(t, match.IsMatchFound)
	})
}

func TestMatchmakingServiceImpl_ProcessQueue(t *testing.T) {
	options := MatchmakingOptions{
		MaxPlayerInMatch: 3,
		MatchInterval:    1 * time.Second,
		EntryTimeout:     2 * time.Second,
		MaxLevelDiff:     0,
	}

	playerData := []PlayerData{
		{
			Id:          "1",
			Level:       1,
			CountryCode: "US",
		},
		{
			Id:          "2",
			Level:       1,
			CountryCode: "US",
		},
	}

	t.Run("Process Queue", func(t *testing.T) {
		ms := NewMatchmakingServiceImpl(options)
		for _, player := range playerData {
			result, err := ms.Join(player)
			assert.Nil(t, err)
			assert.False(t, result.IsMatchFound)
		}

		callsCount := 0
		stop := ms.Start(func(match MatchmakingResult) {
			callsCount++
			assert.True(t, match.IsMatchFound)
		})

		time.Sleep(options.EntryTimeout + options.MatchInterval)
		stop()

		assert.Equal(t, 1, callsCount)
		if len(ms.playersPool) != 0 || len(ms.queue) != 0 {
			t.Errorf("Expected empty player pool and queue, got: %v and %v", ms.playersPool, ms.queue)
		}
	})
}
