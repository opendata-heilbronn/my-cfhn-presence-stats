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
	"net/http"

	"github.com/dghubble/sling"
	_ "github.com/go-sql-driver/mysql"
)

type presence struct {
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	LastSeen string `json:"lastSeen,omitempty"`
}

type tokenRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type tokenResponse struct {
	Token string `json:"token,omitempty"`
}

func fetchPresencesFromAPI() {

	token := fetchPresenceAuthToken()
	if len(token) == 0 {
		log.Fatal("[✘ ] Fatal error could not get auth token \n")
		return
	}

	sqlInsertPresence, err := db.Prepare("INSERT INTO `presences` (`username`, `datetime`) VALUES (?,?)")
	if err != nil {
		log.Fatalf("[✘ ] Fatal error database could not prepare insert: %s \n", err)
		return
	}

	// Get presences
	client := &http.Client{}
	presences := new([]presence)
	_, err = sling.New().
		Client(client).
		Get(config.GetString("presence_api.server")).
		Set("Authorization", token).
		ReceiveSuccess(presences)
	if err != nil {
		log.Fatalf("[✘ ] Fatal error when trying to get presences: %s \n", err)
		return
	}

	// Store the presences
	for _, p := range *presences {
		if _, err = sqlInsertPresence.Exec(p.Username, p.LastSeen); err != nil {
			log.Println("[✘ ] Failed to insert presence", err)
		}
	}

	// Done
	log.Println("[✓ ] Presences inserted")
}

func fetchPresenceAuthToken() string {
	request := tokenRequest{
		Username: config.GetString("presence_api.user"),
		Password: config.GetString("presence_api.password"),
	}
	response := new(tokenResponse)

	_, err := sling.New().
		Client(&http.Client{}).
		BodyJSON(request).
		Post(config.GetString("presence_api.login")).
		ReceiveSuccess(response)
	if err != nil {
		log.Fatalf("[✘ ] Fatal error when trying to get auth token: %s \n", err)
		return ""
	}

	return response.Token
}
