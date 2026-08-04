package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"permisson/internal/database"
	"permisson/internal/pkg/env"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env := env.Env{}
	env.Init()

	sessionTTL, err := env.GetSessionTTL(4 * time.Hour)
	if err != nil {
		panic(err)
	}

	cfg := config{
		addr: env.GetAddr(),
		db: dbConfig{
			dsn: env.GetDbString(),
		},
		prefix:       "time",
		appEnv:       env.GetAppEnv(),
		cookieDomain: env.GetCookieDomain(),
		sessionTTL:   sessionTTL,
	}

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "15:04:05",
	})

	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		fmt.Println("error opening database:", err)
		panic("error con database")
	}
	defer db.Close()

	if err := database.RunMigrations(context.Background(), db); err != nil {
		panic("migrations failed: " + err.Error())
	}

	app := application{
		config: cfg,
		db:     db,
		logger: logger,
	}

	app.run(app.mount())
}
