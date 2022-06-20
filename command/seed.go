package command

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/rs/xid"
)

var logdetail bool

func SeedCommand(db *sql.DB, tblName string) error {
	var err error
	defer db.Close()
	start := time.Now()
	for x := 0; x < 994; x++ {
		err = seed(db, tblName, x)
		if err != nil {
			break
		}
	}
	duration := time.Since(start)
	log.Println("done seeding, elapsed: ", duration)

	return err
}

func seed(db *sql.DB, tblName string, x int) error {
	sql := "INSERT INTO " + tblName + "(name, org, is_good, is_okay, freq, counter) VALUES "
	vals := []interface{}{}

	start := time.Now()

	for i := 0; i < 1000; i++ {
		var org string
		var is_good string
		var is_okay string
		var freq int
		var counter int

		guid := xid.New()
		sId := guid.String()

		s1 := rand.NewSource(time.Now().UnixNano() + int64(guid.Counter()))
		r1 := rand.New(s1)
		i1 := r1.Intn(100)

		s2 := rand.NewSource(int64(guid.Counter()))
		r2 := rand.New(s2)
		i2 := r2.Intn(100)

		if i1 <= 40 {
			org = "TAMPAN"
		} else {
			org = "CANTIK"
		}
		if i2 <= 25 {
			is_good = "N"
			freq = i2
			counter = i2
		} else {
			is_good = "Y"
			freq = i2 + 2
			counter = i2
		}
		if i1 <= 30 {
			is_okay = "N"
		} else {
			is_okay = "Y"
		}

		sql += "(?, ?, ?, ?, ?, ?),"
		vals = append(vals, sId, org, is_good, is_okay, freq, counter)
	}

	sql = sql[0 : len(sql)-1]
	stmt, _ := db.Prepare(sql)
	res, err := stmt.Exec(vals...)

	if err != nil {
		panic(err.Error())
	}

	_, err = res.LastInsertId()

	if err != nil {
		log.Fatal(err)
		return err
	}

	// fmt.Printf("The last inserted row id: %d\n", lastId)
	if logdetail {
		duration := time.Since(start)
		log.Printf("%d. seeded for 1000, elapsed: %s\n", x, duration)
	}

	return nil
}
