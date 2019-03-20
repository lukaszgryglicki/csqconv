package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func execCommand(debug int, output bool, cmdAndArgs []string, env map[string]string) (string, error) {
	// STDOUT pipe size
	pipeSize := 0x100

	// Command & arguments
	command := cmdAndArgs[0]
	arguments := cmdAndArgs[1:]
	if debug > 0 {
		var args []string
		for _, arg := range cmdAndArgs {
			argLen := len(arg)
			if argLen > 0x200 {
				arg = arg[0:0x100] + "..." + arg[argLen-0x100:argLen]
			}
			if strings.Contains(arg, " ") {
				args = append(args, "'"+arg+"'")
			} else {
				args = append(args, arg)
			}
		}
		fmt.Printf("%s\n", strings.Join(args, " "))
	}
	cmd := exec.Command(command, arguments...)

	// Environment setup (if any)
	if len(env) > 0 {
		newEnv := os.Environ()
		for key, value := range env {
			newEnv = append(newEnv, key+"="+value)
		}
		cmd.Env = newEnv
		if debug > 0 {
			fmt.Printf("Environment Override: %+v\n", env)
			if debug > 2 {
				fmt.Printf("Full Environment: %+v\n", newEnv)
			}
		}
	}

	// Capture STDOUT (non buffered - all at once when command finishes), only used on error and when no buffered/piped version used
	// Which means it is used on error when debug <= 1
	// In debug > 1 mode, we're displaying STDOUT during execution, and storing results to 'outputStr'
	// Capture STDERR (non buffered - all at once when command finishes)
	var (
		stdOut    bytes.Buffer
		stdErr    bytes.Buffer
		outputStr string
	)
	cmd.Stderr = &stdErr
	if debug <= 1 {
		cmd.Stdout = &stdOut
	}

	// Pipe command's STDOUT during execution (if debug > 1)
	// Or just starts command when no STDOUT debug
	if debug > 1 {
		stdOutPipe, e := cmd.StdoutPipe()
		if e != nil {
			return "", e
		}
		e = cmd.Start()
		if e != nil {
			return "", e
		}
		buffer := make([]byte, pipeSize, pipeSize)
		nBytes, e := stdOutPipe.Read(buffer)
		for e == nil && nBytes > 0 {
			fmt.Printf("%s", buffer[:nBytes])
			outputStr += string(buffer[:nBytes])
			nBytes, e = stdOutPipe.Read(buffer)
		}
		if e != io.EOF {
			return "", e
		}
	} else {
		e := cmd.Start()
		if e != nil {
			return "", e
		}
	}
	// Wait for command to finish
	err := cmd.Wait()

	// If error - then output STDOUT, STDERR and error info
	if err != nil {
		if debug <= 1 {
			outStr := stdOut.String()
			if len(outStr) > 0 {
				fmt.Printf("%v\n", outStr)
			}
		}
		errStr := stdErr.String()
		if len(errStr) > 0 {
			fmt.Printf("STDERR:\n%v\n", errStr)
		}
		if err != nil {
			return stdOut.String(), err
		}
	}

	// If debug > 1 display STDERR contents as well (if any)
	if debug > 1 {
		errStr := stdErr.String()
		if len(errStr) > 0 {
			fmt.Printf("Errors:\n%v\n", errStr)
		}
	}
	if debug > 0 {
		info := strings.Join(cmdAndArgs, " ")
		lenInfo := len(info)
		if lenInfo > 0x280 {
			info = info[0:0x140] + "..." + info[lenInfo-0x140:lenInfo]
		}
	}
	outStr := ""
	if output {
		if debug <= 1 {
			outStr = stdOut.String()
		} else {
			outStr = outputStr
		}
	}
	return outStr, nil
}

func processFrame(root string, frameNo int, debug int, output bool) error {
	// ffmpeg -f image2 -vcodec jpegls -i "$f" -y -f image2 -vcodec png "$f2"
	root = fmt.Sprintf("%s_%08d", root, frameNo)
	ifn := root + ".jpegls"
	ofn := root + ".png"
	res, err := execCommand(
		debug,
		output,
		[]string{
			"ffmpeg", "-f", "image2",
			"-vcodec", "jpegls", "-i",
			ifn, "-y", "-f", "image2",
			"-vcodec", "png", ofn,
		},
		nil,
	)
	if err != nil {
		if res != "" {
			fmt.Printf("%s -> %s:\n%s\n", ifn, ofn, res)
		}
		return err
	}
	return nil
}

func processCsqFile(fn string, minFrames int) error {
	var err error
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	ary := bytes.Split(data, []byte("\x46\x46\x46\x00\x52\x54"))
	ary = ary[1:]
	nAry := len(ary)
	if nAry < minFrames {
		fmt.Printf("%s: no frames\n", fn)
		return nil
	}
	fna := strings.Split(fn, ".")
	root := strings.Join(fna[0:len(fna)-1], ".")
	fmt.Printf("%s: %d frames --> %s_nnnnnnnn.jpegls\n", fn, nAry, root)

	debug := 0
	if os.Getenv("DEBUG") != "" {
		debug, err = strconv.Atoi(os.Getenv("DEBUG"))
		if err != nil {
			return err
		}
	}
	output := os.Getenv("OUTPUT") != ""

	var ifn string
	ext := []string{".raw", ".jpegls"}
	jpegLS := []byte("\xff\xd8\xff\xf7")
	hdr := [][]byte{[]byte(""), jpegLS}
	for i, fdata := range ary {
		ifn = fmt.Sprintf("%s_%08d", root, i)
		iary := bytes.Split(fdata, jpegLS)
		liary := len(iary)
		if liary != 2 {
			fmt.Printf("%s: broken frame\n", ifn)
			err = ioutil.WriteFile(ifn+".err", fdata, 0644)
			if err != nil {
				return err
			}
			continue
		}
		for k, idata := range iary {
			err = ioutil.WriteFile(ifn+ext[k], append(hdr[k], idata...), 0644)
			if err != nil {
				return err
			}
		}
		err = processFrame(root, i, debug, output)
		if err != nil {
			return err
		}
	}
	return nil
}

// Env:
// MIN_FRAMES (default 5)
// DEBUG (default 0)
// OUTPUT (default false)
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
