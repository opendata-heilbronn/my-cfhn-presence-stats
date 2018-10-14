all:
	go build
	cd web && make build-js
	my-cfhn-presence-stats server

fetch:
	go build
	my-cfhn-presence-stats fetch
