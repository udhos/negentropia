package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path/filepath"
)

func main() {

	if len(os.Args) != 4 {
		basename := filepath.Base(os.Args[0])
		log.Fatalf("usage: %s hostname username password", basename)
	}

	host := os.Args[1]
	user := os.Args[2]
	pass := os.Args[3]

	// Create client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
	}
	// Connect to ssh server
	log.Printf("** opening: %s", host)
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer conn.Close()
	log.Printf("** connected: %s", host)
	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}
	defer session.Close()
	log.Printf("** session open")
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}
	log.Printf("** pseudo-terminal ready")

	/*
		// Start remote shell
		if err := session.Shell(); err != nil {
			log.Fatalf("failed to start shell: %s", err)
		}
		log.Printf("** remote shell ready")
	*/

	cmd := "ls -la"
	log.Printf("** running command: %s", cmd)
	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
