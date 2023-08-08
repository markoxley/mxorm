package mxormtests

import (
	"fmt"
	"testing"
	"time"

	"github.com/markoxley/mxorm/utils"
	"github.com/markoxley/mxorm/where"
)

func TestWhereBetween(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tm2 := time.Date(2020, 2, 7, 22, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Between Int", input: []interface{}{"amount", 7, 3}, out: "`amount` between 3 and 7"},
		{name: "Between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "`miles` between 12987 and 898989"},
		{name: "Between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "`height` between 13.2300 and 45.4000"},
		{name: "Between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "`weight` between 73.321200 and 983.230000"},
		{name: "Between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("`dob` between '%s' and '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},
	}
	for _, tst := range tests {
		if got := where.Between(tst.input[0].(string), tst.input[1], tst.input[2]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereEqual(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Equal Bool", input: []interface{}{"valid", true}, out: "`valid` = 1"},
		{name: "Equal Int", input: []interface{}{"miles", int(829)}, out: "`miles` = 829"},
		{name: "Equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` = 1003322443"},
		{name: "Equal Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` = 73.1200"},
		{name: "Equal Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` = 432.543300"},
		{name: "Equal String", input: []interface{}{"name", "Sally"}, out: "`name` = 'Sally'"},
		{name: "Equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` = '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.Equal(tst.input[0].(string), tst.input[1]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereGreater(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Greater Int", input: []interface{}{"miles", int(829)}, out: "`miles` > 829"},
		{name: "Greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` > 1003322443"},
		{name: "Greater Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` > 73.1200"},
		{name: "Greater Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` > 432.543300"},
		{name: "Greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` > '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.Greater(tst.input[0].(string), tst.input[1]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereLess(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Less Int", input: []interface{}{"miles", int(829)}, out: "`miles` < 829"},
		{name: "Less Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` < 1003322443"},
		{name: "Less Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` < 73.1200"},
		{name: "Less Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` < 432.543300"},
		{name: "Less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` < '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.Less(tst.input[0].(string), tst.input[1]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereIn(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "In Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "`miles` in (829,21,1)"},
		{name: "In Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "`counters` in (1003322443,437216784)"},
		{name: "In Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "`weight` in (73.1200,1.4300,0.7600,32.2000)"},
		{name: "In Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` in (432.543300)"},
		{name: "In String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "`name` in ('Sally','Mark','Jane','Sam','Jack')"},
		{name: "In DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` in ('%s')", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.In(tst.input[0].(string), tst.input[1:]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereLike(t *testing.T) {
	cs := "`name` like '%ma%'"
	if ts := where.Like("name", "%ma%").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereStartsWith(t *testing.T) {
	cs := "`model` like 'atar%'"
	if ts := where.StartsWith("model", "atar").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereEndsWith(t *testing.T) {
	cs := "`product` like 'ole%'"
	if ts := where.StartsWith("product", "ole").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereContains(t *testing.T) {
	cs := "`breed` like '%ige%'"
	if ts := where.Contains("breed", "ige").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereNotBetween(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tm2 := time.Date(2020, 2, 7, 22, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Not between Int", input: []interface{}{"amount", 7, 3}, out: "`amount` not between 3 and 7"},
		{name: "Not between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "`miles` not between 12987 and 898989"},
		{name: "Not between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "`height` not between 13.2300 and 45.4000"},
		{name: "Not between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "`weight` not between 73.321200 and 983.230000"},
		{name: "Not between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("`dob` not between '%s' and '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},
	}
	for _, tst := range tests {
		if got := where.NotBetween(tst.input[0].(string), tst.input[1], tst.input[2]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotEqual(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Not equal Bool", input: []interface{}{"valid", true}, out: "`valid` <> 1"},
		{name: "Not equal Int", input: []interface{}{"miles", int(829)}, out: "`miles` <> 829"},
		{name: "Not equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` <> 1003322443"},
		{name: "Not equal Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` <> 73.1200"},
		{name: "Not equal Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` <> 432.543300"},
		{name: "Not equal String", input: []interface{}{"name", "Sally"}, out: "`name` <> 'Sally'"},
		{name: "Not equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` <> '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.NotEqual(tst.input[0].(string), tst.input[1]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotGreater(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Not greater Int", input: []interface{}{"miles", int(829)}, out: "`miles` <= 829"},
		{name: "Not greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` <= 1003322443"},
		{name: "Not greater Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` <= 73.1200"},
		{name: "Not greater Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` <= 432.543300"},
		{name: "Not greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` <= '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.NotGreater(tst.input[0].(string), tst.input[1]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotLess(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Not less Int", input: []interface{}{"miles", int(829)}, out: "`miles` >= 829"},
		{name: "Not less Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` >= 1003322443"},
		{name: "Not less Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` >= 73.1200"},
		{name: "Not less Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` >= 432.543300"},
		{name: "Not less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` >= '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.NotLess(tst.input[0].(string), tst.input[1]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotIn(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		name  string
		input []interface{}
		out   string
	}{
		{name: "Not in Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "`miles` not in (829,21,1)"},
		{name: "Not in Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "`counters` not in (1003322443,437216784)"},
		{name: "Not in Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "`weight` not in (73.1200,1.4300,0.7600,32.2000)"},
		{name: "Not in Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` not in (432.543300)"},
		{name: "Not in String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "`name` not in ('Sally','Mark','Jane','Sam','Jack')"},
		{name: "Not in DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` not in ('%s')", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		if got := where.NotIn(tst.input[0].(string), tst.input[1:]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
	for _, tst := range tests {
		if got := where.NotIn(tst.input[0].(string), tst.input[1:]).String(); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotLike(t *testing.T) {
	cs := "`name` not like '%ma%'"
	if ts := where.NotLike("name", "%ma%").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereNotStartsWith(t *testing.T) {
	cs := "`model` not like 'atar%'"
	if ts := where.NotStartsWith("model", "atar").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereNotEndsWith(t *testing.T) {
	cs := "`product` not like 'ole%'"
	if ts := where.NotStartsWith("product", "ole").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereNotContains(t *testing.T) {
	cs := "`breed` not like '%ige%'"
	if ts := where.NotContains("breed", "ige").String(); ts != cs {
		t.Errorf("Expecting %v, got %v", cs, ts)
	}
}

func TestWhereAndConjunction(t *testing.T) {
	e := "`Age` = 12 AND `Name` = 'Alex'"
	if r := where.Equal("Age", 12).AndEqual("Name", "Alex").String(); r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
}

func TestWhereOrConjunction(t *testing.T) {
	e := "`Age` = 12 OR `Name` <> 'Alex'"
	if r := where.Equal("Age", 12).OrNotEqual("Name", "Alex").String(); r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
}

func TestWhereIn2(t *testing.T) {
	expected := []string{
		"`ID` in (1,2,3,4)",
		"`Name` in ('Mark','Sally','Oliver')",
		"`Age` in (42)",
		"`Colour` in ('RED')",
	}
	d := []int{1, 2, 3, 4}
	s := []string{"Mark", "Sally", "Oliver"}
	sd := 42
	ss := "RED"
	result := where.In("ID", d).String()
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = where.In("Name", s).String()
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
	result = where.In("Age", sd).String()
	if result != expected[2] {
		t.Errorf("expecting '%s' got '%s'", expected[2], result)
	}
	result = where.In("Colour", ss).String()
	if result != expected[3] {
		t.Errorf("expecting '%s' got '%s'", expected[3], result)
	}

}

func TestWhereNotIn2(t *testing.T) {
	expected := []string{
		"`ID` not in (1,2,3,4)",
		"`Name` not in ('Mark','Sally','Oliver')",
	}
	d := []int{1, 2, 3, 4}
	s := []string{"Mark", "Sally", "Oliver"}
	result := where.NotIn("ID", d).String()
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = where.NotIn("Name", s).String()
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
}

func TestWhereAndIn2(t *testing.T) {
	expected := []string{
		"`ID` = 2 AND `Size` in (2,4,6)",
		"`ID` = 3 AND `Name` in ('Mark','Sally','Oliver')",
	}
	d := []int{2, 4, 6}
	s := []string{"Mark", "Sally", "Oliver"}
	result := where.Equal("ID", 2).AndIn("Size", d).String()
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = where.Equal("ID", 3).AndIn("Name", s).String()
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
}
