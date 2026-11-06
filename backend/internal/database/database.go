package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	//1. parse the conn string
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	//2. Tune the parameters for performance
	poolConfig.MaxConns = 25                      //Max conns in the pool
	poolConfig.MinConns = 5                       //Keep at least 5 warm conns
	poolConfig.MaxConnLifetime = 1 * time.Hour    //Retire connections after 1hr
	poolConfig.MaxConnIdleTime = 30 * time.Minute //Close idle connects after 30mins

	//3 Create context with a 5-sec timeout for establising the
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//4. Create the pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create conn pool:%w", err)
	}
	//5. Ping the DB to veify creds and connectivity
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to create conn pool: %w", err)
	}
	log.Printf("Successfully connected to postgres")
	return pool, nil
}
