package fixtures

import (
	"slices"
	"strconv"

	"fyne.io/fyne/v2/data/binding"
)

type Fixture struct {
	address       int
	channels      []FixtureChannel
	roleToAddress map[ChannelRole]int
	virtualDimmer binding.Int
	color         binding.Bytes
	colorChannels map[ChannelRole]FixtureChannel
	name          string
}

func NewFixture(address int, name string, channels ...ChannelRole) *Fixture {
	fixtureChannels := []FixtureChannel{}
	roleToAddress := map[ChannelRole]int{}
	colorChannels := map[ChannelRole]FixtureChannel{}
	virtualDimmer := binding.NewInt()
	virtualDimmer.Set(128)

	fixture := Fixture{
		address: address,
	}
	for i := 0; i < int(len(channels)); i++ {
		role := channels[i]
		channel := FixtureChannel{
			role:   role,
			source: binding.NewInt(),
			value:  binding.NewInt(),
		}
		fixtureChannels = append(fixtureChannels, channel)
		roleToAddress[role] = i + address
		if role == Dimmer {
			virtualDimmer = nil
		}
		if slices.Contains(ColorRelatedChannelRoles[:], role) {
			colorChannels[role] = channel
			channel.source.AddListener(binding.NewDataListener(func() {
				fixture.refreshColor()
			}))
		}
	}

	if virtualDimmer == nil {
		for _, channel := range fixtureChannels {
			channel.value = channel.source
		}
	}

	color := binding.NewBytes()
	color.Set([]byte{0, 0, 0})

	fixture.channels = fixtureChannels
	fixture.roleToAddress = roleToAddress
	fixture.virtualDimmer = virtualDimmer
	fixture.color = color
	fixture.colorChannels = colorChannels
	if len(name) == 0 {
		name = ""
		for _, channel := range fixtureChannels {
			name += channel.Role().String()[:1]
		}
		name += "_" + strconv.Itoa(address)
	}
	fixture.name = name

	return &fixture
}

func (item *Fixture) Bind(handler func()) {
	for _, channel := range item.channels {
		channel.value.AddListener(binding.NewDataListener(func() {
			handler()
		}))
	}
}

func (item *Fixture) BaseAddress() int {
	return item.address
}

func (item *Fixture) Name() string {
	return item.name
}

func (item *Fixture) Addresses(roles ...int) []int {
	addresses := []int{}
	if len(roles) == 0 {
		for address, _ := range item.channels {
			addresses = append(addresses, address+item.address)
		}
	} else {
		for i := 0; i < len(roles); i++ {
			addresses = append(addresses, item.roleToAddress[ChannelRole(roles[i])])
		}
	}
	return addresses
}

func (item *Fixture) Roles() []ChannelRole {
	roles := []ChannelRole{}
	for _, channel := range item.channels {
		roles = append(roles, channel.role)
	}
	return roles
}

func (item *Fixture) Dimmer() binding.Int {
	if item.virtualDimmer != nil {
		return item.virtualDimmer
	}
	return item.channels[item.roleToAddress[Dimmer]].source
}

func (item *Fixture) Color() binding.Bytes {
	return item.color
}

func (item *Fixture) GetChannelForAddress(address int) *FixtureChannel {
	if item.address <= address && item.address+len(item.channels) > address {
		return &item.channels[address-item.address]
	}
	return nil
}

func (item *Fixture) GetChannelForRole(role ChannelRole) *FixtureChannel {
	return item.getRoleChannel(role)
}

func (item *Fixture) refreshColor() {
	r := item.getRoleValue(Red)
	g := item.getRoleValue(Green)
	b := item.getRoleValue(Blue)
	// w := item.getColorValue(White)

	if item.virtualDimmer != nil {
		dimmer, _ := item.virtualDimmer.Get()
		ratio := (dimmer * 100) / 255
		item.getRoleChannel(Red).value.Set((int(r) * ratio) / 100)
		item.getRoleChannel(Green).value.Set((int(g) * ratio) / 100)
		item.getRoleChannel(Blue).value.Set((int(b) * ratio) / 100)
	}
	item.color.Set([]byte{r, g, b})
}

func (item *Fixture) getRoleValue(role ChannelRole) uint8 {
	channel := item.getRoleChannel(role)
	if channel != nil {
		value, _ := channel.source.Get()
		return uint8(value)
	}

	return 0
}

func (item *Fixture) getRoleChannel(role ChannelRole) *FixtureChannel {
	for _, channel := range item.channels {
		if channel.role == role {
			return &channel
		}
	}
	return nil
}

func FixtureChannelByAddress(fixtures []*Fixture, address int) *FixtureChannel {
	for _, fixture := range fixtures {
		channel := fixture.GetChannelForAddress(address)
		if channel != nil {
			return channel
		}
	}
	return nil
}
