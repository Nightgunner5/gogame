package entity

type world struct{}

func (*world) ID() EntityID {
	return 0
}

func (*world) Parent() Entity {
	return nil
}

func (*world) Tag() string {
	return "world"
}

var World Entity = new(world)

func init() {
	globalEntityList.Add(World)
}
