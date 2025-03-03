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
		sql: `CREATE TABLE IF NOT EXISTS tg_file_tab (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 自增 ID
    file_id INTEGER UNIQUE NOT NULL,         -- 文件 ID, 唯一键
    file_name TEXT NOT NULL,              -- 文件名
    file_size INTEGER NOT NULL,           -- 文件大小
    file_part_count INTEGER NOT NULL,     -- 文件分片数量
    ctime INTEGER NOT NULL, -- 创建时间
    mtime INTEGER NOT NULL, -- 修改时间
    file_state INTEGER NOT NULL           -- 文件状态
);`,
	},
	{
		name: "create_file_part_tab",
		sql: `CREATE TABLE IF NOT EXISTS tg_file_part_tab (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- 自增 ID
    file_id INTEGER UNIQUE NOT NULL,       -- 文件 ID(64 位整数)，唯一键
    file_key TEXT NOT NULL,                -- 文件 Key
    file_part_id INTEGER NOT NULL,         -- 文件分片 ID
    ctime INTEGER NOT NULL,                -- 创建时间（存 UNIX 时间戳）
    mtime INTEGER NOT NULL                 -- 修改时间（存 UNIX 时间戳）
);`,
	},
	{
		name: "create_file_mapping_tab",
		sql: `
CREATE TABLE IF NOT EXISTS tg_file_mapping_tab (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- 自增 ID
    file_name TEXT NOT NULL,               -- 文件名
	file_hash TEXT UNIQUE NOT NULL,        -- 文件名hash
    file_id INTEGER NOT NULL,              -- 文件 ID,唯一键
    ctime INTEGER NOT NULL,                -- 创建时间（存 UNIX 时间戳）
    mtime INTEGER NOT NULL                 -- 修改时间（存 UNIX 时间戳）
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
