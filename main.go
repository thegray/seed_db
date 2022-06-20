package main

import (
	"database/sql"
	"flag"
	"log"
	"seed_db/command"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	queryPtr := flag.Bool("query", true, "a db query mode")
	seedPtr := flag.Bool("seed", false, "a db seeding mode")
	tblPtr := flag.String("table", "hit8", "a db name to be used")
	filterPtr := flag.String("filter", "code", "filter in code or mysql")
	limitPtr := flag.Int("limit", 0, "query limit, 0 is all")
	flag.Parse()

	db, err := sql.Open("mysql", "devroot:devroot@tcp(127.0.0.1:3306)/stress_test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if *seedPtr {
		err = command.SeedCommand(db, *tblPtr)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *queryPtr {
		res, err := command.QueryCommand(db, *tblPtr, *filterPtr, *limitPtr)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(len(res))
		showRes := false
		if showRes {
			for i, hit := range res {
				log.Println(i, hit)
			}
		}

		return
	}

}
