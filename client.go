package main

import (
	"net"
	"log"
	"flag"
	"strings"
	"os/user"
	"encoding/gob"
)

type R struct {
	Run, Scene, Driver, Frames string
}

func main() {
	var run string
	var scene string
	var driver string
	var frames string
	flag.StringVar(&run, "run", "", "a program you want render.")
	flag.StringVar(&scene, "scene", "", "scene for render")
	flag.StringVar(&driver, "driver", "", "which node (if exists) for render")
	flag.StringVar(&frames, "frames", "", "frames for render")
	flag.Parse()
	// expand scene file path
	if scene != "" {
		usr, _ := user.Current()
		home := usr.HomeDir
		if scene[:2] == "~/" {
			scene = strings.Replace(scene, "~", home, 1)
		}
	}
	r := &R {
		Run : run,
		Scene : scene,
		Driver : driver,
		Frames : frames,
	}
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(r)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
