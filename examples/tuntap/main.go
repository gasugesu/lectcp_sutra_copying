package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gasugesu/lectcp/pkg/raw"
	"github.com/gasugesu/lectcp/pkg/raw/tuntap"
)

func setup() (raw.Device, error) {
	// sig handling
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// parse command line params
	name := flag.String("name", "", "device name")
	flag.Parse()
	dev, err := tuntap.NewTap(*name)
	if err != nil {
		return nil, err
	}
	go func() {
		s := <-sig
		fmt.Printf("sig: %s\n", s)
		dev.Close()
	}()
	return dev, nil
}

func main() {
	dev, err := setup()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] %x\n", dev.Name(), dev.Address())
	buf := make([]byte, 1024*4)
	for {
		n, err := dev.Read(buf)
		if n > 0 {
			fmt.Printf("--- [%s] %d bytes ---\n", dev.Name(), n)
			fmt.Printf("%s\n", hex.Dump(buf[:n]))
		}
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
	fmt.Println("good bye")
}
