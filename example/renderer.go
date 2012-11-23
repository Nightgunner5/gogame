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
footer { position: fixed; right: 0; z-index: 1000; bottom: 0; text-align: center; }
button { width: 150px; }
#movement, button { line-height: 1; padding: 2px; border-radius: 3px; border: 1px solid #aaa; display: inline-block; background: -webkit-linear-gradient(top, #eee, #ddd, #aaa); }
#movement button { width: 20px; border: 0; background: transparent; padding: 0; margin: 0; }
button { color: #333; }
button:hover { color: #000; }
</style>
<script src="engine.js"></script>
<script>
var packetID = parseInt(gogame.net.FirstUnusedPacketID, 32);

var Handshake  = (packetID++).toString(32);
var EntityName = (packetID++).toString(32);
var CastSpell  = (packetID++).toString(32);
var KeepAlive  = (packetID++).toString(32);

var VIEW_SCALE = 20;

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

var viewPos = {x: 0, y: 0};

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
					viewPos.x = pos[0], viewPos.y = pos[1];
				}
			} else {
				ctx.font = '24px sans-serif';
				ctx.fillStyle = '#fff';
				ctx.fillText('You are dead!', 0, 0);
			}
		}
		ctx.translate(-viewPos.x * VIEW_SCALE, -viewPos.y * VIEW_SCALE);

		for (var id in gogame.client.Entities) {
			var ent = gogame.client.Entities[id];
			if (!ent.tag || !ent.position || !ent.health)
				continue;
			var x = ent.position[0]*VIEW_SCALE, y = ent.position[1]*VIEW_SCALE;

			if (ent.tag == 'imp') {
				ctx.beginPath();
				ctx.fillStyle = ent.parent == myMagicianID ? 'rgba(0, 255, 0, 0.05)' : 'rgba(255, 0, 0, 0.05)';
				ctx.arc(x, y, VIEW_SCALE * 10, 0, Math.PI * 2);
				ctx.fill();
			}

			ctx.fillStyle = '#f00';
			ctx.fillRect(x, y, {
				imp: VIEW_SCALE / 2,
				magician: VIEW_SCALE * 5
			}[ent.tag], VIEW_SCALE / 10);

			ctx.fillStyle = '#0f0';
			ctx.fillRect(x, y, ent.health * VIEW_SCALE / 20, VIEW_SCALE / 10);

			ctx.font = {
				imp: '12px sans-serif',
				magician: '18px sans-serif'
			}[ent.tag];
			ctx.fillStyle = '#fff';
			ctx.fillText(entityNames[id] || 'Unknown', x, y);
		}
	});
});

function move(x, y) {
	gogame.client.send(new gogame.net.Packet(gogame.net.EntityPosition).set(gogame.net.EntityPosition, [x, y, 0]));
}

function spell(name) {
	gogame.client.send(new gogame.net.Packet(CastSpell).set(CastSpell, name));
}

gogame.client.send(new gogame.net.Packet(Handshake).set(EntityName, /*prompt("YO MAN WHAT'S YO NAME")*/'magician'));
</script>
</head>
<body><footer>
<div id="movement"><button onclick="move(-1,-1)">&nbsp;</button><button onclick="move(0,-1)">&uarr;</button><button onclick="move(1,-1)">&nbsp;</button><br/>
<button onclick="move(-1,0)">&larr;</button><button onclick="move(0, 0)">&nbsp;</button><button onclick="move(1,0)">&rarr;</button><br/>
<button onclick="move(-1,1)">&nbsp;</button><button onclick="move(0,1)">&darr;</button><button onclick="move(1,1)">&nbsp;</button></div><br/>
<button onclick="spell('imp')">SUMMON IMP</button><br/>
<button onclick="spell('shield')">SHIELD SELF</button>
</footer><canvas></canvas></body>
</html>`))
	})

	log.Print("Open a browser to http://localhost:7031/")
	log.Fatal(http.ListenAndServe(":7031", nil))
}
