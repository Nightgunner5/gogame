package entity

import ("runtime";"sync")

type EventData struct {
	Type string

	// The Entity that caused the event (attacker)
	Activator Entity

	// The Entity that sent this input (victim)
	Caller Entity

	// The (optional) data parameter
	Value interface{}
}

type queuedEvent struct {
	target Entity
	data   EventData
}

var eventQueue []queuedEvent
var queueCond = sync.NewCond(new(sync.Mutex))

func queueHandler() {
	for {
		queueCond.L.Lock()
		saveLock.RLock()
		for len(eventQueue) == 0 {
			saveLock.RUnlock()
			queueCond.Wait()
			saveLock.RLock()
		}
		var event queuedEvent
		event, eventQueue = eventQueue[0], eventQueue[1:]
		queueCond.L.Unlock()

		event.target.AcceptEvent(event.data)

		saveLock.RUnlock()
	}
}

func init() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go queueHandler()
	}
}

func QueueEvent(target Entity, data EventData) {
	queueCond.L.Lock()
	defer queueCond.L.Unlock()

	eventQueue = append(eventQueue, queuedEvent{
		target: target,
		data: data,
	})

	queueCond.Signal()
}
