// http://stackoverflow.com/questions/22930510/how-to-retrieve-address-of-current-machine

package main

import (
	"fmt"
	"net"
	"os"
)

/*
https://github.com/mccoyst/myip/blob/master/myip.go
(C) 2012 Steve McCoy. Licensed under the MIT license.
The myip command prints all non-loopback IP addresses associated
with the machine that it runs on, one per line.
*/
func myip() {
	os.Stdout.WriteString("myip:\n")

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Errorf("error: %v\n", err.Error())
		return
	}

	for _, a := range addrs {
		ip := net.ParseIP(a.String())
		fmt.Printf("addr: %v loopback=%v\n", a, ip.IsLoopback())
	}

	fmt.Println()
}

func myip2() {
	os.Stdout.WriteString("myip2:\n")

	tt, err := net.Interfaces()
	if err != nil {
		fmt.Errorf("error: %v\n", err.Error())
		return
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			fmt.Errorf("error: %v\n", err.Error())
			continue
		}
		for _, a := range aa {
			ip := net.ParseIP(a.String())
			fmt.Printf("%v addr: %v loopback=%v\n", t.Name, a, ip.IsLoopback())
		}
	}

	fmt.Println()
}

func main() {
	fmt.Println("myip -- begin")
	myip()
	myip2()
	fmt.Println("myip -- end")
}
