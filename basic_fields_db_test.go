package core_test

import (
	"database/sql"
	core "github.com/scys/service-core-go"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TestDBGetStructA struct {
	core.BasicFields

	A string    `json:"a"`
	B int       `json:"b_field"`
	C time.Time `json:"c"`
	D bool      `json:"d"`
	E core.H    `json:"e"`
}

func (t TestDBGetStructA) TableName() string {
	return "test_struct_a"
}

func __testDBPrepare() *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	// #1
	db.Exec(`CREATE TABLE IF NOT EXISTS test_struct_a (
    	id TEXT primary key, 
    	ts_create time not null default current_timestamp, 
    	ts_update time not null default current_timestamp, 
    	removed bool not null default false,
    	info TEXT not null default '{}',
    	a TEXT, 
    	b_field INTEGER, 
    	c TEXT, 
    	d INTEGER not null DEFAULT 0,  
    	e TEXT
	);`)

	return db
}

func TestDBInsert(t *testing.T) {
	db := __testDBPrepare()
	defer db.Close()

	type args struct {
		db   *sql.DB
		item core.BasicFieldsInterface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"simple usage",
			args{
				db,
				&TestDBGetStructA{
					BasicFields: core.NewBasicFields(),
					A:           "a",
					B:           1,
					C:           time.Now(),
					D:           true,
					E:           core.H{"a": "b", "c": 123},
				},
			},
			false,
		},
		{
			"simple usage without pointer",
			args{
				db,
				TestDBGetStructA{
					BasicFields: core.NewBasicFields(),
					A:           "a_field",
					B:           2,
					C:           time.Now(),
					D:           true,
					E:           core.H{"a2": "b2", "c2": 321},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := core.DBInsert(tt.args.db, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("DBInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestDBCount(t *testing.T) {
//	type args struct {
//		db   *sql.DB
//		item core.BasicFieldsInterface
//		raw  string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    int
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := core.DBCount(tt.args.db, tt.args.item, tt.args.raw)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("DBCount() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("DBCount() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestDBFilter(t *testing.T) {
//	type args struct {
//		db          *sql.DB
//		item        core.BasicFieldsInterface
//		raw         string
//		order       string
//		offset      int64
//		limit       int64
//		scanWrapper func(*sql.Rows) error
//		params      []any
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := DBFilter(tt.args.db, tt.args.item, tt.args.raw, tt.args.order, tt.args.offset, tt.args.limit, tt.args.scanWrapper, tt.args.params...); (err != nil) != tt.wantErr {
//				t.Errorf("DBFilter() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestDBGet(t *testing.T) {
//	type args struct {
//		db    *sql.DB
//		item  core.BasicFieldsInterface
//		raw   string
//		order string
//		key   any
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := core.DBGet(tt.args.db, tt.args.item, tt.args.raw, tt.args.order, tt.args.key); (err != nil) != tt.wantErr {
//				t.Errorf("DBGet() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

//func TestDBRemove(t *testing.T) {
//	type args struct {
//		db   *sql.DB
//		item BasicFieldsInterface
//		key  any
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := DBRemove(tt.args.db, tt.args.item, tt.args.key); (err != nil) != tt.wantErr {
//				t.Errorf("DBRemove() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestDBUpdate(t *testing.T) {
//	type args struct {
//		db   *sql.DB
//		item BasicFieldsInterface
//		key  any
//		data H
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := DBUpdate(tt.args.db, tt.args.item, tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
//				t.Errorf("DBUpdate() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
