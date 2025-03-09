package db

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/xxxsen/common/database/kv"
	"github.com/xxxsen/common/database/kv/bolt"
)

var (
	dbClient kv.IKvDataBase
)

var (
	tableList = []string{
		"tg_file_tab",
		"tg_file_part_tab",
		"tg_file_part_tab",
	}
)

func InitDB(file string) error {
	db, err := bolt.New(file, tableList...)
	if err != nil {
		return err
	}
	dbClient = db
	return nil
}

func GetClient() kv.IKvDataBase {
	return dbClient
}
