package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	//arango "github.com/joselitofilho/gorm-arango/pkg"
	// "github.com/arangodb/go-driver/http"
	// driver "github.com/arangodb/go-driver"

	//arango "github.com/arangodb/go-driver"
	// "github.com/arangodb/go-driver/http"
	"fmt"
	// "github.com/fatih/structs"
	// "strings"
	// "os"
	"encoding/json"
	// "reflect"
	"database/sql"
	"time"
)

var ActiveUser string

// var db *mysql.
var OpenDB *gorm.DB

const (
	Username string = "pusdikl7_jawievent"
	Password string = "Wonderwoman122.."
	DbName   string = "pusdikl7_jawievent"
	Tcp      string = "203.161.184.81:3306"
)

func DbConnect() (*gorm.DB, error) {
	// Capture connection properties.
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v", Username, Password, Tcp, DbName)

	// dbMaster, err := gorm.Open(cfg.FormatDSN())
	// Get a database handle.
	// connSql, _ := sql.Open("mysql", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil

}

// parse null string on model
type NullString struct {
	sql.NullString
}

// parse null string on model
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte(`null`), nil
}

func NullStringInput(s string) NullString {
	if len(s) == 0 {
		return NullString{sql.NullString{`null`, false}}
	}
	return NullString{sql.NullString{s, true}}
}

// parse null int on model
type NullInt64 struct {
	sql.NullInt64
}

// parse null int on model
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return []byte(`null`), nil
}

func NullInt64Input(s int64) NullInt64 {
	if s == 0 {
		return NullInt64{sql.NullInt64{0, false}}
	}
	return NullInt64{sql.NullInt64{s, true}}
}

// parse null time on model
type NullDateTime struct {
	sql.NullTime
}

// parse null time on model
func (nt NullDateTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		t := nt.Time
		return json.Marshal(t.Format("2006-01-02 15:04:05"))
	}
	return []byte(`null`), nil
}

func NullDateTimeInput(t time.Time) NullDateTime {
	if t.IsZero() {
		return NullDateTime{sql.NullTime{time.Time{}, false}}
	}
	return NullDateTime{sql.NullTime{t, true}}
}

// date
type NullDate struct {
	sql.NullTime
}

func (nt NullDate) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		t := nt.Time
		return json.Marshal(t.Format("2006-01-02"))
	}
	return []byte(`null`), nil
}

func NullDateInput(t time.Time) NullDate {
	if t.IsZero() {
		return NullDate{sql.NullTime{time.Time{}, false}}
	}
	return NullDate{sql.NullTime{t, true}}
}

// get Columns Name from Model
func GetColumnName(db *gorm.DB, model []interface{}) *gorm.DB {
	for _, v := range model {
		fmt.Println(v)
	}
	return db
}
