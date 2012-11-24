var client = gogame['client'] = {};

var entities = client['Entities'] = {};

client['disconnected'] = false;

client['start'] = function(url) {
	client = gogame['client'] = new net['Socket'](url);
	client['Entities'] = entities;
	client['disconnected'] = false;

	client.socket['onerror'] = client.socket['onclose'] = function(event) {
		client['disconnected'] = true;
	};

	client.listen(net.EntitySpawned, function(packet) {
		entities[packet.get(net.EntityID)] = {
			'parent':   packet.get(net.OtherEntID),
			'tag':      packet.get(net.Tag),
			'position': [0, 0, 0],
			'effects':  []
		};
	});

	client.listen(net.EntityDespawned, function(packet) {
		(function despawnRecursive(parentID) {
			delete entities[parentID];
			for (var id in entities) {
				if (entities[id]['parent'] == parentID) {
					despawnRecursive(id);
				}
			}
		})(packet.get(net.EntityID));
	});

	client.listen(net.ChangeResource, function(packet) {
		entities[packet.get(net.EntityID)]['resource'] = packet.get(net.Amount);
	});

	client.listen(net.ChangeHealth, function(packet) {
		entities[packet.get(net.EntityID)]['health'] = packet.get(net.Amount);
	});

	client.listen(net.CastSpell, function(packet) {
		entities[packet.get(net.EntityID)]['spell'] = {
			'target':    packet.get(net.OtherEntID),
			'timeLeft':  packet.get(net.TimeLeft),
			'totalTime': packet.get(net.TotalTime),
			'tag':       packet.get(net.Tag)
		};
	});

	client.listen(net.EntityPosition, function(packet) {
		entities[packet.get(net.EntityID)]['position'] = packet.get(net.EntityPosition);
	});

	client.listen(net.EntityEffects, function(packet) {
		entities[packet.get(net.EntityID)]['effects'] = packet.get(net.EntityEffects);
	});
};

setInterval(function() {
	for (var id in entities) {
		if ('spell' in entities[id]) {
			if ((entities[id]['spell']['timeLeft'] -= 0.1) <= 0) {
				delete entities[id]['spell'];
			}
		}
	}
}, 100);
