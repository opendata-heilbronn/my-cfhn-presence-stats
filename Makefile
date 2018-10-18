all:
	go build
	cd web && make build-js
	my-cfhn-presence-stats server

fetch:
	go build
	my-cfhn-presence-stats fetch

streaks:
	go build
	my-cfhn-presence-stats streaks

test:
	go build
	my-cfhn-presence-stats test
