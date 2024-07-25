package mxormtests

import (
	"testing"

	"github.com/markoxley/mxorm"
)

type Updateable struct {
	mxorm.Model
	Name string `mxorm:"size:64"`
	Size int    `mxorm:""`
}

func (u *Updateable) Update() {
	u.Size = len(u.Name)
}

func createUpdateableTest() {
	names := []string{"Mark", "Sally", "Oliver"}
	configuremxorm()
	sql := "Drop table if exists Updateable;"
	mxorm.RawExecute(sql)
	for i := range names {
		tm := &Updateable{Name: names[i]}
		mxorm.Save(tm)
	}
}

func TestUpdate(t *testing.T) {
	createUpdateableTest()
	recs, err := mxorm.Fetch[Updateable](nil)
	if err != nil {
		t.Errorf("unable to read records: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("expected 3 records, got %d", len(recs))
	}
	for _, r := range recs {
		if r.Size != len(r.Name) {
			t.Errorf("expected size %d, got %d", len(r.Name), r.Size)
		}
	}
}
