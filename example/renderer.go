package main

import (
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
html, body, canvas { height: 100%; margin: 0; padding: 0; overflow: hidden; }
header, footer { position: fixed; right: 0; z-index: 1000; background: rgba(255, 255, 255, 0.7); }
header { top: 0; line-height: 0; width: 100px; }
footer { bottom: 0; }
#health, #mana { height: 5px; width: 50px; display: inline-block; background: #000; }
#health div, #mana div { height: 5px; background: #0f0; }
#mana div { background: #00f; }
#spellprogress { height: 10px; width: 100px; background: #000; }
#spellprogress div { height: 10px; background: #fff; }
header p { font: 10px/1.1 sans-serif; padding: 3px; }
</style>
<script src="engine.js"></script>
<script>
var Handshake = (parseInt(gogame.net.FirstUnusedPacketID, 32) + 0).toString(32);
var CastSpell = (parseInt(gogame.net.FirstUnusedPacketID, 32) + 1).toString(32);
var KeepAlive = (parseInt(gogame.net.FirstUnusedPacketID, 32) + 2).toString(32);

gogame.client.start('ws://' + location.host + '/socket');

var myMagicianID;
gogame.client.listen(Handshake, function(packet) {
	myMagicianID = packet.get(gogame.net.EntityID);
});

requestAnimationFrame = window.requestAnimationFrame || window.mozRequestAnimationFrame || window.webkitRequestAnimationFrame || window.msRequestAnimationFrame;

setInterval(function() {
	gogame.client.send(new gogame.net.Packet(KeepAlive));
}, 30000);
gogame.client.listen(KeepAlive, function(packet) {
	// do nothing
});

requestAnimationFrame(function() {
	requestAnimationFrame(function render() {
		requestAnimationFrame(render);

		var canvas = document.querySelector('canvas'),
			ctx = canvas.getContext('2d');
		canvas.width = window.innerWidth;
		canvas.height = window.innerHeight;

		ctx.translate(canvas.width / 2, canvas.height / 2);

		if (gogame.client.disconnected) {
			ctx.font = '24px sans-serif';
			ctx.fillStyle = '#000';
			ctx.fillText('Disconnected', 0, 0);
			return;
		}

		if (myMagicianID && gogame.client.Entities[myMagicianID] && gogame.client.Entities[myMagicianID].position) {
			var pos = gogame.client.Entities[myMagicianID].position;
			ctx.translate(-pos[0]*10, -pos[1]*10);
		} else {
			ctx.font = '24px sans-serif';
			ctx.fillStyle = '#000';
			ctx.fillText('You are dead!', 0, 0);
			return;
		}

		for (var id in gogame.client.Entities) {
			var ent = gogame.client.Entities[id];
			if (!ent.tag || !ent.position || !ent.health)
				continue;
			var x = ent.position[0]*10, y = ent.position[1]*10;

			ctx.fillStyle = '#f00';
			ctx.fillRect(x, y, {
				imp: 10,
				magician: 100
			}[ent.tag], 1);

			ctx.fillStyle = '#0f0';
			ctx.fillRect(x, y, ent.health, 1);

			ctx.font = {
				imp: '12px sans-serif',
				magician: '18px sans-serif'
			}[ent.tag];
			ctx.fillStyle = '#000';
			ctx.fillText(ent.tag, x, y);
		}
	});
});

function spell(name) {
	gogame.client.send(new gogame.net.Packet(CastSpell).set(CastSpell, name));
}

gogame.client.send(new gogame.net.Packet(Handshake));
</script>
</head>
<body><header><div id="health"><div></div></div><div id="mana"><div></div></div><div id="spellprogress"><div></div></div><p></p></header><footer><button onclick="spell('imp')">SUMMON IMP</button> <button onclick="spell('shield')">SHIELD SELF</button></footer><canvas width="1000" height="1000"></canvas></body>
</html>`))
	})

	log.Print("Open a browser to http://localhost:7031/")
	log.Fatal(http.ListenAndServe(":7031", nil))
}
