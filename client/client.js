var client = gogame['client'] = {};

var entities = client['Entities'] = {};

client.start = function(url) {
	client = gogame['client'] = new net['Socket'](url);
	client['Entities'] = entities;

	client.listen(net.EntitySpawned, function(packet) {
		entities[packet.get(net.EntityID)] = {
			'parent': packet.get(net.ParentID),
			'tag': packet.get(net.Tag)
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

	client.listen(net.ChangeHealth, function(packet) {
		entities[packet.get(net['VictimID'])]['health'] = packet.get(net.Amount);
	});

	client.listen(net.EntityPosition, function(packet) {
		entities[packet.get(net.EntityID)]['position'] = packet.get(net.EntityPosition);
	});
};
