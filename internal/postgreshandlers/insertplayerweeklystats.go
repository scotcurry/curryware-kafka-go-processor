package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/statsclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"fmt"
	"strconv"
	"strings"
)

const statsBatchSize = 500
const statsColumnsPerRow = 5

// InsertPlayerWeeklyStats performs a bulk parameterized insert of player weekly stats.
// Rows are batched to stay within Postgres parameter limits.
func InsertPlayerWeeklyStats(statsJson []statsclasses.PlayerWeeklyStatsInfo) int64 {
	if len(statsJson) == 0 {
		return 0
	}

	var totalCount int64
	for batchStart := 0; batchStart < len(statsJson); batchStart += statsBatchSize {
		batchEnd := batchStart + statsBatchSize
		if batchEnd > len(statsJson) {
			batchEnd = len(statsJson)
		}
		batch := statsJson[batchStart:batchEnd]

		count, err := insertStatsBatch(batch)
		if err != nil {
			logger.LogError("Error inserting player stats batch",
				"error", err.Error(),
				"batch_start", batchStart,
				"batch_size", len(batch))
			return totalCount
		}
		totalCount += count
	}

	logger.LogInfo("Player stats inserted", "total_rows", totalCount)
	return totalCount
}

func insertStatsBatch(batch []statsclasses.PlayerWeeklyStatsInfo) (int64, error) {
	valuePlaceholders := make([]string, len(batch))
	params := make([]any, 0, len(batch)*statsColumnsPerRow)

	for i, stat := range batch {
		offset := i * statsColumnsPerRow
		valuePlaceholders[i] = "($" + strconv.Itoa(offset+1) +
			", $" + strconv.Itoa(offset+2) +
			", $" + strconv.Itoa(offset+3) +
			", $" + strconv.Itoa(offset+4) +
			", $" + strconv.Itoa(offset+5) + ")"

		params = append(params,
			stat.PlayerId,
			stat.PlayerGameKey,
			stat.PlayerStatWeek,
			stat.StatId,
			stat.StatValue,
		)
	}

	sqlStatement := fmt.Sprintf(
		"INSERT INTO player_weekly_stats (player_id, player_game_id, player_stats_week, stat_id, stat_value) VALUES %s ON CONFLICT DO NOTHING",
		strings.Join(valuePlaceholders, ", "),
	)

	return ExecStatement(sqlStatement, params...)
}
