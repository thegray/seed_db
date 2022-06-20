package command

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Hit struct {
	Id      int64
	Name    string
	Org     string
	IsGood  string
	IsOkay  string
	Freq    int
	Counter int
}

// var show_result = false

func QueryCommand(db *sql.DB, tblName string, filterMode string, limit int) ([]Hit, error) {
	var filter bool
	switch filterMode {
	case "code":
		filter = false
	case "mysql":
		filter = true
	default:
		err := errors.New("unknown filter value: " + filterMode)
		log.Println(err)
		return nil, err
	}

	start := time.Now()
	res, err := query(db, tblName, filter, limit)
	if err != nil {
		return nil, err
	}
	checkPoint1 := time.Since(start)
	log.Println("query time: ", checkPoint1)

	arrResult := make([]Hit, 0)
	if res != nil {
		hit := Hit{}
		for res.Next() {
			var id int64
			var name, org, isGood, isOkay string
			var freq, counter int
			err = res.Scan(&id, &name, &org, &isGood, &isOkay, &freq, &counter)
			if err != nil {
				panic(err.Error())
			}
			if !filter {
				if freq <= counter {
					continue
				}
			}
			hit.Id = id
			hit.Name = name
			hit.Org = org
			hit.IsGood = isGood
			hit.IsOkay = isOkay
			hit.Freq = freq
			hit.Counter = counter

			arrResult = append(arrResult, hit)
			// if show_result {
			// 	log.Println(hit)
			// }
		}
		checkPoint2 := time.Since(start)
		log.Println("query+processing time: ", checkPoint2)
		log.Println("processing time: ", checkPoint2-checkPoint1)
	} else {
		log.Println("empty result")
	}

	return arrResult, nil
}

func query(db *sql.DB, tblName string, filter bool, limit int) (*sql.Rows, error) {
	defer db.Close()

	sql := "SELECT * FROM " + tblName + " WHERE is_good = 'Y' AND is_okay = 'Y'"
	if filter {
		sql += " AND freq > counter"
	}
	if limit > 0 {
		l := fmt.Sprintf(" LIMIT %d", limit)
		sql += l
	}
	log.Println(sql)

	stmt, _ := db.Prepare(sql)
	res, err := stmt.Query()

	if err != nil {
		return nil, err
	}
	return res, nil
}
