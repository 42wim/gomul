package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var iMap map[string]int
var flagList bool

func parseInterfaces() {
	if flagList {
		fmt.Println("Available interfaces")
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			//windows
			case *net.IPAddr:
				iMap[v.String()] = i.Index
				if flagList {
					fmt.Printf("%v) %v : %s\n", i.Index, i.Name, v.String())
				}
			//linux
			case *net.IPNet:
				iMap[v.String()] = i.Index
				if flagList {
					fmt.Printf("%v) %v : %s\n", i.Index, i.Name, v.String())
				}
			}

		}
	}
}

func main() {
	var flagGroup, flagIP string
	var flagIndex int
	var c net.PacketConn
	var p6 *ipv6.PacketConn
	var p4 *ipv4.PacketConn

	iMap = make(map[string]int)

	flag.BoolVar(&flagList, "li", false, "show available interfaces")
	flag.IntVar(&flagIndex, "interface", 0, "interface to listen on (number)")
	flag.StringVar(&flagGroup, "group", "ff02::42:1 239.42.42.1", "multicast groups to join (space seperated)")
	flag.StringVar(&flagIP, "ip", "", "use interface where the specified ip is bound on")
	flag.Parse()
	if flag.NFlag() < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if flagList {
		parseInterfaces()
		os.Exit(1)
	}

	parseInterfaces()
	if flagIndex == 0 {
		fmt.Print("searching interface for ip ", flagIP)
		for k, v := range iMap {
			if strings.HasPrefix(k, flagIP) {
				fmt.Println(" using interface", v, "with ip", k)
				flagIndex = v
				break
			}
		}
	}

	fmt.Println("listening on index", flagIndex)
	iface, err := net.InterfaceByIndex(flagIndex)
	if err != nil {
		log.Fatal(err)
	}

	groups := strings.Fields(flagGroup)

	if strings.Contains(flagGroup, ":") {
		c, err = net.ListenPacket("udp6", "[::]:1024")
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		p6 = ipv6.NewPacketConn(c)
	}
	if strings.Contains(flagGroup, ".") {
		c, err = net.ListenPacket("udp4", ":1024")
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		p4 = ipv4.NewPacketConn(c)
	}

	for _, group := range groups {
		IPgroup := net.ParseIP(group)
		if strings.Contains(group, ":") {
			fmt.Println("joining ipv6 group", group)
			if err := p6.JoinGroup(iface, &net.UDPAddr{IP: IPgroup}); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("joining ipv4 group", group)
			if err := p4.JoinGroup(iface, &net.UDPAddr{IP: IPgroup}); err != nil {
				log.Fatal(err)
			}
		}
	}
	select {}
}
