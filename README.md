## Getting Started

**Prerequisites**

- Go 1.22 or later

**Run**

```shell
git clone https://github.com/Miwwa/futureplay
cd futureplay
go test
```

Look at `matchmakingserver_test.go` for usage examples

## Key design decisions

**Player's waiting time**: the success of a matchmaking request, or the length of time needed to find a match,
cannot be guaranteed as both are dependent on the pool of users active in the matchmaker.
This solution focused on player's waiting time and guaranties to start a new competition in 30s for new joined player.
In this case, competition starts with less than 10 players. We can backfill it with new joined players or bots.

**In-memory storage**: for simplicity and demonstration purposes, service using in-memory storage

### Search by country proposal

Let's map every country with a number. Number depends on country's region - Europe, Asia, America, etc.
Countries from the same region should have close numbers. For example:
```json
{
    "FI": 100,
    "SE": 101,
    "NO": 102,

    "US": 200,
    "CA": 201,
    ...
}
```
In matchmaking function, we calculate a *distance* between player's countries and make a matchmaking decision based on result.
Also, multiple countries can be mapped to the same number if we want to group players from these countries in any case.
We can manipulate countries index and max distance between players to optimize the search
