package fixtures

type ChannelRole int

const (
	None ChannelRole = iota
	Dimmer
	Red
	Green
	Blue
	White
	Color
	Gobo
)

var channelRolesNames = [...]string{
	"None",
	"Dimmer",
	"Red",
	"Green",
	"Blue",
	"White",
	"Color",
	"Gobo",
}

var ColorRelatedChannelRoles = [...]ChannelRole{
	Red,
	Green,
	Blue,
	White,
}

func (item ChannelRole) String() string {
	return channelRolesNames[item]
}

func (item ChannelRole) EnumIndex() int {
	return int(item)
}

func IsChannelRole(value int) bool {
	return value <= len(channelRolesNames)
}
