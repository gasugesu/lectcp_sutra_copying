package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gasugesu/lectcp/pkg/arp"
	"github.com/gasugesu/lectcp/pkg/ethernet"
	"github.com/gasugesu/lectcp/pkg/ip"
	"github.com/gasugesu/lectcp/pkg/net"
	"github.com/gasugesu/lectcp/pkg/raw/tuntap"
)

var sig chan os.Signal

func init(){
	arp.Init()
}

func setup() (*net.Device, error) {
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// parse command line params
	name := flag.String("name", "", "device name")
	addr := flag.String("addr", "", "hardware address")
	flag.Parse()
	raw, err := tuntap.NewTap(*name)
	if err != nil {
		return nil, err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return nil, err
	}
	if *addr != "" {
		link.SetAddress(ethernet.ParseAddress(*addr))
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return nil, err
	}
	iface, err := ip.CreateInterface(dev, "192.0.2.1", "255.255.255.0", "")
	if err != nil {
		return nil, err
	}
	dev.RegisterInterface(iface)
	return dev, nil
}

func main() {
	dev, err := setup()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] %s\n", dev.Name(), dev.Address())
	select {
	case s := <-sig:
		fmt.Printf("sig: %s\n", s)
		dev.Shutdown()
	// TODO: create PR
	case err = <-dev.Errors:
		dev.Shutdown()
	}
	fmt.Println("good bye")
}
