package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	//Version is build number
	Version string

	isZoneConnected  = false
	isWorldConnected = false
	isLoginConnected = false
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
	fmt.Println("starting EQPort", Version)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	login, err := net.ListenPacket("udp", ":5999")
	if err != nil {
		return fmt.Errorf("listen loginserver on port 5999: %w", err)
	}
	go udpLoop(ctx, "login", login)
	defer login.Close()
	world, err := net.ListenPacket("udp", ":9000")
	if err != nil {
		return fmt.Errorf("listen world on port 9000: %w", err)
	}
	defer world.Close()
	go udpLoop(ctx, "world", world)
	zone, err := net.ListenPacket("udp", ":7000")
	if err != nil {
		return fmt.Errorf("listen zone on port 7000: %w", err)
	}
	defer zone.Close()
	go udpLoop(ctx, "zone", zone)
	fmt.Println("checking login, world, and zone ports...")

	tryPing()
	time.Sleep(3 * time.Second)
	isAllConnected := true
	if !isZoneConnected {
		fmt.Println("zone port 7000 is not open")
		isAllConnected = false
	}
	if !isWorldConnected {
		fmt.Println("world port 9000 is not open")
		isAllConnected = false
	}
	if !isLoginConnected {
		fmt.Println("login port 5999 is not open")
		isAllConnected = false
	}
	if isAllConnected {
		fmt.Println("world, zone, and login are all open")
	}
	os.Exit(0)
	return nil
}

func udpLoop(ctx context.Context, id string, conn net.PacketConn) {
	data := []byte{}
	for {
		data = make([]byte, 512)
		_, _, err := conn.ReadFrom(data)
		if err != nil {
			fmt.Printf("failed %s read: %s\n", id, err.Error())
			return
		}
		//fmt.Printf("external ip %s connected to %s successfully!\n", addr.String(), id)
		switch id {
		case "zone":
			isZoneConnected = true
		case "login":
			isLoginConnected = true
		case "world":
			isWorldConnected = true
		}
		return
	}
}

func tryPing() {
	host := "mandalorianquest.thegrandpackard.com:10000"
	conn, err := net.Dial("udp", host)
	if err != nil {
		fmt.Printf("failed to connect to wan host: %s\n", err.Error())
		os.Exit(1)
	}
	_, err = conn.Write([]byte("rawr"))
	if err != nil {
		fmt.Printf("failed to write to wan host: %s\n", err.Error())
		os.Exit(1)
	}
	conn.Close()
}
