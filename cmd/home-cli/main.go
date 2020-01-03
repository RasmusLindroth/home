package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/RasmusLindroth/home/grpchome/golang"
	"github.com/RasmusLindroth/home/home"
	"google.golang.org/grpc"
)

func showHelp() {
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tlamps <room> [<lamp>] <action> [value]\n\n")
	fmt.Printf("\t<room> = room name or all (if all, skip <lamp>\n")
	fmt.Printf("\t<lamp> = lamp name or all\n")
	fmt.Printf("\t<action> =\n")
	fmt.Printf("\t\tOn\n")
	fmt.Printf("\t\tOff\n")
	fmt.Printf("\t\tToggle\n")
	fmt.Printf("\t\tDim = [1-254]\n")
}

func parseFlags() (room, lamp, action, value string, err error) {
	if len(os.Args) < 3 {
		return room, lamp, action, value, errors.New("Not enough argumetns")
	}

	actions := home.GetActions()

	room = strings.ToLower(os.Args[1])
	i := 2

	if room == "all" {
		lamp = "all"
	}

	if lamp == "" {
		lamp = strings.ToLower(os.Args[2])
		i++
	} else {
		action = strings.ToLower(os.Args[2])
		i++
	}

	if action == "" {
		if i+1 > len(os.Args) {
			return room, lamp, action, value, errors.New("Not enough argumetns")
		}

		i++
		action = strings.ToLower(os.Args[3])
	}

	if _, ok := actions[action]; !ok {
		return room, lamp, action, value, fmt.Errorf("Invalid action %s", action)
	}

	if actions[action].Value {
		if i+1 > len(os.Args) {
			return room, lamp, action, value, fmt.Errorf("No value for action %s", action)
		}

		value = os.Args[i]
		i++
		if actions[action].ValidateValue(value) == false {
			return room, lamp, action, value, fmt.Errorf("Invalid value %s for action %s", value, action)
		}
	}

	if i != len(os.Args) {
		return room, lamp, action, value, errors.New("Invalid number of arguments")
	}
	return room, lamp, action, value, nil
}

func cli(address string) {
	room, lamp, action, value, err := parseFlags()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		showHelp()
		return
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*10))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHomeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.RunAction(ctx, &pb.RunActionRequest{
		Room:   room,
		Lamp:   lamp,
		Action: action,
		Value:  value,
	})
	if err != nil {
		fmt.Println(grpc.ErrorDesc(err))
	}
}

func main() {
	conf, err := home.ParseConfig()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	cli(conf.GetAddress())
}
