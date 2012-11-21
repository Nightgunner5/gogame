/** @const */
gogame.net = gogame.net || {};

/** @constructor */
gogame.net.Socket = function(uri) {
	this.socket = new WebSocket(uri);
	var that = this;
	this.socket.addEventListener('message', function(message) {
		that.message_(JSON.parse(message['data']));
	}, false);
};

gogame.net.Socket.prototype.send = function(packet) {
	this.socket.send(JSON.stringify(packet));
};

/** @private */
gogame.net.Socket.prototype.message_ = function(packet) {
	console.log(packet);
};

/** @private */
gogame.net.base32 = function(n) {
	return n.toString(32);
};

/** @private */
gogame.net.iota = function() {
	return gogame.net.base32(gogame.net.iota_++);
};

/** @private */
gogame.net.iota_ = 0;

/** @const */ gogame.net.AttackerID   = gogame.net.iota();
/** @const */ gogame.net.VictimID     = gogame.net.iota();
/** @const */ gogame.net.Amount       = gogame.net.iota();
/** @const */ gogame.net.ChangeHealth = gogame.net.iota();
/** @const */ gogame.net.debugEcho    = gogame.net.iota();

/** @constructor */
gogame.net.Packet = function(id) {
	this['i'] = id;
	this['p'] = {};
};

gogame.net.Packet.prototype.set = function(key, value) {
	this['p'][key] = value;
	return this;
};

gogame.net.Packet.prototype.get = function(key) {
	return this['p'][key];
};
