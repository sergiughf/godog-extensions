package extension

import (
	"database/sql"
	"log"

	"github.com/cucumber/godog"
)

// NewPostgresCleanup executes a query before each scenario in order to clean the postgres db
func NewPostgresCleanup(ctx *godog.ScenarioContext, postgresDSN string) {
	cleanupQuery := `
        DO
        $func$
        begin
            EXECUTE
            (SELECT
                'TRUNCATE TABLE '
                || string_agg(format('%I.%I', schemaname, tablename), ', ')
                || ' RESTART IDENTITY CASCADE'
            FROM   pg_tables
            WHERE  schemaname = 'public'
            );
        END
        $func$ LANGUAGE plpgsql;
    `

	var db *sql.DB
	var err error
	ctx.BeforeScenario(func(scenario *godog.Scenario) {
		db, err = sql.Open("postgres", postgresDSN)
		if err != nil {
			log.Fatalf("failed to connect to postgres while executing the db clean up: %+v", err.Error())
		}

		if _, err := db.Exec(cleanupQuery); err != nil {
			log.Fatalf("failed to execute db cleanup before scenarios: %+v", err.Error())
		}
	})

	ctx.AfterScenario(func(scenario *godog.Scenario, err error) {
		if db != nil {
			if err := db.Close(); err != nil {
				log.Fatalf("failed to execute db cleanup after scenarios: %+v", err.Error())
			}
		}
	})
}
