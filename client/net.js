var net = gogame['net'] = gogame['net'] || {};

var Socket = net['Socket'] = /** @constructor */ function(uri) {
	var that = this;

	/** @private */ that.socket = new WebSocket(uri);
	/** @private */ that.listeners = {};

	that.socket['onmessage'] = function(message) {
		var packet = JSON.parse(message['data']);
		if (packet['i'] in that.listeners) {
			that.listeners[packet['i']](new Packet(packet));
		} else {
			console.log(packet);
		}
	};
	var toSend = "";
	that.socket['onopen'] = function(event) {
		that.socket.send(toSend);
		that['send'] = function(packet) {
			that.socket.send(JSON.stringify(packet));
		};
	};
	that.send = that['send'] = function(packet) {
		toSend += JSON.stringify(packet) + '\n';
	};
	that.listen = that['listen'] = function(packetID, listener) {
		that.listeners[packetID] = listener;
	};
};

/** @private */

/** @private */
net.base32 = function(n) {
	return n.toString(32);
};

/** @private */
net.iota = function() {
	net.iota_++;
	return net.base32(net.iota_ - 1);
};

/** @private */
net.iota_ = 0;

/** @const */ net.AttackerID   = net['AttackerID']   = net.iota();
/** @const */ net.VictimID     = net['VictimID']     = net.iota();
/** @const */ net.Amount       = net['Amount']       = net.iota();
/** @const */ net.ChangeHealth = net['ChangeHealth'] = net.iota();

/** @const */ net.debugEcho = net.iota();

/** @const */ net.EntityID        = net['EntityID']        = net.iota();
/** @const */ net.ParentID        = net['ParentID']        = net.iota();
/** @const */ net.Tag             = net['Tag']             = net.iota();

/** @const */ net.EntitySpawned   = net['EntitySpawned']   = net.iota();
/** @const */ net.EntityDespawned = net['EntityDespawned'] = net.iota();
/** @const */ net.EntityPosition  = net['EntityPosition']  = net.iota();

/** @const */ net['FirstUnusedPacketID'] = net.iota();

var Packet = net['Packet'] = /** @constructor */ function(id) {
	if (typeof id == 'object') {
		this['i'] = id['i'];
		this['p'] = id['p'];
	} else {
		this['i'] = id;
		this['p'] = {};
	}

	this.set = this['set'] = function(key, value) {
		this['p'][key] = value;
		return this;
	};

	this.get = this['get'] = function(key) {
		return this['p'][key];
	};
};
