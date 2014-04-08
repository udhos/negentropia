// http://stackoverflow.com/questions/22930510/how-to-retrieve-address-of-current-machine

package main

import (
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
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			os.Stdout.WriteString(ipnet.IP.String())
			if ipnet.IP.IsLoopback() {
				os.Stdout.WriteString(" loopback\n")
			} else {
				os.Stdout.WriteString(" non-loopback\n")
			}
		}
	}
}

func myip2() {
	os.Stdout.WriteString("myip2:\n")

	tt, err := net.Interfaces()
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			os.Stdout.WriteString(err.Error())
			continue
		}
		for _, a := range aa {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				os.Stdout.WriteString(t.Name + " " + v4.String() + " loopback\n")
			} else {
				os.Stdout.WriteString(t.Name + " " + v4.String() + " non-loopback\n")
			}
		}
	}
}

func main() {
	os.Stdout.WriteString("myip -- begin\n")
	myip()
	myip2()
	os.Stdout.WriteString("myip -- end\n")	
}
