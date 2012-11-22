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
html, body, canvas { height: 100%; margin: 0; padding: 0; overflow: hidden; background: #333; }
header, footer { position: fixed; right: 0; z-index: 1000; background: rgba(255, 255, 255, 0.7); }
footer { bottom: 0; }
</style>
<script src="engine.js"></script>
<script>
var packetID = parseInt(gogame.net.FirstUnusedPacketID, 32);

var Handshake  = (packetID++).toString(32);
var EntityName = (packetID++).toString(32);
var CastSpell  = (packetID++).toString(32);
var KeepAlive  = (packetID++).toString(32);

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

var entityNames = {};
gogame.client.listen(EntityName, function(packet) {
	entityNames[packet.get(gogame.net.EntityID)] = packet.get(EntityName);
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
			ctx.fillStyle = '#fff';
			ctx.fillText('Disconnected', 0, 0);
			return;
		}

		if (myMagicianID) {
			if (gogame.client.Entities[myMagicianID]) {
				var pos = gogame.client.Entities[myMagicianID].position;
				if (pos) {
					ctx.translate(-pos[0]*20, -pos[1]*20);
				}
			} else {
				ctx.font = '24px sans-serif';
				ctx.fillStyle = '#fff';
				ctx.fillText('You are dead!', 0, 0);
				return;
			}
		}

		for (var id in gogame.client.Entities) {
			var ent = gogame.client.Entities[id];
			if (!ent.tag || !ent.position || !ent.health)
				continue;
			var x = ent.position[0]*20, y = ent.position[1]*20;

			ctx.fillStyle = '#f00';
			ctx.fillRect(x, y, {
				imp: 10,
				magician: 100
			}[ent.tag], 2);

			ctx.fillStyle = '#0f0';
			ctx.fillRect(x, y, ent.health, 2);

			ctx.font = {
				imp: '12px sans-serif',
				magician: '18px sans-serif'
			}[ent.tag];
			ctx.fillStyle = '#fff';
			ctx.fillText(entityNames[id] || 'Unknown', x, y);
		}
	});
});

var lastMoveX = 0, lastMoveY = 0;
function move(x, y) {
	if (lastMoveX == x && lastMoveY == y) {
		x = 0;
		y = 0;
	}
	gogame.client.send(new gogame.net.Packet(gogame.net.EntityPosition).set(gogame.net.EntityPosition, [x, y, 0]));
	lastMoveX = x;
	lastMoveY = y;
}

function spell(name) {
	gogame.client.send(new gogame.net.Packet(CastSpell).set(CastSpell, name));
}

gogame.client.send(new gogame.net.Packet(Handshake).set(EntityName, /*prompt("YO MAN WHAT'S YO NAME")*/'magician'));
</script>
</head>
<body><footer><button onclick="move(0,-1)">&uarr;</button><button onclick="move(-1,0)">&larr;</button><button onclick="move(1,0)">&rarr;</button><button onclick="move(0,1)">&darr;</button>
<br><button onclick="spell('imp')">SUMMON IMP</button> <button onclick="spell('shield')">SHIELD SELF</button></footer><canvas></canvas></body>
</html>`))
	})

	log.Print("Open a browser to http://localhost:7031/")
	log.Fatal(http.ListenAndServe(":7031", nil))
}
