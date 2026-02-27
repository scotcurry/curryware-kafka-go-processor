package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed sqltemplates/player_weekly_stats.sql
var playerWeeklyStatsInsertSQL string

//go:embed sqltemplates/league_stat_info_insert.sql
var leagueStatInfoInsertSQL string

var sqlTemplates = map[string]string{
	"multiple_player_stats_input_statement": playerWeeklyStatsInsertSQL,
	"league_stat_info_insert_statement":     leagueStatInfoInsertSQL,
}

func GetSqlTemplate(templateName string) string {
	tmpl, ok := sqlTemplates[templateName]
	if !ok {
		logger.LogError(fmt.Sprintf("SQL template not found: %s", templateName))
		return ""
	}
	return strings.TrimSpace(tmpl)
}
