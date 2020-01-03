package home

import (
	"fmt"
	"strconv"

	"github.com/eriklupander/tradfri-go/tradfri"
)

type Action struct {
	Value         bool
	ValidateValue func(string) bool
}

func GetActions() map[string]Action {
	actions := make(map[string]Action)

	noValidate := Action{
		Value: false,
	}
	actions["on"] = noValidate
	actions["off"] = noValidate
	actions["toggle"] = noValidate
	actions["dim"] = Action{
		Value: true,
		ValidateValue: func(s string) bool {
			i, err := strconv.Atoi(s)
			if err != nil {
				return false
			}

			if i < 1 || i > 254 {
				return false
			}

			return true
		},
	}

	return actions
}

func actionAllLamps(client *tradfri.TradfriClient, lamps []Lamp, action, value string) {
	for _, l := range lamps {
		switch action {
		case "on":
			l.On(client)
		case "off":
			l.Off(client)
		case "toggle":
			l.Toggle(client)
		case "dim":
			l.Dim(client, value)
		}
	}
}

func ExecAction(client *tradfri.TradfriClient, room, lamp, action, value string) error {
	rooms := GetRooms(client)

	actions := GetActions()
	if _, ok := actions[action]; !ok {
		return fmt.Errorf("Invalid action %s", action)
	}

	if actions[action].Value {
		if value == "" {
			return fmt.Errorf("No value for action %s", action)
		}

		if actions[action].ValidateValue(value) == false {
			return fmt.Errorf("Invalid value %s for action %s", value, action)
		}
	}

	if room == "all" {
		for _, r := range rooms {
			actionAllLamps(client, r.Lamps, action, value)
		}

		return nil
	}

	if _, ok := rooms[room]; !ok {
		return fmt.Errorf(
			fmt.Sprintf("No room named %s", room),
		)
	}

	if lamp == "all" {
		actionAllLamps(client, rooms[room].Lamps, action, value)
		return nil
	}

	var targetLamp Lamp
	foundLamp := false

	for _, l := range rooms[room].Lamps {
		if l.Name != lamp {
			continue
		}
		targetLamp = l
		foundLamp = true
		break
	}

	if foundLamp == false {
		return fmt.Errorf(
			fmt.Sprintf("No lamp named %s in room %s", lamp, room),
		)
	}

	actionAllLamps(client, []Lamp{targetLamp}, action, value)
	return nil
}
