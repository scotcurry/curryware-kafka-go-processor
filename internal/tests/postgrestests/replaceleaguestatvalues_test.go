// Runs ReplaceLeagueStatValues against a real, throwaway Postgres. Skipped unless
// CURRYWARE_TEST_POSTGRES=1, and it refuses to run against currywarefantasy, because it drops
// and recreates league_scoring_value_information from ../curryware-fantasy-database. To run:
//
//	podman run -d --rm --name go-processor-test-pg -p 55435:5432 -e POSTGRES_USER=test \
//	    -e POSTGRES_PASSWORD=test -e POSTGRES_DB=curryware_test docker.io/library/postgres:17
//	CURRYWARE_TEST_POSTGRES=1 POSTGRES_SERVER=localhost POSTGRES_PORT=55435 \
//	    POSTGRES_USERNAME=test POSTGRES_PASSWORD=test POSTGRES_DATABASE=curryware_test \
//	    go test ./internal/tests/postgrestests/...
//	podman stop go-processor-test-pg
//
// The DDL is read from ../curryware-fantasy-database (override with CURRYWARE_SCHEMA_DIR).
package postgrestests

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	"curryware-kafka-go-processor/internal/postgreshandlers"
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const productionDatabase = "currywarefantasy"

func setUp(t *testing.T) *sql.DB {
	t.Helper()
	if os.Getenv("CURRYWARE_TEST_POSTGRES") != "1" {
		t.Skip("set CURRYWARE_TEST_POSTGRES=1 and POSTGRES_* to run against a throwaway Postgres")
	}
	if os.Getenv("POSTGRES_DATABASE") == productionDatabase {
		t.Fatalf("POSTGRES_DATABASE must not be %s", productionDatabase)
	}
	schemaDir := os.Getenv("CURRYWARE_SCHEMA_DIR")
	if schemaDir == "" {
		// go test runs in the package directory: internal/tests/postgrestests.
		schemaDir = filepath.Join("..", "..", "..", "..", "curryware-fantasy-database")
	}
	ddl, err := os.ReadFile(filepath.Join(schemaDir, "create_league_scoring_value_information.sql"))
	if err != nil {
		t.Fatalf("reading DDL: %v", err)
	}
	ctx := context.Background()
	db, err := postgreshandlers.GetDB(ctx)
	if err != nil {
		t.Fatalf("connecting: %v", err)
	}
	if _, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS league_scoring_value_information"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
		t.Fatalf("creating table: %v", err)
	}
	return db
}

type row struct {
	LeagueId int
	StatId   int
	Value    *string
}

func storedRows(t *testing.T, db *sql.DB) []row {
	t.Helper()
	rows, err := db.Query(`SELECT league_id, league_stat_id, league_stat_value::text
		FROM league_scoring_value_information ORDER BY league_id, league_stat_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var result []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.LeagueId, &r.StatId, &r.Value); err != nil {
			t.Fatal(err)
		}
		result = append(result, r)
	}
	return result
}

func value(v float64) *float64 { return &v }
func text(s string) *string    { return &s }

func stat(league, statId int, v *float64) leagueclasses.LeagueStatsValueInfo {
	return leagueclasses.LeagueStatsValueInfo{GameId: 470, LeagueId: league, StatId: statId, StatValue: v}
}

func TestReplaceInsertsExactDecimalsAndNulls(t *testing.T) {
	db := setUp(t)

	count, err := postgreshandlers.ReplaceLeagueStatValues(context.Background(), []leagueclasses.LeagueStatsValueInfo{
		stat(838859, 4, value(0.04)), stat(838859, 6, value(-1)), stat(838859, 57, nil)})
	if err != nil || count != 3 {
		t.Fatalf("got count=%d err=%v", count, err)
	}

	want := []row{{838859, 4, text("0.04")}, {838859, 6, text("-1")}, {838859, 57, nil}}
	if got := storedRows(t, db); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestReplaceDropsStatsTheLeagueNoLongerScoresAndLeavesOtherLeagues(t *testing.T) {
	db := setUp(t)
	ctx := context.Background()
	if _, err := postgreshandlers.ReplaceLeagueStatValues(ctx, []leagueclasses.LeagueStatsValueInfo{
		stat(838859, 4, value(0.04)), stat(838859, 11, value(0.5)), stat(523652, 11, value(1))}); err != nil {
		t.Fatal(err)
	}

	// 838859 switches to full PPR and stops scoring passing yards.
	if _, err := postgreshandlers.ReplaceLeagueStatValues(ctx, []leagueclasses.LeagueStatsValueInfo{
		stat(838859, 11, value(1))}); err != nil {
		t.Fatal(err)
	}

	want := []row{{523652, 11, text("1")}, {838859, 11, text("1")}}
	if got := storedRows(t, db); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestReplaceWithAMissingIdKeepsTheExistingRules(t *testing.T) {
	db := setUp(t)
	ctx := context.Background()
	if _, err := postgreshandlers.ReplaceLeagueStatValues(ctx, []leagueclasses.LeagueStatsValueInfo{
		stat(838859, 4, value(0.04))}); err != nil {
		t.Fatal(err)
	}

	_, err := postgreshandlers.ReplaceLeagueStatValues(ctx, []leagueclasses.LeagueStatsValueInfo{
		stat(838859, 5, value(4)), stat(838859, 0, value(1))})
	if err == nil {
		t.Fatal("expected an error for stat id 0")
	}

	want := []row{{838859, 4, text("0.04")}}
	if got := storedRows(t, db); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestReplaceWithNoValuesChangesNothing(t *testing.T) {
	db := setUp(t)
	ctx := context.Background()
	if _, err := postgreshandlers.ReplaceLeagueStatValues(ctx, []leagueclasses.LeagueStatsValueInfo{
		stat(838859, 4, value(0.04))}); err != nil {
		t.Fatal(err)
	}

	count, err := postgreshandlers.ReplaceLeagueStatValues(ctx, nil)
	if err != nil || count != 0 {
		t.Fatalf("got count=%d err=%v", count, err)
	}
	if got := storedRows(t, db); len(got) != 1 {
		t.Errorf("expected the existing row to stay, got %+v", got)
	}
}
