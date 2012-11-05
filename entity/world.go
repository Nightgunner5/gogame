package entity

type world struct{}

func (*world) ID() EntityID {
	return 0
}

func (*world) Parent() Entity {
	return nil
}

var World Entity = new(world)

func init() {
	Spawn(World)
}
