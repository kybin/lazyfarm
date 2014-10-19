package main

import (
	"fmt"
	"net"
	"log"
	"flag"
	"strings"
	"strconv"
	"errors"
	"sort"
	"encoding/gob"
	"os"
)

func main() {
	var cmd string
	var framestr string
	var server string
	var group string
	flag.StringVar(&cmd, "cmd", "", "render command")
	flag.StringVar(&framestr, "frames", "", "frames for render")
	flag.StringVar(&server, "server", "", "server address")
	flag.StringVar(&group, "group", "", "worker in the group will serve this job")
	flag.Parse()
	if server == "" {
		fmt.Println("please specify server address")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// expand scene file path
	//if scene != "" {
	//	usr, _ := user.Current()
	//	home := usr.HomeDir
	//	if scene[:2] == "~/" {
	//		scene = strings.Replace(scene, "~", home, 1)
	//	}
	//}
	frames, err := parseFrames(framestr)
	fmt.Println(frames)
	if err != nil {
		log.Fatal(err)
	}
	r := &Job {
		Cmd : cmd,
		Frames : frames,
		Group : group,
	}
	fmt.Println(r)
	conn, err := net.Dial("tcp", server)
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

// Here we parse frames flag. The result will list of frames. 
// Any ambiguity in the flag leads error.
func parseFrames(framestr string) ([]int, error) {
	frames := make([]int, 0)
	if framestr == "" {
		err := errors.New("Cannot parse empty frame")
		return nil, err
	}
	splited := strings.Split(framestr, ",")
	for _, f := range splited {
		f = strings.TrimSpace(f)
		if strings.Contains(f, "-") {
			fs := strings.Split(f, "-")
			if len(fs) != 2 {
				err := errors.New(fmt.Sprintf("Cannot parse frames flag : %v", framestr))
				return nil, err
			}
			fstart, err := strconv.Atoi(fs[0])
			if err != nil {
				return nil, err
			}
			fend, err := strconv.Atoi(fs[1])
			if err != nil {
				return nil, err
			}
			if fstart >= fend {
				err := errors.New(fmt.Sprintf("Cannot parse frames flag : %v", framestr))
				return nil, err
			}
			for i := fstart; i <= fend; i++ {
				frames = append(frames, i)
			}
		} else {
			i, err := strconv.Atoi(f)
			if err != nil {
				return nil, err
			}
			frames = append(frames, i)
		}
	}
	sort.Sort(sort.IntSlice(frames))
	//removeDuplicates(frames)
	return frames, nil
}

