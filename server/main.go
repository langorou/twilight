package server

import (
	"flag"
	"log"
	"time"
)

func StartServer(
	mapPath string,
	useRand bool,
	rows int,
	columns int,
	humans int,
	monster int,
	timeout time.Duration,
	useRandomPort bool,
	portUsed chan int,
	noWebApp bool,
	gameOutcome chan state,
) {
	flag.Parse()
	var names [2]string
	var m *Map
	if !useRand {
		if mapPath != "" {
			m = newMap(mapPath)
		} else {
			log.Println("Please either use -map with a valid filename or -rand for a random map")
			return
		}
	} else {
		m = generate(mapPath, rows, columns, humans, monster)
	}
	m.updateHistory()
	s := server{m, names}
	go s.run(timeout, useRandomPort, portUsed, gameOutcome)
	if !noWebApp {
		startWebApp(s.Map)
	}
}
