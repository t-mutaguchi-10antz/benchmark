package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/t-mutaguchi-10antz/benchmark/sqlite/boiler"
	"github.com/volatiletech/null/v8"
)

const (
	dbPath      string = "./sqlite/test.db"
	RecordCount int    = 100000
	LoadCount   int    = 100000
)

var db *sql.DB
var text string

var Heap = map[string]Sample{}

type Sample struct {
	ID     string
	Field1 string
	Field2 string
	Field3 string
}

func main() {
	if err := setup(); err != nil {
		panic(err)
	}

	startMeasure("heap memory")
	if err := load1(); err != nil {
		panic(err)
	}
	endMeasure()

	startMeasure("database/sql")
	if err := load2(); err != nil {
		panic(err)
	}
	endMeasure()

	startMeasure("volatiletech/sqlboiler ( ORM )")
	if err := load3(); err != nil {
		panic(err)
	}
	endMeasure()
}

func load1() error {
	for i := 1; i <= LoadCount; i++ {
		id := strconv.Itoa(i)
		_ = Heap[id]
	}
	return nil
}

func load2() error {
	query := "SELECT * FROM `Sample` WHERE `ID` = ?"
	for i := 1; i <= LoadCount; i++ {
		id := strconv.Itoa(i)
		var sample Sample
		_ = db.QueryRow(query, id).Scan(&sample.ID, &sample.Field1, &sample.Field2, &sample.Field3)
	}
	return nil
}

func load3() error {
	ctx := context.Background()

	for i := 1; i <= LoadCount; i++ {
		id := strconv.Itoa(i)
		if _, err := boiler.FindSample(ctx, db, null.NewString(id, true)); err != nil {
			panic(err)
		}
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
	log.Printf("%.9f sec/load [%s]", n, text)
}

func setup() error {
	if v, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared&mode=memory", dbPath)); err != nil {
		return fmt.Errorf("failed to open: %w", err)
	} else {
		db = v
	}

	gen := true
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		gen = false
	}

	if gen {
		if _, err := db.Exec("CREATE TABLE `Sample` (`ID` String PRIMARY KEY, `Field1` String, `Field2` String, `Field3` String)"); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	values := []string{}
	for i := 1; i <= RecordCount; i++ {
		id := strconv.Itoa(i)
		s := Sample{
			ID:     id,
			Field1: fmt.Sprintf("foo-%d", i),
			Field2: fmt.Sprintf("bar-%d", i),
			Field3: fmt.Sprintf("baz-%d", i),
		}
		Heap[id] = s
		values = append(values, fmt.Sprintf("('%s', '%s', '%s', '%s')", s.ID, s.Field1, s.Field2, s.Field3))
	}

	if gen {
		if _, err := db.Exec(fmt.Sprintf("INSERT INTO `Sample` (`ID`, `Field1`, `Field2`, `Field3`) VALUES %s", strings.Join(values, ","))); err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}
	}

	return nil
}
