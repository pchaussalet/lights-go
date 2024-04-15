package model

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/pchaussalet/lights-go/lib/model/fixtures"
)

type Closeable interface {
	Close()
}

type AppState struct {
	pending  []Closeable
	Fixtures []*fixtures.Fixture
	Mappings map[uint8]map[uint8]*fixtures.FixtureChannel
}

func (item *AppState) Add(p Closeable) {
	if item.pending == nil {
		item.pending = []Closeable{p}
		return
	}
	item.pending = append(item.pending, p)
}

func (item *AppState) Close() {
	if item.pending != nil {
		for i, p := range item.pending {
			log.Printf("Cleaning up (%v/%v)", i+1, len(item.pending))
			p.Close()
		}
	}
}

type serialized struct {
	Fixtures []string `json:fixtures`
	Mappings []string `json:mappings`
}

func (item *AppState) ToJSON() ([]byte, error) {
	serializedState := serialized{
		Fixtures: []string{},
		Mappings: []string{},
	}
	for _, fixture := range item.Fixtures {
		encoded := []string{
			strconv.Itoa(fixture.BaseAddress()),
			fixture.Name(),
		}
		for _, role := range fixture.Roles() {
			encoded = append(encoded, strconv.Itoa(role.EnumIndex()))
		}
		serializedState.Fixtures = append(serializedState.Fixtures, strings.Join(encoded, ":"))
	}
	for channel, channelConfig := range item.Mappings {
		for control, fixtureChannel := range channelConfig {
			serializedState.Mappings = append(serializedState.Mappings,
				strconv.Itoa(int(channel))+":"+
					strconv.Itoa(int(control))+":"+
					strconv.Itoa(fixtureChannel.Role().EnumIndex()),
			)
		}
	}
	return json.Marshal(serializedState)
}

func LoadShow(showJson []byte) (*AppState, error) {
	serializedState := serialized{
		Fixtures: []string{},
		Mappings: []string{},
	}
	err := json.Unmarshal(showJson, &serializedState)
	if err != nil {
		return nil, err
	}
	fixturesList := loadFixtures(serializedState.Fixtures)
	show := AppState{
		Fixtures: fixturesList,
		Mappings: loadMappings(serializedState.Mappings, fixturesList),
	}
	return &show, nil
}

func loadFixtures(config []string) []*fixtures.Fixture {
	fixturesList := []*fixtures.Fixture{}
	for i, entry := range config {
		parts := strings.Split(entry, ":")
		address, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			log.Printf("Invalid address in fixture %v: %v", i, entry)
			continue
		}
		name := parts[1]
		channels := []fixtures.ChannelRole{}
		isValid := true
		for j := 2; j < len(parts); j++ {
			offset, err := strconv.ParseInt(parts[j], 10, 32)
			if err != nil {
				log.Printf("Invalid offset %v in fixture %v: %v", j, i, entry)
				isValid = false
				break
			}
			channels = append(channels, fixtures.ChannelRole(offset))
		}
		if isValid {
			fixturesList = append(fixturesList, fixtures.NewFixture(int(address), name, channels...))
		} else {
			fixturesList = append(fixturesList, nil)
		}
	}
	return fixturesList
}

func loadMappings(mappings []string, fixturesList []*fixtures.Fixture) map[uint8]map[uint8]*fixtures.FixtureChannel {
	controlToFixtureChannels := map[uint8]map[uint8]*fixtures.FixtureChannel{}
	if len(fixturesList) == 0 {
		log.Println("No fixtures configured, skipping mappings loading")
		return controlToFixtureChannels
	}

	for i, mapping := range mappings {
		parts := strings.Split(mapping, ":")
		channel, err := strconv.ParseUint(parts[0], 10, 8)
		if err != nil {
			log.Printf("Invalid channel in mapping %v", i)
			continue
		}

		control, err := strconv.ParseUint(parts[1], 10, 8)
		if err != nil {
			log.Printf("Invalid control in mapping %v", i)
			continue
		}

		fixtureIndex, err := strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			log.Printf("Invalid fixture in mapping %v", i)
			continue
		}
		if int(fixtureIndex) >= len(fixturesList) {
			log.Printf("Unknown fixture in mapping %v: %v", i, fixtureIndex)
			continue
		}

		role, err := strconv.ParseInt(parts[3], 10, 32)
		if err != nil {
			log.Printf("Invalid role in mapping %v", i)
			continue
		}
		if !fixtures.IsChannelRole(int(role)) {
			log.Printf("Unknown role in mapping %v: %v", i, role)
			continue
		}

		channelConfig, ok := controlToFixtureChannels[uint8(channel)]
		if !ok || channelConfig == nil {
			channelConfig = map[uint8]*fixtures.FixtureChannel{}
			controlToFixtureChannels[uint8(channel)] = channelConfig
		}
		fixture := fixturesList[fixtureIndex]
		fixtureChannel := fixture.GetChannelForRole(fixtures.ChannelRole(int(role)))
		if fixtureChannel == nil {
			log.Printf("No Channel Role %v in fixture %v", fixtures.ChannelRole(int(role)).String(), role)
			continue
		}
		channelConfig[uint8(control)] = fixtureChannel
	}
	return controlToFixtureChannels
}
