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
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type apiResponse struct {
	GeneratedAt      string       `json:"generated_at"`
	Total            []userVisits `json:"user_total"`
	LastWeek         []userVisits `json:"user_lastweek"`
	TotalAlone       []userVisits `json:"user_alone"`
	LongestStreaks   []userVisits `json:"user_streaks"`
	UserVisits       []userVisits `json:"user_visits"`
	DaysUsers        []userVisits `json:"day_users"`
	DaysVisits       []userVisits `json:"day_visits"`
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
	response.GeneratedAt = time.Now().Format(time.UnixDate)
	response.Total = totalVisits()
	response.LastWeek = lastWeekVisits()
	response.TotalAlone = totalAloneVisits()
	response.UserVisits = userVisitCount()
	response.LongestStreaks = longestStreaks()
	response.DaysUsers = daysUsers()
	response.DaysVisits = daysVisits()
	response.LastWeekOverview = lastWeekUserCount()

	json.NewEncoder(w).Encode(response)
}

func totalVisits() []userVisits {

	sql, err := db.Prepare("SELECT `username`, COUNT(`datetime`) AS `presences`" +
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

	sql, err := db.Prepare("SELECT `username`, COUNT(`datetime`) AS `presences`" +
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

func longestStreaks() []userVisits {

	sql, err := db.Prepare("SELECT `username`, MAX(`ticks`) AS `presences`" +
		" FROM `streaks`" +
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

func userVisitCount() []userVisits {

	sql, err := db.Prepare("SELECT `username`, COUNT(`ticks`) AS `presences`" +
		" FROM `streaks`" +
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

func totalAloneVisits() []userVisits {

	sql, err := db.Prepare("SELECT `username`, COUNT(`datetime`) AS `presences`" +
		" FROM `presences`" +
		" WHERE `datetime` IN (SELECT `datetime`" +
		"  FROM `presences`" +
		"  GROUP BY `datetime`" +
		"  HAVING COUNT(*) = 1)" +
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

	sql, err := db.Prepare("SELECT DATE_FORMAT(`datetime`, '%Y.%m.%d %H'), COUNT(DISTINCT `username`) AS `presences`" +
		" FROM `presences`" +
		" WHERE `datetime` > DATE_ADD(CURRENT_DATE(), INTERVAL -7 DAY)" +
		" GROUP BY DATE_FORMAT(`datetime`, '%Y.%m.%d %H')" +
		" ORDER BY DATE_FORMAT(`datetime`, '%Y.%m.%d %H') ASC")
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

	loc, _ := time.LoadLocation("Europe/Berlin")

	visits := make([]timeVists, 0)
	for rows.Next() {
		var v timeVists
		var date string

		if err := rows.Scan(&date, &v.Visits); err != nil {
			log.Fatal(err)
		}

		datetime, _ := time.Parse("2006.01.02 15", date)
		datetime = datetime.In(loc)
		v.Day = datetime.Format("Mon 02.Jan.2006")
		v.Hour, _ = strconv.Atoi(datetime.Format("15"))

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return visits
}

func daysUsers() []userVisits {

	sql, err := db.Prepare("SELECT DATE_FORMAT(`datetime`, '%d.%m.%Y'), COUNT(DISTINCT `username`) AS `presences`" +
		" FROM `presences`" +
		" GROUP BY DATE_FORMAT(`datetime`, '%d.%m.%Y')" +
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

	loc, _ := time.LoadLocation("Europe/Berlin")

	visits := make([]userVisits, 0)
	for rows.Next() {
		var v userVisits
		var date string

		if err := rows.Scan(&date, &v.Visits); err != nil {
			log.Fatal(err)
		}

		datetime, _ := time.Parse("02.01.2006", date)
		datetime = datetime.In(loc)
		v.Username = datetime.Format("Mon 02.Jan.2006")

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return visits
}

func daysVisits() []userVisits {

	sql, err := db.Prepare("SELECT DATE_FORMAT(`datetime`, '%d.%m.%Y'), COUNT(*) AS `presences`" +
		" FROM `presences`" +
		" GROUP BY DATE_FORMAT(`datetime`, '%d.%m.%Y')" +
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

	loc, _ := time.LoadLocation("Europe/Berlin")

	visits := make([]userVisits, 0)
	for rows.Next() {
		var v userVisits
		var date string

		if err := rows.Scan(&date, &v.Visits); err != nil {
			log.Fatal(err)
		}

		datetime, _ := time.Parse("02.01.2006", date)
		datetime = datetime.In(loc)
		v.Username = datetime.Format("Mon 02.Jan.2006")

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return visits
}
