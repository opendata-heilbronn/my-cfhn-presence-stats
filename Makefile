all:
	go build
	cd web && make build-js
	my-cfhn-presence-stats server
