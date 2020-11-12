package main

import (
	"context"
	"fmt"
	"net"
	"os"
)

var (
	//Version is build number
	Version string
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("failed:", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	fmt.Println("Starting EQPort", Version)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//5999
	//9000
	//7000

	login, err := net.ListenPacket("udp", ":5999")
	if err != nil {
		return fmt.Errorf("listen loginserver on port 5999: %w", err)
	}
	defer login.Close()
	world, err := net.ListenPacket("udp", ":9000")
	if err != nil {
		return fmt.Errorf("listen world on port 9000: %w", err)
	}
	defer world.Close()
	zone, err := net.ListenPacket("udp", ":7000")
	if err != nil {
		return fmt.Errorf("listen zone on port 7000: %w", err)
	}
	defer zone.Close()
	fmt.Println("listening on udp 5999, 9000, and 7000.")
	fmt.Println("visit https://portchecker.co/ and put the ports in, see if it says open")
	fmt.Println("just exit this window to close the ports")
	select {
	case <-ctx.Done():
		return nil
	}

}
