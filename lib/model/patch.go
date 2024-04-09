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

func (item *Patch) AddDimmer(address int) (*Patch, *fixtures.Fixture) {
	fixture := fixtures.NewFixture(address, fixtures.Dimmer)
	patch.Fixtures = append(patch.Fixtures, fixture)
	return item, fixture
}
