package server

import (
	"flag"
	"log"
)

func StartServer(mapPath string, useRand bool, rows int, columns int, humans int, monsters int) {
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
	go s.run()
	startWebApp(s.Map)
}
