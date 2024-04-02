package fixtures

import (
	"fyne.io/fyne/v2/data/binding"
)

type FixtureChannel struct {
	role   ChannelRole
	source binding.Int
	value  binding.Int
}

func (item *FixtureChannel) Role() ChannelRole {
	return item.role
}

func (item *FixtureChannel) Value() binding.Int {
	return item.value
}

func (item *FixtureChannel) Set(value byte) {
	item.source.Set(int(value))
}
