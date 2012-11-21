gogame.client = new gogame.net.Socket('ws://localhost:7031/socket');

gogame.client.send(new gogame.net.Packet('a')
		.set(gogame.net.AttackerID, 1)
		.set(gogame.net.VictimID, 2)
		.set(gogame.net.Amount, 10)
		);
