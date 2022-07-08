package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	RecordCount int = 100000
	LoadCount   int = 100000
)

var db *sql.DB
var text string

var FooMap = map[string]Foo{}

type Foo struct {
	ID     string
	Field1 string
	Field2 string
	Field3 string
}

func main() {
	log.Println("setup...")
	if err := setup(); err != nil {
		panic(err)
	}

	startMeasure("load from heap memory")
	if err := loadFromHeapMemory(); err != nil {
		panic(err)
	}
	endMeasure()

	startMeasure("load from sqlite memory ( not scan )")
	if err := loadFromSQLiteMemoryNotScan(); err != nil {
		panic(err)
	}
	endMeasure()

	startMeasure("load from sqlite memory ( scan )")
	if err := loadFromSQLiteMemoryScan(); err != nil {
		panic(err)
	}
	endMeasure()
}

func loadFromHeapMemory() error {
	for i := 1; i <= LoadCount; i++ {
		id := strconv.Itoa(i)
		_ = FooMap[id]
	}
	return nil
}

func loadFromSQLiteMemoryNotScan() error {
	query := "SELECT * FROM `Foo` WHERE `ID` = ?"
	for i := 1; i <= LoadCount; i++ {
		id := strconv.Itoa(i)
		_ = db.QueryRow(query, id)
	}
	return nil
}

func loadFromSQLiteMemoryScan() error {
	query := "SELECT * FROM `Foo` WHERE `ID` = ?"
	for i := 1; i <= LoadCount; i++ {
		id := strconv.Itoa(i)
		var foo Foo
		_ = db.QueryRow(query, id).Scan(&foo.ID, &foo.Field1, &foo.Field2, &foo.Field3)
	}
	return nil
}

var t time.Time

func startMeasure(v string) {
	text = v
	t = time.Now()
}

func endMeasure() {
	n := float64(time.Since(t).Seconds()) / float64(LoadCount)
	log.Printf("[%s] %.9f seconds per load", text, n)
}

func setup() error {
	dbPath := "./sqlite/test.db"

	if err := os.Remove(dbPath); err != nil {
		return fmt.Errorf("failed to remove db file: %v", err)
	}

	if v, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared&mode=memory", dbPath)); err != nil {
		return fmt.Errorf("failed to open: %w", err)
	} else {
		db = v
	}

	if _, err := db.Exec("CREATE TABLE `Foo` (`ID` String PRIMARY KEY, `Field1` String, `Field2` String, `Field3` String)"); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	values := []string{}
	for i := 1; i <= RecordCount; i++ {
		id := strconv.Itoa(i)
		foo := Foo{
			ID:     id,
			Field1: fmt.Sprintf("foo-%d", i),
			Field2: fmt.Sprintf("foo-%d", i),
			Field3: fmt.Sprintf("foo-%d", i),
		}
		FooMap[id] = foo
		values = append(values, fmt.Sprintf("('%s', '%s', '%s', '%s')", foo.ID, foo.Field1, foo.Field2, foo.Field3))
	}
	if _, err := db.Exec(fmt.Sprintf("INSERT INTO `Foo` (`ID`, `Field1`, `Field2`, `Field3`) VALUES %s", strings.Join(values, ","))); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	return nil
}
