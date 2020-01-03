package home

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/eriklupander/tradfri-go/model"
	"github.com/eriklupander/tradfri-go/tradfri"
)

type Lamp struct {
	Name   string
	ID     string
	Type   string
	Vendor string
	Power  bool
	Xcolor int
	Ycolor int
	RGB    string
	Dimmer int
}

func (l Lamp) On(client *tradfri.TradfriClient) {
	if l.Power == false {
		client.PutDevicePower(l.ID, 1)
	}
}

func (l Lamp) Off(client *tradfri.TradfriClient) {
	if l.Power {
		client.PutDevicePower(l.ID, 0)
	}
}

func (l Lamp) Toggle(client *tradfri.TradfriClient) {
	if l.Power {
		l.Off(client)
	} else {
		l.On(client)
	}
}

func (l Lamp) Dim(client *tradfri.TradfriClient, value string) {
	i, err := strconv.Atoi(value)

	if err != nil {
		log.Fatalln("Dim value is not an int")
	}

	if i < 1 || i > 254 {
		log.Fatalln("Dim value must be between 1-254")
	}

	client.PutDeviceDimming(l.ID, i)
}

func formatLamp(lamp model.Device) Lamp {
	return Lamp{
		Name:   strings.ToLower(lamp.Name),
		ID:     fmt.Sprintf("%d", lamp.DeviceId),
		Type:   lamp.Metadata.TypeName,
		Vendor: lamp.Metadata.Vendor,
		Power:  lamp.LightControl[0].Power == 1,
		Xcolor: lamp.LightControl[0].CIE_1931_X,
		Ycolor: lamp.LightControl[0].CIE_1931_Y,
		RGB:    lamp.LightControl[0].RGBHex,
		Dimmer: lamp.LightControl[0].Dimmer,
	}
}

func GetLamps(client *tradfri.TradfriClient, ids []int) []Lamp {
	var ll []Lamp

	for _, dID := range ids {
		dres, err := client.GetDevice(fmt.Sprintf("%d", dID))

		if err != nil {
			log.Fatalln(err)
		}

		//Skip non lamps
		if dres.LightControl == nil || len(dres.LightControl) == 0 {
			continue
		}

		l := formatLamp(dres)

		ll = append(ll, l)
	}

	return ll
}
