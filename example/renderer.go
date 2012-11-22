package main

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	"log"
	"net/http"
)

func spellProgress(s spell.Spell) float64 {
	return (s.TotalTime() - s.TimeLeft()) / s.TotalTime()
}

func renderer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>GoGame</title>
<style>
html {
	overflow: hidden;
	background: #333;
}
body>span {
	position: absolute;
	width: 5px;
	height: 5px;
	background: #c00;
	display: block;
	border: 1px solid #000;
}
body>div {
	position: absolute;
	width: 20px;
	height: 10px;
	background: #000;
	display: block;
	border: 1px solid #000;
}
div>span, span>span {
	background: #0f0;
	height: 5px;
	display: block;
}
header, footer {
	position: fixed;
	right: 0;
	z-index: 1000;
	background: rgba(255, 255, 255, 0.7);
}
header {
	top: 0;
	line-height: 0;
	width: 100px;
}
footer {
	bottom: 0;
}
#health, #mana {
	height: 5px;
	width: 50px;
	display: inline-block;
	background: #000;
}
#health div, #mana div {
	height: 5px;
	background: #0f0;
}
#mana div {
	background: #00f;
}
#spellprogress {
	height: 10px;
	width: 100px;
	background: #000;
}
#spellprogress div {
	height: 10px;
	background: #fff;
}
header p {
	font: 10px/1.1 sans-serif;
}
</style>
<script src="engine.js"></script>
<script>
gogame.client.start('ws://localhost:7031/socket');
</script>
</head>
<body><header><div id="health"><div></div></div><div id="mana"><div></div></div><div id="spellprogress"><div></div></div><p></p></header><footer><button onclick="spell('imp')">SUMMON IMP</button> <button onclick="spell('shield')">SHIELD SELF</button></footer></body>
</html>`))
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"self":%d,"magicians":[`, getMagician(r.Header.Get("X-Forwarded-For")).ID())
		first := true
		entity.ForEach(func(e entity.Entity) {
			if m, ok := e.(Magician); ok {
				if first {
					first = false
				} else {
					fmt.Fprintf(w, ",")
				}

				currentSpell := `null`
				if s := m.CurrentSpell(); s != nil {
					tar := s.(*spell.BasicSpell).Target_
					currentSpell = fmt.Sprintf(`{"id":%q,"progress":%v,"target":%d}`, s.(*spell.BasicSpell).Tag, spellProgress(s), tar)
				}

				x, y, z := m.Position()
				fmt.Fprintf(w, `{"id":%d,"x":%v,"y":%v,"z":%v,"health":%v,"mana":%v,"spell":%v,"effects":%q}`, m.ID(), x, y, z, m.Health(), m.Resource(), currentSpell, m.EffectDescription())
			}
		})

		fmt.Fprintf(w, `],"imps":[`)

		first = true
		entity.ForEach(func(e entity.Entity) {
			if i, ok := e.(Imp); ok {
				if first {
					first = false
				} else {
					fmt.Fprintf(w, ",")
				}

				currentSpell := `null`
				if s := i.CurrentSpell(); s != nil {
					tar := s.(*spell.BasicSpell).Target_
					currentSpell = fmt.Sprintf(`{"id":"impfire","progress":%v,"target":%d}`, spellProgress(s), tar)
				}

				x, y, z := i.Position()
				fmt.Fprintf(w, `{"id":%d,"x":%v,"y":%v,"z":%v,"health":%v,"spell":%v}`, i.ID(), x, y, z, i.Health(), currentSpell)
			}
		})

		fmt.Fprintf(w, `]}`)
	})

	log.Print("Open a browser to http://localhost:7031/")
	log.Fatal(http.ListenAndServe(":7031", nil))
}
