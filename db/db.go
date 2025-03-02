package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

type initSql struct {
	name string
	sql  string
}

var (
	dbClient *sql.DB
)

var initList = []initSql{
	{
		name: "create_file_info_tab",
		sql: `CREATE TABLE IF NOT EXISTS file_info_tab (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_name TEXT NOT NULL,
    hash TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    create_time INTEGER NOT NULL,
    down_key INTEGER NOT NULL,
    file_key TEXT NOT NULL,
    extra BLOB NOT NULL,
    st_type INTEGER NOT NULL,
    UNIQUE(down_key)
);`,
	},
	{
		name: "create_mapping_info_tab",
		sql: `CREATE TABLE IF NOT EXISTS mapping_info_tab (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_name TEXT NOT NULL,
    hash_code INTEGER NOT NULL,
    check_sum TEXT NOT NULL,
    create_time INTEGER NOT NULL,
    modify_time INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    UNIQUE(check_sum)
);`,
	},
}

func InitDB(file string) error {
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return err
	}
	if err := tableInit(db); err != nil {
		return err
	}
	dbClient = db
	return nil
}

func tableInit(db *sql.DB) error {
	ctx := context.Background()
	for _, initItem := range initList {
		if _, err := db.ExecContext(ctx, initItem.sql); err != nil {
			return fmt.Errorf("init sql failed, name:%s, err:%w", initItem.name, err)
		}
	}
	return nil
}

func GetClient() *sql.DB {
	return dbClient
}
