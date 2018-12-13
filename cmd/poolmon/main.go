package main

import (
	"flag"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PgpoolPool struct {
	PoolPid        int    `db:"pool_pid"`
	PoolID         int    `db:"pool_id"`
	BackendID      int    `db:"backend_id"`
	Database       string `db:"database"`
	Username       string `db:"username"`
	CreateTime     string `db:"create_time"`
	StartTime      string `db:"start_time"`
	MajorVersion   int    `db:"majorversion"`
	MinorVersion   int    `db:"minorversion"`
	PoolCounter    int    `db:"pool_counter"`
	PoolBackendPid int    `db:"pool_backendpid"`
	PoolConnected  int    `db:"pool_connected"`
}

func main() {
	connString := flag.String("uri", "postgresql://postgres@127.0.0.1/postgres", "Postgres URI to check for pgpool status")
	flag.Parse()

	db, err := sqlx.Connect("postgres", *connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Queryx("SHOW pool_pools")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var pools []PgpoolPool
	for rows.Next() {
		pool := PgpoolPool{}

		if err := rows.StructScan(&pool); err != nil {
			panic(err)
		}

		pools = append(pools, pool)
	}

	connected := 0
	for _, pool := range pools {
		connected += pool.PoolConnected
	}
	log.Printf("connected:%d\n", connected)
}
