package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func processCsqFile(fn string, minFrames int) error {
	var err error
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
	fna := strings.Split(fn, ".")
	root := strings.Join(fna[0:len(fna)-1], ".")
	fmt.Printf("%s: %d frames --> %s_nnnnnnnn.jpg\n", fn, nAry, root)
	var ifn string
	for i, fdata := range ary {
		ifn = fmt.Sprintf("%s_%08d.jpg", root, i)
		err = ioutil.WriteFile(ifn, fdata, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var err error
	minFrames := os.Getenv("MIN_FRAMES")
	if minFrames == "" {
		minFrames = "5"
	}
	mf, err := strconv.Atoi(minFrames)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
	for _, arg := range os.Args[1:] {
		dtStart := time.Now()
		err = processCsqFile(arg, mf)
		if err != nil {
			fmt.Printf("%s: error: %+v\n", arg, err)
		}
		dtEnd := time.Now()
		fmt.Printf("%s: %v\n", arg, dtEnd.Sub(dtStart))
	}
}
