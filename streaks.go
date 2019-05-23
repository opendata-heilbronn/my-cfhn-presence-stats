package main

// Copyright (C) 2018 Joas Schilling <coding@schilljs.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type streak struct {
	Username  string `json:"username,omitempty"`
	Arrival   string `json:"arrival,omitempty"`
	Departure string `json:"departure,omitempty"`
}

func recalculateStreaksFromPresences() {

	sqlOffset, err := db.Prepare("SELECT `departure`" +
		" FROM `streaks`" +
		" ORDER BY `departure` DESC" +
		" LIMIT 1")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare query: %s \n", err)
		return
	}

	dateOffset := ""
	row := sqlOffset.QueryRow()
	row.Scan(&dateOffset)

	sqlCount, err := db.Prepare("SELECT COUNT(*) AS `num_presences`" +
		" FROM `presences`" +
		" WHERE `datetime` > ?")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare query: %s \n", err)
		return
	}

	numRows := 0
	row = sqlCount.QueryRow(dateOffset)
	row.Scan(&numRows)

	log.Printf("[✓ ] Number of presence: %d", numRows)

	// db.Query("TRUNCATE `streaks`")

	sqlInsertStreak, err := db.Prepare("INSERT INTO `streaks` (`username`, `arrival`, `departure`, `ticks`) VALUES (?,?,?,1)")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare insert: %s \n", err)
		return
	}

	sqlUpdateStreak, err := db.Prepare("UPDATE `streaks` SET `departure` = ?, `ticks` = `ticks` + 1 WHERE `username` = ? AND `departure` = ?")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare update: %s \n", err)
		return
	}

	sqlSelect, err := db.Prepare("SELECT `username`, `datetime`" +
		" FROM `presences`" +
		" WHERE `datetime` > ?" +
		" ORDER BY `datetime` ASC")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare query: %s \n", err)
		return
	}

	rows, err := sqlSelect.Query(dateOffset)
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not run query: %s \n", err)
		return
	}
	defer rows.Close()

	bar := pb.StartNew(numRows)
	for rows.Next() {
		bar.Increment()
		var username, datetime string
		if err := rows.Scan(&username, &datetime); err != nil {
			log.Fatal(err)
		}

		// datetime
		lastTick, _ := time.Parse("2006-01-02 15:04:05", datetime)
		lastTick = lastTick.Add(time.Minute * -5)

		result, err := sqlUpdateStreak.Exec(datetime, username, lastTick)
		if err != nil {
			log.Println("[✘ ] Failed to insert presence", err)
			continue
		}

		if ok, _ := result.RowsAffected(); ok == 0 {
			sqlInsertStreak.Exec(username, datetime, datetime)
		}
		// log.Printf("[✓ ] Added %s", p.Username)
	}

	bar.FinishPrint("Done")
	log.Printf("[✓ ] Finished")

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
