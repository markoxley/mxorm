package mxormtests

import (
	"testing"

	"github.com/markoxley/mxorm"
)

func TestRawSelect(t *testing.T) {
	names := []string{"Oliver", "Mark", "Sally"}
	populateTable()
	sql := "SELECT * FROM `testmodel` ORDER By Age ASC"
	rows := mxorm.RawSelect(sql)
	if rows == nil {
		t.Errorf("RawSelect() returned nil")
	}
	if len(rows) != len(names) {
		t.Errorf("RawSelect() returned %d rows, expected %d", len(rows), len(names))
	}

	for i, row := range rows {
		rowName := *(row["Name"].(*string))
		if rowName != names[i] {
			t.Errorf("RawSelect() row %d returned Name = %s, expected %s", i, row["Name"], names[i])
		}
	}

}
