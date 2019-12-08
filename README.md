
Simulates a fight between aliens written in golang.

This code serves no practical purpose, as true aliens would operate in a much less predictable manner.

Simple execution:

`go run cli/main.go --map maps/official.map run 4 10000`

Run tests:

`make test`

Execute with logging turned on:

`go run cli/main.go --map maps/official.map run 4 10 --info-logging`

Verify a 'map' file:

`go run cli/main.go --map maps/official.map verify`


Milestones:

* Load and parse the map file.

* Alien structure and initial random location.

* Alien moving and fight detection

* Optimization


