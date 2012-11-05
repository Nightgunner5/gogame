package entity

import "encoding/gob"

type world struct{}

func (*world) ID() EntityID {
	return 0
}

func (*world) Parent() Entity {
	return nil
}

var World Entity = new(world)

func init() {
	gob.Register(World)
	Spawn(World)
}
