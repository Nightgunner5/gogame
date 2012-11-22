var net = gogame['net'] = gogame['net'] || {};

var Socket = net['Socket'] = /** @constructor */ function(uri) {
	/** @private */ this.socket = new WebSocket(uri);
	/** @private */ this.listeners = {};
	var that = this;


	this.socket['onmessage'] = function(message) {
		if (packet['i'] in this.listeners) {
			this.listeners[packet['i']](new gogame.net.Packet(JSON.parse(message['data'])));
		}
	};
	this['send'] = function(packet) {
		this.socket.send(JSON.stringify(packet));
	};
	this['listen'] = function(packetID, listener) {
		this.listeners[packetID] = listener;
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

/** @const */ net['AttackerID']   = net.iota();
/** @const */ net['VictimID']     = net.iota();
/** @const */ net['Amount']       = net.iota();
/** @const */ net['ChangeHealth'] = net.iota();

/** @const */ net.debugEcho = net.iota();

/** @const */ net['EntityID']        = net.iota();
/** @const */ net['ParentID']        = net.iota();
/** @const */ net['EntitySpawned']   = net.iota();
/** @const */ net['EntityDespawned'] = net.iota();

var Packet = net['Packet'] = /** @constructor */ function(id) {
	if (typeof id == 'object') {
		this['i'] = id['i'];
		this['p'] = id['p'];
	} else {
		this['i'] = id;
		this['p'] = {};
	}

	this['set'] = function(key, value) {
		this['p'][key] = value;
		return this;
	};

	this['get'] = function(key) {
		return this['p'][key];
	};
};
