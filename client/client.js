gogame.client = new gogame.net.Socket('ws://localhost:7031/socket');

gogame.client.send(new gogame.net.Packet(debugEcho)
		.set(gogame.net.Amount, 'abc'));
