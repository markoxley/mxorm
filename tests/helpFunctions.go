package mxormtests

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/markoxley/mxorm"
	"github.com/markoxley/mxorm/config"
	"github.com/markoxley/mxorm/utils"
)

const (
	dbName     = "mxormtest"
	dbUser     = "root"
	dbPassword = "Dantooine2020!"
	dbPort     = 3306
	dbAddress  = "127.0.0.1"
)

// TestModel for testing database
type TestModel struct {
	mxorm.Model

	Age   int    `mxorm:""`
	Name  string `mxorm:"size:20"`
	Death *int   `mxorm:""`
}

func getConnectionDetails() *config.Config {
	address := fmt.Sprintf("tcp(%s:%d)", dbAddress, dbPort)
	c := config.New(address, dbName, dbUser, dbPassword, true)
	return c
}

func getConnection() (*sql.DB, bool) {
	cs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset-utfs", dbUser, dbPassword, dbAddress, dbPort, dbName)
	if tdb, err := sql.Open("mysql", cs); err == nil {
		return tdb, true
	}
	return nil, false
}

func closeConnection(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func testTableExists(t string) bool {
	if c, ok := getConnection(); ok {
		defer closeConnection(c)
		sql := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM information_schema.TABLES WHERE  TABLE_SCHEMA = 'mxormtest' AND TABLE_NAME = '%s');`, t)
		r, err := c.Query(sql)
		if err == nil {
			if r.Next() {
				return true
			}
		}
		return false
	}
	return false
}

func configuremxorm() {
	mxorm.Configure(getConnectionDetails())
}

func reset() {
	configuremxorm()
	sql := "Delete from TestModel;"
	if c, ok := getConnection(); ok {
		defer closeConnection(c)
		c.Exec(sql)
	}
}

func compareDates(s time.Time, d time.Time) bool {
	d1 := utils.TimeToSQL(s)
	d2 := utils.TimeToSQL(d)

	return d1 == d2
}

func populateTable() {
	names := []string{"Mark", "Sally", "Oliver"}
	ages := []int{52, 57, 22}
	configuremxorm()
	sql := "Drop table if exists TestModel;"
	mxorm.RawExecute(sql)
	for i := range names {
		tm := &TestModel{Name: names[i], Age: ages[i]}
		mxorm.Save(tm)
	}
}
