package mxorm

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/markoxley/mxorm/utils"
	uuid "github.com/satori/go.uuid"
)

// ModelState is the state of the current model
type ModelState int

// Model is the base for all database models
type Model struct {
	ID         *string
	CreateDate time.Time
	LastUpdate time.Time
	DeleteDate *time.Time
	tableName  *string
}

// CreateModel sets the default parameters for the Model
// @return Model
func CreateModel() Model {
	return Model{
		CreateDate: time.Now(),
		LastUpdate: time.Now(),
	}
}

// StandingData returns the standing data for when the table is created
// @receiver m
// @return []Modeller
func (m Model) StandingData() []Modeller {
	return nil
}

// GetID returns the ID of the model
// @receiver m
// @return *string
func (m Model) GetID() *string {
	return m.ID
}

// IsNew returns true if the model has yet to be stored
// @receiver m
// @return bool
func (m Model) IsNew() bool {
	return m.ID == nil
}

// IsDeleted returns true if teh model has been marked as deleted
// @receiver m
// @return bool
func (m Model) IsDeleted() bool {
	return m.DeleteDate == nil
}

// getTableName getTableNAme returns the name of the table based on the Model specified
// @param m
// @return string
func getTableName(m Modeller) string {
	if reflect.TypeOf(m).Kind() == reflect.Pointer {
		return reflect.Indirect(reflect.ValueOf(m).Elem()).Type().Name()
	}
	return reflect.ValueOf(m).Type().Name()

}

// tableTest tests for the existence of the specified table
// If the table does not exist, it is created
// @param m
// @return []field
// @return string
// @return bool
func tableTest(m Modeller) ([]field, string, bool) {
	sql, required := tableDefinition(m)
	if required {
		te := tableExists(getTableName(m))
		knownTables = append(knownTables, getTableName(m))
		if !te {
			for _, s := range sql {
				if !executeQuery(s) {
					log.Panicf(`Error executing "%s"`, s)
				}
			}
			if standingData := m.StandingData(); standingData != nil {
				for _, data := range standingData {
					Save(data)
				}
			}
		}
	}
	flds, ok := tableDef[getTableName(m)]
	return flds, getTableName(m), ok
}

// tableDefinition Returns a slice of strings with the sql statements and boolean to
// indicate if the table needs to be created
// @param m
// @return []string
// @return bool
func tableDefinition(m Modeller) ([]string, bool) {
	sql := make([]string, 0, 3)

	n := getTableName(m)
	if _, ok := tableDef[n]; ok {
		return nil, false
	}

	t := reflect.TypeOf(m)
	var nm interface{}
	if t.Kind() == reflect.Ptr {
		nm = reflect.New(t.Elem()).Elem().Interface()
	} else {
		nm = reflect.New(t).Elem().Interface()
	}
	fs := getDefs(nm, true)

	tableDef[n] = fs
	if len(fs) == 0 {
		return nil, false
	}
	flds := ""
	keys := make([]string, 0, 5)
	for _, f := range fs {
		if flds != "" {
			flds += ", "
		}
		flds += fmt.Sprintf("`%s` %s", f.name, pgFieldNames[f.fType])
		if f.fType != tUUID && f.fType != tChar && f.size.size > 0 {
			flds += fmt.Sprintf("(%s)", f.size.String())
		}
		if f.fType == tString && f.size.size == 0 {
			flds += "(256)"
		}
		if f.unsigned {
			flds += " UNSIGNED"
		}
		if !f.allowNull {
			flds += " NOT NULL"
		}
		if f.key {
			keys = append(keys, f.name)
		}
	}
	sql = append(sql, fmt.Sprintf(sqlTableCreate, n, flds))
	kn := strings.ReplaceAll(n, ".", "_")
	for _, k := range keys {
		sql = append(sql, fmt.Sprintf(sqlIndexCreate, kn, k, n, k))
	}
	return sql, true
}

// insertCommand returns the sql command to insert the current model into the database
// @param m
// @return string
func insertCommand(m Modeller) string {
	flds, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	uid := uuid.NewV4()

	fds := "ID, CreateDate, LastUpdate"
	now := time.Now()
	dbNow := utils.TimeToSQL(now)
	updateModel(m, fmt.Sprintf("%s", uid), now, now, nil)
	q := fmt.Sprintf("'%s', '%s', '%s'", uid, dbNow, dbNow)
	v := reflect.ValueOf(m).Elem()
	for _, f := range flds {
		if f.name == "ID" || f.name == "CreateDate" || f.name == "LastUpdate" || f.name == "DeleteDate" {
			continue
		}
		vi := v.FieldByName(f.name)

		if f.allowNull {
			if vi.IsNil() {
				continue
			}
			vi = vi.Elem()
		}

		vf := vi.Interface()

		if vl, ok := utils.MakeValue(vf); ok {
			fds += fmt.Sprintf(", `%s`", f.name)
			q += fmt.Sprintf(", %s", vl)
		}
	}

	def := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", n, fds, q)
	return def
}

// updateCommand returns the SQL command to update the
// current model in the database
// @param m
// @return string
func updateCommand(m Modeller) string {
	flds, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	now := time.Now()
	updateLastUpdate(m, now)
	res := fmt.Sprintf("UPDATE %s SET", n)
	v := reflect.ValueOf(m)
	first := true
	for _, f := range flds {
		if f.name != "ID" && f.name != "CreateDate" {
			if !first {
				res += ","
			}
			first = false
			var value interface{}
			if f.allowNull {
				if v.Elem().FieldByName(f.name).IsNil() {
					res += fmt.Sprintf(" `%s` = null", f.name)
					continue
				}
				value = v.Elem().FieldByName(f.name).Elem().Interface()
			} else {
				value = v.Elem().FieldByName(f.name).Interface()
			}
			if vl, ok := utils.MakeValue(value); ok {
				res += fmt.Sprintf(" `%s` = %s", f.name, vl)
			}
		}
	}
	def := res + fmt.Sprintf(" WHERE `Id` = '%s'", *m.GetID())
	return def
}

// deleteCommand returns the SQL command to remove the model from the database
// @param m
// @return string
func deleteCommand(m Modeller) string {
	_, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	def := fmt.Sprintf("DELETE FROM %s WHERE `Id` = '%s'", n, *m.GetID())
	return def
}

// refreshCommand returns the SQL query to refresh the data in the model
// @param m
// @return string
func refreshCommand(m Modeller) string {
	_, n, ok := tableTest(m)
	if !ok {
		return ""
	}
	def := fmt.Sprintf("SELECT * FROM %s WHERE `Id` = '%s'", n, *m.GetID())
	return def
}

// updateLastUpdate updates the LastUpdate field in the model
// @param m
// @param date
func updateLastUpdate(m Modeller, date time.Time) {
	v := reflect.ValueOf(m)
	dateValue := reflect.ValueOf(date)
	v.Elem().FieldByName("LastUpdate").Set(dateValue)
}
