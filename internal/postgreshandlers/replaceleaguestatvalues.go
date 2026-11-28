package postgreshandlers

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"fmt"
)

type leagueKey struct {
	gameId   int
	leagueId int
}

// ReplaceLeagueStatValues replaces the scoring rules in league_scoring_value_information for
// every league in values. A StatValueTopic message holds a league's complete set of rules, so
// the league's old rows are deleted first; a stat the league stopped scoring doesn't linger.
// It all runs in one transaction, so readers (the Java API's best-available endpoint) never see
// a league with no rules, and a failure keeps the previous rules. Returns the rows written.
func ReplaceLeagueStatValues(ctx context.Context, values []leagueclasses.LeagueStatsValueInfo) (int, error) {
	if len(values) == 0 {
		logger.LogInfo(ctx, "No league stat values to replace")
		return 0, nil
	}
	leagues := make(map[leagueKey]bool)
	for _, value := range values {
		if value.GameId <= 0 || value.LeagueId <= 0 || value.StatId <= 0 {
			return 0, fmt.Errorf("league stat value has a missing id: game_id=%d league_id=%d stat_id=%d",
				value.GameId, value.LeagueId, value.StatId)
		}
		leagues[leagueKey{value.GameId, value.LeagueId}] = true
	}

	db, err := GetDB(ctx)
	if err != nil {
		return 0, err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// A no-op once the transaction is committed.
	defer func() { _ = tx.Rollback() }()

	for league := range leagues {
		deleted, err := tx.ExecContext(ctx,
			`DELETE FROM league_scoring_value_information WHERE game_id = $1 AND league_id = $2`,
			league.gameId, league.leagueId)
		if err != nil {
			return 0, fmt.Errorf("deleting league %d.l.%d stat values: %w", league.gameId, league.leagueId, err)
		}
		deletedCount, _ := deleted.RowsAffected()
		logger.LogInfo(ctx, "Deleted league stat values", "game_id", league.gameId,
			"league_id", league.leagueId, "count", deletedCount)
	}
	// ON CONFLICT keeps the last value if a message repeats a stat.
	for _, value := range values {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO league_scoring_value_information (game_id, league_id, league_stat_id, league_stat_value)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (game_id, league_id, league_stat_id) DO UPDATE SET league_stat_value = EXCLUDED.league_stat_value`,
			value.GameId, value.LeagueId, value.StatId, value.StatValue)
		if err != nil {
			return 0, fmt.Errorf("inserting league %d.l.%d stat %d: %w", value.GameId, value.LeagueId, value.StatId, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	logger.LogInfo(ctx, "Replaced league stat values", "leagues", len(leagues), "count", len(values))
	return len(values), nil
}
