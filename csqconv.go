package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func processCsqFile(fn string, minFrames int) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	ary := bytes.Split(data, []byte("\x46\x46\x46\x00\x52\x54"))
	nAry := len(ary)
	if nAry < minFrames {
		fmt.Printf("%s: no frames\n", fn)
		return nil
	}
	fmt.Printf("%s: %d frames\n", fn, nAry)
	return nil
}

func main() {
	minFrames := os.Getenv("MIN_FRAMES")
	if minFrames == "" {
		minFrames = "2"
	}
	mf, err := strconv.Atoi(minFrames)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
	for _, arg := range os.Args[1:] {
		dtStart := time.Now()
		err := processCsqFile(arg, mf)
		if err != nil {
			fmt.Printf("%s: error: %+v\n", arg, err)
		}
		dtEnd := time.Now()
		fmt.Printf("%s: %v\n", arg, dtEnd.Sub(dtStart))
	}
}
