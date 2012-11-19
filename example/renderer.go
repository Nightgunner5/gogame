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
</style>
<script>
// Array Remove - By John Resig (MIT Licensed)
Array.prototype.remove = function(from, to) {
	var rest = this.slice((to || from) + 1 || this.length);
	this.length = from < 0 ? this.length + from : from;
	return this.push.apply(this, rest);
};

setInterval(function() {
	var req = new XMLHttpRequest();
	req.open('GET', 'state', true);
	req.addEventListener('load', function() {
		var have = {};

		var response = JSON.parse(req.responseText);
		response.magicians.forEach(function(m) {
			update(have, true, 'ent' + m.id, m.x, m.y, m.z, m.health, m.mana, m.spell);
		});
		response.imps.forEach(function(i) {
			update(have, false, 'ent' + i.id, i.x, i.y, i.z, i.health, 0, i.spell);
		});

		var entities = [];
		Array.prototype.push.apply(entities, document.body.children);
		entities.forEach(function(ent) {
			if (!have[ent.id]) {
				ent.parentNode.removeChild(ent);
			}
		});
	});
	req.send();
}, 50);


function update(have, big, id, x, y, z, health, mana, spell) {
	var ent = document.getElementById(id);
	if (!ent) {
		ent = document.createElement(big ? 'div' : 'span');
		ent.id = id;
		ent.appendChild(document.createElement('span'));
		if (big) ent.appendChild(document.createElement('span'));
		document.body.appendChild(ent);
	}
	have[id] = true;
	ent.style.left = ((x + 10) * 5) + '%';
	ent.style.top = ((y + 10) * 5) + '%';
	ent.firstChild.style.width = (health * (big ? 1 : 10)) + '%';
	if (big)
		ent.lastChild.style.width = (mana / 10) + '%';

	if (spell) {
		ent.style.boxShadow = '0 0 ' + (100 - spell.progress * 100) + 'px #000';
	} else {
		ent.style.boxShadow = '';
	}
}
</script>
</head>
<body>
</body>
</html>`))
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"magicians":[`)
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
				fmt.Fprintf(w, `{"id":%d,"x":%v,"y":%v,"z":%v,"health":%v,"mana":%v,"spell":%v}`, m.ID(), x, y, z, m.Health(), m.Resource(), currentSpell)
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
