package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// configuration settings
type config struct {
	addr string
	dsn string
}

// application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	var cfg config

	// command-line flags
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse();

	// initialize structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// open db connection pool
	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// initialize instance of application with our dependencies
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", cfg.addr))

	// listens on the passed TCP network address with our servemux
	err = http.ListenAndServe("localhost" + cfg.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// returns an sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}