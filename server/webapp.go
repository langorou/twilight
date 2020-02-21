package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/langorou/twilight/statik"
	"github.com/rakyll/statik/fs"
)

func startWebApp(m *Map) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(statikFS)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	index, err := template.New("index.html").Parse(Index)
	if err != nil {
		panic(err.Error())
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := index.Execute(w, packMap(m))
		if err != nil {
			log.Fatalf(err.Error())
		}
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mov := r.FormValue("mov")
		var offset int
		offset, _ = strconv.Atoi(mov)
		if offset > len(m.history) {
			// Just return empty
			json.NewEncoder(w).Encode([]int{})
		}
		json.NewEncoder(w).Encode(m.history[offset:])
	})

	log.Println("Web server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type packed struct {
	X, Y   int
	Humans []cell
	Vamps  []cell
	Wolfs  []cell
	State  string
	Mov    int
}

func packMap(m *Map) packed {
	p := packed{
		X:      m.Columns*80 + 1,
		Y:      m.Rows*80 + 1,
		Humans: make([]cell, 0),
		Vamps:  make([]cell, 0),
		Wolfs:  make([]cell, 0),
		Mov:    m.mov,
	}
	for _, i := range m.humans {
		p.Humans = append(p.Humans, scale(m.cells[i]))
	}
	for _, i := range m.monster[wolf-1] {
		p.Wolfs = append(p.Wolfs, scale(m.cells[i]))
	}
	for _, i := range m.monster[vamp-1] {
		p.Vamps = append(p.Vamps, scale(m.cells[i]))
	}
	switch m.state {
	case waiting:
		p.State = "Waiting"
	case ready:
		p.State = "Playing"
	case win0:
		p.State = "Player 0 won"
	case win1:
		p.State = "Player 1 won"
	}
	return p
}

func scale(c cell) cell {
	c.X = c.X * 80
	c.Y = c.Y * 80
	return c
}

const Index = `
<!DOCTYPE html>
<html lang="en">
	<head>
	<meta charset="utf-8">
	<script src="/static/js/vue.min.js"></script>
	<link rel="stylesheet" href="/static/css/bootstrap.min.css">
	</head>
	<body>
	<div class="row" id="game">
		<div id="info" class="col-md-3" style="height:80vh">
			<div style="background-color:rgb(42,79,110);">
				<img src="/static/img/human.svg">
				<h2>{{"{{totalH}}"}}</h2>
			</div>
			<div style="background-color:rgb(170,62,57);">
				<img src="/static/img/vamp.svg">
				<h2>{{"{{totalV}}"}}</h2>
			</div>
			<div style="background-color:rgb(141,163,54);">
				<img src="/static/img/wolf.svg">
				<h2>{{"{{totalW}}"}}</h2>
			</div>
		</div>

		<div id="grid" class="col-md-6">
			<div v-bind:style="dim">
				<svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
					<g v-for="h in history[idx].Humans">
						<rect v-bind:x="h.X" v-bind:y="h.Y" width="80" height="80"
									 stroke="black" fill="rgb(42,79,110)"/>
						<text v-bind:x="h.X+20" v-bind:y="h.Y+40">
							{{"{{h.c}}"}}
						</text>
					</g>
					<g v-for="v in history[idx].Vamps">
						<rect v-bind:x="v.X" v-bind:y="v.Y" width="80" height="80"
									 stroke="black" fill="rgb(170,62,57)"/>
						<text v-bind:x="v.X+20" v-bind:y="v.Y+40">
							{{"{{v.c}}"}}
						</text>
					</g>
					<g v-for="w in history[idx].Wolfs">
						<rect v-bind:x="w.X" v-bind:y="w.Y" width="80" height="80"
									 stroke="black" fill="rgb(141,163,54)"/>
						<text v-bind:x="w.X+20" v-bind:y="w.Y+40">
							{{"{{w.c}}"}}
						</text>
					</g>
					<defs>
						<pattern id="sgrid" width="80" height="80" patternUnits="userSpaceOnUse">
							<path d="M 80 0 L 0 0 0 80" fill="none" stroke="black" stroke-width="3"/>
						</pattern>
					</defs>
					<rect width="100%" height="100%" fill="url(#sgrid)" stroke-width="3"/>
					</svg>
			</div>
		</div>

		<div  class="col-md-3">
			<h3>Welcome to twilight</h3>
			<div style="display: flex;justify-content:space-around">
				<div >
					<h4>{{"{{history[idx].state}}"}}</h4>
					<h4>Mouvement {{"{{history[idx].mov}}"}}</h4>
				</div>
				<button class="btn btn-default" v-on:click="playpause">
					⏯️
				</button>
			</div>
			<div style="width:90%;height:90vh;overflow:auto;margin-top:3px">
				<table class="table table-striped">
					<tbody>
						<tr v-for="i in history.length" style="height:1.2em"> 
							<td v-on:click="setState(i-1)" role="button">
								Movement {{"{{i-1}}"}}
							</td>
					</tbody>
				</table>
			</div>
		</div>
	</div>
	</body>

	<script>
		var game = new Vue({
		  el: '#game',
		  data: {
			dim: {
				width: {{.X}} + 'px',
				height: {{.Y}} + 'px',
			},
			history: [{
				Humans: {{.Humans}},
				Vamps: {{.Vamps}},
				Wolfs: {{.Wolfs}},
				State: {{.State}},
				Mov: {{.Mov}},
			}],
			idx: 0,
			playing: true,
		  },
		  computed: {
			  totalH: function() {
				  return this.history[this.idx].Humans.reduce(function(acc, val) {
					  return acc + val.c
				  }, 0);
			  },
			  totalV: function() {
				  return this.history[this.idx].Vamps.reduce(function(acc, val) {
					  return acc + val.c
				  }, 0);
			  },
			  totalW: function() {
				  return this.history[this.idx].Wolfs.reduce(function(acc, val) {
					  return acc + val.c
				  }, 0);
			  },
		  },
		  methods: {
			playpause: function() {
				this.playing =  !this.playing;
			},
			setState: function(i) {
				this.playing = false;
				this.idx = i
			}
		  }
		});
		(function poll(){
			setTimeout(function(){
				var xhttp = new XMLHttpRequest();
				xhttp.onreadystatechange = function() {
					if (this.readyState == 4 && this.status == 200) {
						data = JSON.parse(this.responseText);
						if (data.length > 0) {
							i = data.length-1
							game.history = game.history.concat(data)
							if (game.playing) {
								game.idx = game.history.length-1;
							}
						}
						poll()
					}
				};
				xhttp.open("GET", "/data?mov="+game.history.length, true);
				xhttp.send();
			}, 300);
		})();
	</script>
	<style type="text/css">
	img{
		height:60%;
	}
	#info div {
		height:33%;
		display: flex;
		flex-direction: row;
		justify-content: space-around;
		align-items: center;
	}
	</style>
</html>
`
