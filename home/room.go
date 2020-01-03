package home

import (
	"fmt"
	"log"
	"strings"

	"github.com/eriklupander/tradfri-go/tradfri"
)

type Room struct {
	Name  string
	ID    string
	Lamps []Lamp
}

func (r Room) On(client *tradfri.TradfriClient) {
	for _, l := range r.Lamps {
		l.On(client)
	}
}

func (r Room) Off(client *tradfri.TradfriClient) {
	for _, l := range r.Lamps {
		l.Off(client)
	}
}

func GetRooms(client *tradfri.TradfriClient) map[string]Room {
	groups, err := client.ListGroups()
	if err != nil {
		log.Fatalln(err)
	}

	rooms := make(map[string]Room)
	for _, g := range groups {
		r := Room{
			Name: strings.ToLower(g.Name),
			ID:   fmt.Sprintf("%d", g.DeviceId),
		}

		r.Lamps = GetLamps(client, g.Content.DeviceList.DeviceIds)

		rooms[r.Name] = r
	}

	return rooms
}
