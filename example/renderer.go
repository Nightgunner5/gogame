package main

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"log"
	"net/http"
)

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
		var entities = [];
		Array.prototype.push.apply(entities, document.querySelectorAll('body>*'));
		req.responseText.split('\n').forEach(function(line) {
			line = line.split(' ');
			if (line.length != 7)
				return;
			var id = 'ent'+line[0], type = line[1], x = line[2]-0, y = line[3]-0, z = line[4]-0, health = line[5]-0, mana = line[6]-0;
			switch (type) {
			case 'magician':
				update(entities, true, id, x, y, z, health, mana);
				break;
			case 'imp':
				update(entities, false, id, x, y, z, health, mana);
				break;
			}
		});
		entities.forEach(function(ent) {
			ent.parentNode.removeChild(ent);
		});
	});
	req.send();
}, 50);


function update(entities, big, id, x, y, z, health, mana) {
	var ent = document.getElementById(id);
	if (ent) {
		entities.remove(entities.indexOf(ent), entities.indexOf(ent));
	} else {
		ent = document.createElement(big ? 'div' : 'span');
		ent.id = id;
		ent.appendChild(document.createElement('span'));
		if (big) ent.appendChild(document.createElement('span'));
		document.body.appendChild(ent);
	}
	ent.style.left = ((x + 10) * 5) + '%';
	ent.style.top = ((y + 10) * 5) + '%';
	ent.firstChild.style.width = (health * (big ? 1 : 10)) + '%';
	if (big)
		ent.lastChild.style.width = (mana / 10) + '%';
}
</script>
</head>
<body>
</body>
</html>`))
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		entity.ForEach(func(e entity.Entity) {
			if m, ok := e.(Magician); ok {
				x, y, z := m.Position()
				fmt.Fprintf(w, "%d magician %f %f %f %f %f\n", m.ID(), x, y, z, m.Health(), m.Resource())
				return
			}
			if i, ok := e.(Imp); ok {
				x, y, z := i.Position()
				fmt.Fprintf(w, "%d imp %f %f %f %f %f\n", i.ID(), x, y, z, i.Health(), 0)
				return
			}
		})
	})

	log.Print("Open a browser to http://localhost:7031/")
	log.Fatal(http.ListenAndServe(":7031", nil))
}
