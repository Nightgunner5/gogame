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
var Handshake = (parseInt(gogame.net.FirstUnusedPacketID, 32) + 0).toString(32);
var CastSpell = (parseInt(gogame.net.FirstUnusedPacketID, 32) + 1).toString(32);

gogame.client.start('ws://localhost:7031/socket');

var myMagicianID;
gogame.client.listen(Handshake, function(packet) {
	myMagicianID = packet.get(gogame.net.EntityID);
});

requestAnimationFrame = window.requestAnimationFrame || window.mozRequestAnimationFrame || window.webkitRequestAnimationFrame || window.msRequestAnimationFrame;

requestAnimationFrame(function render() {
	console.log(gogame.client.Entities);

	requestAnimationFrame(render);
});

function spell(name) {
	gogame.client.send(new gogame.net.Packet(CastSpell).set(CastSpell, name));
}

gogame.client.send(new gogame.net.Packet(Handshake));
</script>
</head>
<body><header><div id="health"><div></div></div><div id="mana"><div></div></div><div id="spellprogress"><div></div></div><p></p></header><footer><button onclick="spell('imp')">SUMMON IMP</button> <button onclick="spell('shield')">SHIELD SELF</button></footer></body>
</html>`))
	})

	log.Print("Open a browser to http://localhost:7031/")
	log.Fatal(http.ListenAndServe(":7031", nil))
}
