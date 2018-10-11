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
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type apiResponse struct {
	Total            []userVisits `json:"user_total"`
	LastWeek         []userVisits `json:"user_lastweek"`
	LastWeekOverview []timeVists  `json:"overview_lastweek"`
}

type userVisits struct {
	Username string `json:"username"`
	Visits   int    `json:"visits"`
}

type timeVists struct {
	Day    string `json:"day"`
	Hour   int    `json:"hour"`
	Visits int    `json:"visits"`
}

func apiGetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	var response apiResponse
	response.Total = totalVisits()
	response.LastWeek = lastWeekVisits()
	response.LastWeekOverview = lastWeekUserCount()
	json.NewEncoder(w).Encode(response)
}

func totalVisits() []userVisits {

	sql, err := db.Prepare("SELECT `username`, COUNT(`datetime`) AS presences" +
		" FROM `presences`" +
		" GROUP BY `username`" +
		" ORDER BY `presences` DESC LIMIT 10")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare query: %s \n", err)
		return make([]userVisits, 0)
	}

	rows, err := sql.Query()
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not run query: %s \n", err)
		return make([]userVisits, 0)
	}
	defer rows.Close()

	visits := make([]userVisits, 0)
	for rows.Next() {
		var v userVisits
		if err := rows.Scan(&v.Username, &v.Visits); err != nil {
			log.Fatal(err)
		}

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return visits
}

func lastWeekVisits() []userVisits {

	sql, err := db.Prepare("SELECT `username`, COUNT(`datetime`) AS presences" +
		" FROM `presences`" +
		" WHERE `datetime` > DATE_ADD(CURRENT_TIMESTAMP(), INTERVAL -7 DAY)" +
		" GROUP BY `username`" +
		" ORDER BY `presences` DESC LIMIT 10")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare query: %s \n", err)
		return make([]userVisits, 0)
	}

	rows, err := sql.Query()
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not run query: %s \n", err)
		return make([]userVisits, 0)
	}
	defer rows.Close()

	visits := make([]userVisits, 0)
	for rows.Next() {
		var v userVisits
		if err := rows.Scan(&v.Username, &v.Visits); err != nil {
			log.Fatal(err)
		}

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return visits
}

func lastWeekUserCount() []timeVists {

	sql, err := db.Prepare("SELECT DATE_FORMAT(`datetime`, '%W %d.%m.%Y'), HOUR(`datetime`), COUNT(DISTINCT `username`) AS presences" +
		" FROM `presences`" +
		" WHERE `datetime` > DATE_ADD(CURRENT_DATE(), INTERVAL -7 DAY)" +
		" GROUP BY DAY(`datetime`), DATE_FORMAT(`datetime`, '%W %d.%m.%Y'), HOUR(`datetime`)" +
		" ORDER BY DAY(`datetime`) ASC, HOUR(`datetime`) ASC")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare query: %s \n", err)
		return make([]timeVists, 0)
	}

	rows, err := sql.Query()
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not run query: %s \n", err)
		return make([]timeVists, 0)
	}
	defer rows.Close()

	visits := make([]timeVists, 0)
	for rows.Next() {
		var v timeVists

		if err := rows.Scan(&v.Day, &v.Hour, &v.Visits); err != nil {
			log.Fatal(err)
		}

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return visits
}
