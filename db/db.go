package db

import (
	"database/sql"
	"fmt"

	"github.com/xxxsen/common/database"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbFileInfo *sql.DB
)

func InitFileDB(c *database.DBConfig) error {
	client, err := database.InitDatabase(c)
	if err != nil {
		return fmt.Errorf("open db fail, err:%w", err)
	}
	dbFileInfo = client
	return nil
}

func GetFileDB() *sql.DB {
	return dbFileInfo
}
