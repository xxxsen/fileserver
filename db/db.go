package db

import (
	"database/sql"

	"github.com/xxxsen/common/database"
	"github.com/xxxsen/common/errs"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbMediaInfo *sql.DB
)

func InitFileDB(c *database.DBConfig) error {
	client, err := database.InitDatabase(c)
	if err != nil {
		return errs.Wrap(errs.ErrDatabase, "open db fail", err)
	}
	dbMediaInfo = client
	return nil
}

func GetMediaDB() *sql.DB {
	return dbMediaInfo
}
