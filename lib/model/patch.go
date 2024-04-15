package model

import "github.com/pchaussalet/lights-go/lib/model/fixtures"

type Patch struct {
	Fixtures []*fixtures.Fixture
}

var patch Patch = Patch{
	Fixtures: []*fixtures.Fixture{},
}

func GetPatch() *Patch {
	return &patch
}

func (patch *Patch) AddFixture(address int, name string, subs ...fixtures.ChannelRole) {
	patch.Fixtures = append(patch.Fixtures, fixtures.NewFixture(address, name, subs...))
}
