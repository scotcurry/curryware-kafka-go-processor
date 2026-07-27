package postgreshandlers

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/teamclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertTeamWeeklyStandingsInfo(ctx context.Context, standingsInfo []teamclasses.TeamWeeklyStandingsInfo) int {
	sqlStatement := `INSERT INTO team_weekly_standings (team_key, team_id, week_key, rank, points_for,
                         points_against, wins, losses, ties, number_of_moves)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                ON CONFLICT (team_key, week_key) DO UPDATE SET
                    team_id = EXCLUDED.team_id,
                    rank = EXCLUDED.rank,
                    points_for = EXCLUDED.points_for,
                    points_against = EXCLUDED.points_against,
                    wins = EXCLUDED.wins,
                    losses = EXCLUDED.losses,
                    ties = EXCLUDED.ties,
                    number_of_moves = EXCLUDED.number_of_moves`

	for counter := range standingsInfo {
		teamKey := standingsInfo[counter].TeamKey
		teamId := standingsInfo[counter].TeamId
		weekKey := standingsInfo[counter].WeekKey
		rank := standingsInfo[counter].Rank
		pointsFor := standingsInfo[counter].PointsFor
		pointsAgainst := standingsInfo[counter].PointsAgainst
		wins := standingsInfo[counter].Wins
		losses := standingsInfo[counter].Losses
		ties := standingsInfo[counter].Ties
		numberOfMoves := standingsInfo[counter].NumberOfMoves

		count, err := ExecStatement(ctx, sqlStatement, teamKey, teamId, weekKey, rank, pointsFor,
			pointsAgainst, wins, losses, ties, numberOfMoves)
		if err != nil {
			logger.LogError(ctx, "Error inserting team weekly standings record", "error", err.Error(), "team_key", teamKey, "week_key", weekKey)
			continue
		}
		logger.LogInfo(ctx, "Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo(ctx, "Done inserting team weekly standings records", "total_records", strconv.Itoa(len(standingsInfo)))
	return len(standingsInfo)
}
