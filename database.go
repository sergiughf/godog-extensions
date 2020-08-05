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

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to postgres while executing the db clean up: %+v", err.Error())
	}

	ctx.BeforeScenario(func(scenario *godog.Scenario) {
		if _, err := db.Exec(cleanupQuery); err != nil {
			log.Fatalf("failed to execute db cleanup before scenarios: %+v", err.Error())
		}
	})
}
