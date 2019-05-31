package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	port := flag.String("p", "8000", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Printf("Current dir: " + pwd)

	ips, err := externalIPs()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Please visit: http://%s:%s\n", ips[0], *port)
	if len(ips) > 1 {
		for i := 1; i < len(ips); i++ {
			log.Printf("           or http://%s:%s\n", ips[i], *port)
		}
	}
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}

func externalIPs() ([]string, error) {
	ipArr := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return []string{}, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return []string{}, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ipArr = append(ipArr, ip.String())
		}
	}

	if len(ipArr) < 1 {
		return []string{}, errors.New("are you connected to the network?")
	} else {
		return ipArr, nil
	}
}