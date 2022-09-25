package db

import (
	"database/sql"

	"github.com/xxxsen/common/database"
	"github.com/xxxsen/common/errs"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbFileInfo *sql.DB
)

func InitFileDB(c *database.DBConfig) error {
	client, err := database.InitDatabase(c)
	if err != nil {
		return errs.Wrap(errs.ErrDatabase, "open db fail", err)
	}
	dbFileInfo = client
	return nil
}

func GetFileDB() *sql.DB {
	return dbFileInfo
}
