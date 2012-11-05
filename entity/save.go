package entity

import (
	"encoding/gob"
	"io"
)

func Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)

	idLock.Lock()
	defer idLock.Unlock()
	entities := globalEntityList.(*concurrentEntityList)
	entities.m.Lock()
	defer entities.m.Unlock()

	err = enc.Encode(uint64(nextID))
	if err != nil {
		return
	}

	err = enc.Encode(uint64(len(entities.l)))
	if err != nil {
		return
	}

	for _, e := range entities.l {
		err = enc.Encode(e)
		if err != nil {
			return
		}
	}
	return
}

func Load(r io.Reader) (err error) {
	dec := gob.NewDecoder(r)

	idLock.Lock()
	defer idLock.Unlock()
	entities := globalEntityList.(*concurrentEntityList)
	entities.m.Lock()
	defer entities.m.Unlock()

	err = dec.Decode(&nextID)
	if err != nil {
		return
	}

	var entityListSize uint64
	err = dec.Decode(&entityListSize)
	if err != nil {
		return
	}

	entities.l = make(entityList, 0, entityListSize)
	for i := uint64(0); i < entityListSize; i++ {
		var ent Entity
		err = dec.Decode(&ent)
		if err != nil {
			return
		}
		entities.l.Add(ent)
	}
	return
}
