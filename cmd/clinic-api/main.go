package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model/filler"
	"github.com/Zhassulan1/Go_Project/pkg/jsonlog"

	// "github.com/Zhassulan1/Go_Project/pkg/vcs"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/peterbourgon/ff/v3"

	_ "github.com/lib/pq"
)

// var (
// 	version = vcs.Version()
// )

type config struct {
	port       string
	env        string
	fill       bool
	migrations string
	db         struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	wg     sync.WaitGroup
	logger *jsonlog.Logger
}

func main() {
	// <<<<<<< docking
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	// var cfg config
	// flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	// flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1234@localhost/medicalclinic?sslmode=disable", "PostgreSQL DSN")
	// flag.Parse()
	var (
		cfg        config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.String("port", ":8060", "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://postgres:1234@localhost:5432/medicalclinic?sslmode=disable", "PostgreSQL DSN")
	)
	// =======
	// 	var cfg config
	// 	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	// 	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// 	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:Yekanai11@localhost/medicalclinic?sslmode=disable", "PostgreSQL DSN")
	// 	flag.Parse()
	// >>>>>>> main

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	cfg.port = *port
	cfg.env = *env
	cfg.fill = *fill
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	logger.PrintInfo("starting application with configuration", map[string]string{
		"port":       fmt.Sprintf(cfg.port),
		"fill":       fmt.Sprintf("%t", cfg.fill),
		"env":        cfg.env,
		"db":         cfg.db.dsn,
		"migrations": cfg.migrations,
	})

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	// Defer a call to db.Close() so that the connection pool is closed before the main()
	// function exits.
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}

	if cfg.fill {
		err = filler.PopulateDatabase(app.models)
		if err != nil {
			logger.PrintFatal(err, nil)
			return
		}
	}

	// Call app.server() to start the server.
	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Check if the connection is working by executing a test query
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	if cfg.migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.migrations,
			"postgres", driver)
		if err != nil {
			return nil, err
		}
		m.Up()
	}

	return db, nil
}
