package mxorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/markoxley/mxorm/config"
	"github.com/markoxley/mxorm/utils"
)

const (
	majorVersion   = 1
	minorVersion   = 0
	releaseVersion = 0
)

type fieldMap struct {
	name      string
	fieldType string
}

var (
	conf        *config.Config
	configured  bool
	knownTables []string
)

const (
	mySQLConnectionPattern = "%s:%s@%s/%s"
)

func init() {
	knownTables = make([]string, 0, 20)
}

// Configure attempts to configure the database
// @param c This is the configuration
// @return error
func Configure(c *config.Config) error {
	conf = c
	db, err := connect()
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	db.Close()
	configured = true
	return nil
}

func Configured() bool {
	return configured
}

// connect attempts to connect to the database
// @return *sql.DB
// @return error
func connect() (*sql.DB, error) {
	cs := fmt.Sprintf(mySQLConnectionPattern, conf.User, conf.Password, conf.Host, conf.Name)
	tdb, err := sql.Open("mysql", cs)
	if err != nil {
		return nil, err
	}
	return tdb, nil
}

// disconnect from the database
// @param db
func disconnect(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// beginTransaction begins the transaction process
// @param db
// @return *sql.Tx
// @return error
func beginTransaction(db *sql.DB) (*sql.Tx, error) {
	return db.Begin()
}

// commitTransaction commites the transaction to the database
// @param tx
func commitTransaction(tx *sql.Tx) {
	if tx != nil {
		tx.Commit()
	}
}

// selectScalar atempts to execute the specified query and returns
// the value of the first column of the first row
// @param q
// @return interface{}
// @return bool
func selectScalar(q string) (interface{}, bool) {
	db, err := connect()
	if err != nil {
		return nil, false
	}
	defer disconnect(db)

	tx, err := beginTransaction(db)
	if err != nil {
		return nil, false
	}
	defer commitTransaction(tx)

	res, err := db.Query(q)
	if err != nil {
		return nil, false
	}
	defer res.Close()
	if res.Next() {
		var cols string
		vl := &cols
		//var vl interface{}
		res.Scan(vl)
		return cols, true
	}
	return nil, false

}

// selectQuery attempts to execute the query passed, returning
// a slice of the type specified by the type parameter
// @param q
// @return []*T
// @return bool
func selectQuery[T Modeller](q string) ([]*T, bool) {
	db, err := connect()
	if err != nil {
		return nil, false
	}
	defer disconnect(db)

	tx, err := beginTransaction(db)
	if err != nil {
		return nil, false
	}
	defer commitTransaction(tx)

	res, err := db.Query(q)
	if err != nil {
		return nil, false
	}
	defer res.Close()
	return populateModel[T](res)
}

// populateModel creates a new slice of models of the type
// specified by the type parameter and populates the fields from the sql query
// @param r
// @return []*T
// @return bool
func populateModel[T Modeller](r *sql.Rows) ([]*T, bool) {
	res := make([]*T, 0, 100)
	ok := true
	m := *new(T)
	// Get the column count
	cc, _ := r.Columns()

	// Make them all uppercase
	for i := range cc {
		cc[i] = strings.ToUpper(cc[i])
	}

	// Get the fields of the model and build a map of them
	//t := reflect.TypeOf(*m)
	flds, ok := tableDef[getTableName(m)]
	if !ok {
		return nil, false
	}
	fMap := make(map[string]field, len(flds))
	for _, f := range flds {
		fMap[strings.ToUpper(f.name)] = f
	}

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))
	s := reflect.TypeOf(m)
	rowCount := 0
	for r.Next() {
		// s := reflect.TypeOf(m)
		v := reflect.New(s)

		for i := range cols {
			vls[i] = &cols[i]
		}
		r.Scan(vls...)

		for i := 0; i < len(cc); i++ {
			if cols[i] == nil {
				continue
			}
			if cc[i] == "ID" {
				tmpID := cols[i]
				v.Elem().FieldByName("ID").Set(reflect.ValueOf(tmpID))
			} else if cc[i] == "CREATEDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpCreate := *tm
						v.Elem().FieldByName("CreateDate").Set(reflect.ValueOf(tmpCreate))
					}
				}
			} else if cc[i] == "LASTUPDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpUpdate := *tm
						v.Elem().FieldByName("LastUpdate").Set(reflect.ValueOf(tmpUpdate))
					}
				}
			} else if cc[i] == "DELETEDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpDeleted := tm
						v.Elem().FieldByName("DeleteDate").Set(reflect.ValueOf(tmpDeleted))
					}
				}
			} else if fld, ok := fMap[cc[i]]; ok {
				switch fld.fType {
				case tInt, tLong:
					if fld.unsigned {
						if val, err := strconv.ParseUint(*cols[i], 10, 64); err != nil {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).SetUint(val)
							}
						}
					} else {
						if val, err := strconv.ParseInt(*cols[i], 10, 64); err == nil {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).SetInt(val)
							}
						}
					}
				case tBool:
					if val, err := strconv.ParseInt(*cols[i], 10, 64); err == nil {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Elem().SetBool(val == 1)
						} else {
							v.Elem().FieldByName(fld.name).SetBool(val == 1)
						}
					}
				case tDecimal, tFloat, tDouble:
					if val, err := strconv.ParseFloat(*cols[i], 64); err == nil {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
						} else {
							v.Elem().FieldByName(fld.name).SetFloat(val)
						}
					}
				case tDateTime:
					if cols[i] != nil {
						if val, ok := utils.SQLToTime(*cols[i]); ok {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(*val))
							}
						}
					}
				case tChar:
					st := *cols[i]
					strVal := st[:1]
					if fld.allowNull {
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(strVal)
					}
				case tString:
					if fld.allowNull {
						strVal := *cols[i]
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(*cols[i])
					}
				case tUUID:
					if fld.allowNull {
						strVal := *cols[i]
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(*cols[i])
					}
				}
			}

		}
		newObj := interface{}(v.Elem().Interface()).(T)
		res = append(res, &newObj)
		rowCount++

	}
	return res[:rowCount], ok
}

// updateModel updates the date fields of the specified model
// @param m
// @param id
// @param createdate
// @param updatedate
// @param deletedate
func updateModel(m Modeller, id string, createdate time.Time, updatedate time.Time, deletedate *time.Time) {
	v := reflect.ValueOf(m)
	createdateValue := reflect.ValueOf(createdate)
	updatedateValue := reflect.ValueOf(updatedate)
	deletedateValue := reflect.ValueOf(deletedate)
	rv := reflect.New(reflect.TypeOf(id))
	rv.Elem().Set(reflect.ValueOf(id))

	v.Elem().FieldByName("ID").Set(rv)
	v.Elem().FieldByName("CreateDate").Set(createdateValue)
	v.Elem().FieldByName("LastUpdate").Set(updatedateValue)
	v.Elem().FieldByName("DeleteDate").Set(deletedateValue)
}

// executeQuery attempts to execute the passed sql query
// @param q
// @return bool
func executeQuery(q string) bool {
	db, err := connect()
	if err != nil {
		return false
	}
	defer disconnect(db)

	tx, err := beginTransaction(db)
	if err != nil {
		return false
	}
	defer commitTransaction(tx)

	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// tableExists tests for the existence of the specified table
// @param t
// @return bool
func tableExists(t string) bool {
	for _, tn := range knownTables {
		if tn == t {
			return true
		}
	}
	if _, ok := selectScalar(fmt.Sprintf("SHOW TABLES WHERE Tables_in_%s = '%s'", conf.Name, t)); ok {
		knownTables = append(knownTables, t)
		return true
	}
	return false
}

// RawExecute executes a sql statement on the database, without returning a value
// Not recommended for general use - can break shadowing
// @param sql
// @return bool
func RawExecute(sql string) bool {
	return executeQuery(sql)
}

// RawScalar exeutes a raw sql statement that returns a single value
// Not recommended for general use
// @param sql
// @return interface{}
// @return bool
func RawScalar(sql string) (interface{}, bool) {
	return selectScalar(sql)
}

// RawSelect executes a raw sql statement on the database
// Not recommended for general use
// @param sql
// @return []map
func RawSelect(sql string) []map[string]interface{} {
	db, err := connect()
	if err != nil {
		return nil
	}
	defer disconnect(db)
	res, err := db.Query(sql)
	if err != nil {
		return nil
	}
	defer res.Close()
	data := make([]map[string]interface{}, 0, 10)

	// Get the column count
	cc, _ := res.Columns()

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))

	for res.Next() {

		for i := range cols {
			vls[i] = &cols[i]
		}
		res.Scan(vls...)
		row := make(map[string]interface{})
		for i, n := range cc {
			row[n] = vls[i]
		}
		data = append(data, row)
	}
	return data
}

// getCriteria returns the criteria for a query in SQL format
// @param criteria
// @return *Criteria
// @return error
func getCriteria(criteria []interface{}) (*Criteria, error) {
	for _, cr := range criteria {
		if cr == nil {
			continue
		}

		if c, ok := cr.(*Criteria); ok {
			return c, nil
		} else if c, ok := cr.(Criteria); ok {
			return &c, nil
		} else if c, ok := cr.(fmt.Stringer); ok {
			return &Criteria{Where: c}, nil
		} else if c, ok := cr.(string); ok {
			return &Criteria{Where: c}, nil
		}
		return nil, errors.New("invalid criteria format")
	}
	return &Criteria{}, nil
}

// Fetch populates the slice with models from the database that match the criteria.
// Returns an error if this fails
// @param criteria
// @return []*T
// @return error
func Fetch[T Modeller](criteria ...interface{}) ([]*T, error) {
	c, err := getCriteria(criteria)
	if err != nil {
		return nil, err
	}
	m := *new(T)
	t := reflect.TypeOf(m)
	n := t.Name()
	_, n, ok := tableTest(m)
	if !ok {
		return nil, errors.New("failed table check")
	}
	s := fmt.Sprintf("select * from `%s`", n)
	s += c.String()
	res, ok := selectQuery[T](s)
	if !ok {
		return nil, errors.New("error selecting data")
	}
	return res, nil
}

// First returns the first model that matches the criteria
// @param criteria
// @return *T
// @return error
func First[T Modeller](criteria ...interface{}) (*T, error) {
	c, err := getCriteria(criteria)
	if err != nil {
		return nil, err
	}
	c.Limit = 1
	c.Offset = 0
	ml, err := Fetch[T](c)
	if err != nil {
		return nil, err
	}
	if len(ml) > 0 {
		return ml[0], nil
	}

	return nil, nil
}

// Count returns the number of rows in the database that match the criteria
// @param criteria
// @return int
func Count[T Modeller](criteria ...interface{}) int {
	c, err := getCriteria(criteria)
	if err != nil {
		return -1
	}
	m := *new(T)
	t := getTableName(m)
	if !tableExists(t) {
		return 0
	}
	s := fmt.Sprintf("Select Count(*) from `%s`", t)
	s += c.WhereString()
	if i, ok := selectScalar(s); ok {
		if vl, vlok := i.(string); vlok {
			if res, err := strconv.Atoi(vl); err == nil {
				return res
			}
		}
	}
	return 0

}

// Save stores the model in the database.
// Depending on the status of the model, this is either
// an update or an insert command
// @param m
// @return bool
func Save(m Modeller) bool {
	if m.IsNew() {
		return executeQuery(insertCommand(m))
	}
	return executeQuery(updateCommand(m))
}

// Remove removes the passed model from the database
// @param m
// @return bool
func Remove(m Modeller) bool {
	if m.GetID() == nil {
		return true
	}
	if conf.Deletable {
		return executeQuery(fmt.Sprintf("delete from `%s` where id = '%s'", getTableName(m), *(m.GetID())))
	}
	now := time.Now()
	return executeQuery(fmt.Sprintf("update `%s` set `deleteDate` = %v where `id` = '%s'", getTableName(m), utils.TimeToSQL(now), *(m.GetID())))
}

// RemoveMany removes all models of the specified type that match the criteria
// @param c
// @return int
// @return bool
func RemoveMany[T Modeller](c *Criteria) (int, bool) {
	t := getTableName(*new(T))
	if !tableExists(t) {
		return 0, true
	}
	r := Count[T](c)
	if r == 0 {
		return 0, true
	}
	s := ""
	if conf.Deletable {
		s = fmt.Sprintf("delete from %s", t)
	} else {
		tm := time.Now()
		s = fmt.Sprintf("update %s set `deleteDate` = '%v'", t, utils.TimeToSQL(tm))
	}
	whereAdded := false
	if c != nil && c.Where != "" {
		s += fmt.Sprintf(" where %s", c.Where)
		whereAdded = true
	}

	if whereAdded {
		s += " AND DeleteDate is null"
	} else {
		s += " WHERE DeleteDate is null"
	}
	ok := executeQuery(s)
	return r, ok
}

func Version() string {
	return fmt.Sprintf("Batty mxorm v%d.%d.%d ©2023 I Have a Hat", majorVersion, minorVersion, releaseVersion)
}
