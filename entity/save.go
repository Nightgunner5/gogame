package entity

import (
	"encoding/gob"
	"io"
	"fmt"
	"sync"
)

var saveLock sync.RWMutex

func Save(w io.Writer) error {
	saveLock.Lock()
	defer saveLock.Unlock()

	enc := gob.NewEncoder(w)

	if err := enc.Encode(uint64(1)); err != nil {
		return err
	}
	if err := enc.Encode(lastID); err != nil {
		return err
	}
	if err := enc.Encode(entityList); err != nil {
		return err
	}
	if err := enc.Encode(eventQueue); err != nil {
		return err
	}
	return nil
}

func Load(r io.Reader) error {
	saveLock.Lock()
	dec := gob.NewDecoder(r)
	var version uint64
	if err := dec.Decode(&version); err != nil {
		return err
	}

	switch version {
	case 1:
		if err := dec.Decode(&lastID); err != nil {
			return err
		}
		if err := dec.Decode(&entityList); err != nil {
			return err
		}
		if err := dec.Decode(&eventQueue); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown version %d", version)
	}

	return nil
}
