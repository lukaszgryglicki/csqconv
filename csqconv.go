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
	// Execution time
	dtStart := time.Now()
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
		dtEnd := time.Now()
		fmt.Printf("%s: %+v\n", info, dtEnd.Sub(dtStart))
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

func processFrame(root string, frameNo int, debug int, output bool, norm bool) error {
	// ffmpeg -f image2 -vcodec jpegls -i "$f" -y -f image2 -vcodec png "$f2"
	root = fmt.Sprintf("%s%06d", root, frameNo)
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
	if !norm {
		res, err = execCommand(
			debug,
			output,
			[]string{"rm", "-f", ifn},
			nil,
		)
		if err != nil {
			if res != "" {
				fmt.Printf("rm %s:\n%s\n", ifn, res)
			}
			return err
		}
	}
	return nil
}

func processCsqFile(fn string, minFrames int) error {
	mode := os.Getenv("MODE")
	crf := ""
	mpng := "mpng"
	if mode != "" && mode != mpng {
		crf = os.Getenv("CRF")
		if crf != "" {
			icrf, err := strconv.Atoi(os.Getenv("CRF"))
			if err != nil {
				return err
			}
			if icrf < 0 || icrf > 51 {
				return fmt.Errorf("crf must be from 0-51 range, got %d", icrf)
			}
		} else {
			crf = "17"
		}
	}
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
	fmt.Printf("%s: %d frames\n", fn, nAry)

	debug := 0
	if os.Getenv("DEBUG") != "" {
		debug, err = strconv.Atoi(os.Getenv("DEBUG"))
		if err != nil {
			return err
		}
	}
	output := os.Getenv("OUTPUT") != ""
	norm := os.Getenv("NORM") != ""
	hint := os.Getenv("HINT") != ""

	sr := 0
	if os.Getenv("SR") != "" {
		sr, err = strconv.Atoi(os.Getenv("SR"))
		if err != nil {
			return err
		}
		if sr < 2 {
			return fmt.Errorf("SR must be >= 2: %d", sr)
		}
	}

	var ifn string
	ext := []string{".raw", ".jpegls"}
	jpegLS := []byte("\xff\xd8\xff\xf7")
	hdr := [][]byte{[]byte(""), jpegLS}
	indices := []int{}
	for i, fdata := range ary {
		if i > 0 && i%100 == 99 {
			fmt.Printf("%s: frame %d/%d\n", fn, i+1, nAry)
		}
		ifn = fmt.Sprintf("%s%06d", root, i)
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
			if !norm && k == 0 {
				continue
			}
			err = ioutil.WriteFile(ifn+ext[k], append(hdr[k], idata...), 0644)
			if err != nil {
				return err
			}
		}
		err = processFrame(root, i, debug, output, norm)
		if err != nil {
			return err
		}
		indices = append(indices, i)
	}

	nIndices := len(indices)
	h := 0x1fff9 / (len(root) + 14)
	packs := nIndices / h
	if nIndices%h > 0 {
		packs++
	}

	// Postprocess all frames in h-sized packs
	indicesA := [][]int{}
	for p := 0; p < packs; p++ {
		f := h * p
		t := f + h
		if t > nIndices {
			t = nIndices
		}
		indicesA = append(indicesA, indices[f:t])
	}

	for i, ind := range indicesA {
		if sr > 1 {
			fmt.Printf("Postprocessing %d/%d pack: %d frames (sr command)\n", i+1, packs, len(ind))
			cmd := []string{"sr", fmt.Sprintf("%d", sr)}
			for _, idx := range ind {
				cmd = append(cmd, fmt.Sprintf("%s%06d.png", root, idx))
			}
			res, err := execCommand(debug, output, cmd, map[string]string{"GS": "1", "INPL": "1"})
			if err != nil {
				if res != "" {
					fmt.Printf("postprocessing frames via 'sr:%d' tool:\n%s\n", sr, res)
				}
				return err
			}
		}
		if hint {
			fmt.Printf("Postprocessing %d/%d pack: %d frames (hist command)\n", i+1, packs, len(ind))
			cmd := []string{"hist"}
			for _, idx := range ind {
				cmd = append(cmd, fmt.Sprintf("%s%06d.png", root, idx))
			}
			res, err := execCommand(debug, output, cmd, nil)
			if err != nil {
				if res != "" {
					fmt.Printf("postprocessing frames via 'hist' tool:\n%s\n", res)
				}
				return err
			}
		}
		fmt.Printf("Postprocessing %d/%d pack: %d frames (jpeg command)\n", i+1, packs, len(ind))
		cmd := []string{"jpeg"}
		for _, idx := range ind {
			cmd = append(cmd, fmt.Sprintf("%s%06d.png", root, idx))
		}
		res, err := execCommand(debug, output, cmd, nil)
		if err != nil {
			if res != "" {
				fmt.Printf("postprocessing frames via 'jpeg' tool:\n%s\n", res)
			}
			return err
		}

		// Remove intermediate PNGs
		if !norm {
			cmdRm := []string{"rm", "-f"}
			cmdRm = append(cmdRm, cmd[1:]...)
			res, err = execCommand(debug, output, cmdRm, nil)
			if err != nil {
				if res != "" {
					fmt.Printf("rm intermediate PNGs:\n%s\n", res)
				}
				return err
			}
			cmdRm = []string{"rm", "-f"}
			for _, idx := range ind {
				cmdRm = append(cmdRm, fmt.Sprintf("%s%06d.png.hint", root, idx))
			}
			res, err = execCommand(debug, output, cmdRm, nil)
			if err != nil {
				if res != "" {
					fmt.Printf("rm intermediate Hints:\n%s\n", res)
				}
				return err
			}
		}
	}

	// Create final X.264 MP4
	// To check frames info and encoded info
	// ffmpeg -i "co_small%06d.png" -f framehash -
	// ffmpeg -i small.mp4 -map 0:v -f framehash -
	// ffmpeg -f image2 -vcodec png -r 30 -i "co_small%06d.png" -q:v 0 -y -vcodec libx264 -y "small.mp4"
	// Truly lossless mpng (keeps 16 bit color data - each frama checksum match original input PNGs) - but takes 150x more space
	// ffmpeg -f image2 -vcodec png -r 30 -i "co_small%06d.png" -codec copy -y "small.mp4"
	// Truly lossless H.264 - slow and still 50x times more space than -q:v 0
	// ffmpeg -f image2 -vcodec png -r 30 -i "co_small%06d.png" -vcodec libx264 -crf 0 -preset veryslow -y "small.mp4"
	fmt.Printf("Creating final video\n")
	a := strings.Split(root, "/")
	lA := len(a)
	last := a[lA-1]
	a[lA-1] = "co_" + last
	nroot := strings.Join(a, "/")
	//pattern := nroot + "%06d.png"
	pattern := nroot + "*.png"
	vidfn := root + ".mp4"
	vCmd := []string{}
	if mode == "" {
		// ffmpeg -f image2 -vcodec png -r 30 -i "co_small%06d.png" -q:v 0 -y -vcodec libx264 -y "small.mp4"
		vCmd = []string{
			"ffmpeg", "-f", "image2", "-pattern_type", "glob",
			"-vcodec", "png", "-r", "30", "-i", pattern,
			"-q:v", "0", "-vcodec", "libx264", "-y", vidfn,
		}
	} else if mode == mpng {
		// ffmpeg -f image2 -vcodec png -r 30 -i "co_small%06d.png" -codec copy -y "small.mp4"
		vCmd = []string{
			"ffmpeg", "-f", "image2", "-pattern_type", "glob",
			"-vcodec", "png", "-r", "30", "-i", pattern,
			"-codec", "copy", "-y", vidfn,
		}
	} else {
		// ffmpeg -f image2 -vcodec png -r 30 -i "co_small%06d.png" -vcodec libx264 -crf 0 -preset veryslow -y "small.mp4"
		vCmd = []string{
			"ffmpeg", "-f", "image2", "-pattern_type", "glob",
			"-vcodec", "png", "-r", "30", "-i", pattern,
			"-vcodec", "libx264", "-crf", crf,
			"-preset", mode, "-y", vidfn,
		}
	}
	res, err := execCommand(debug, output, vCmd, nil)
	if err != nil {
		if res != "" {
			fmt.Printf("%s -> %s:\n%s\n", pattern, vidfn, res)
		}
		return err
	}

	// Cleanup postprocessed frames
	if !norm {
		rmpattern := nroot + "*.png"
		for _, ind := range indicesA {
			cmdRm := []string{"rm", "-f"}
			for _, idx := range ind {
				cmdRm = append(cmdRm, fmt.Sprintf("%s%06d.png", nroot, idx))
			}
			res, err = execCommand(debug, output, cmdRm, nil)
			if err != nil {
				if res != "" {
					fmt.Printf("rm %s:\n%s\n", rmpattern, res)
				}
				return err
			}
		}
	}
	return nil
}

func help() {
	helpStr := `
Environment variables:
MIN_FRAMES - minimum number of frames that must be present in CSQ file (default 5)
DEBUG - enabled debug mode (default 0)
OUTPUT - enabled additional output
NORM - do not remove temporary files used
HINT - use hint program to generate moving live histogram data and use it to create video
  this should make flashing less visible (histogram calculated over 16 frames by default)
MODE - set video encoding mode: if empty then uses '-q:v 0', else:
  mpng (loseless PNG sequence), libx264 preset name otherwise (CRF mode):
  ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow
CRF - set quantizer parameter: 0-51 when running in CRF mode, default 17
SR - set scale factor (2 will combine 4 images, 3 will combine 9 images, N will combine N^2 images)
`
	fmt.Printf("%s\n", helpStr)
}

// Env:
// MIN_FRAMES (default 5)
// DEBUG (default 0)
// OUTPUT (default false)
// NORM (default false)
// HINT (default false)
func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Need at least one file name\n")
		help()
		return
	}
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
