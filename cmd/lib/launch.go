package lib

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Param struct {
	Url, Title string
}

var saveDir string

func GenerateCmd(param Param) (*exec.Cmd, io.Reader, *sync.WaitGroup) {
	if saveDir == "" {
		saveDir = os.Getenv("SAVE_PATH")
		if saveDir == "" {
			saveDir = filepath.Join(".", "download")
		}
	}

	cmd := exec.Command(`N_m3u8DL-RE`, `--ui-language`, `en-US`, `--auto-select`, "--save-dir", saveDir, "--tmp-dir", saveDir, param.Url)
	if param.Title != "" {
		cmd.Args = append(cmd.Args, `--save-name`, param.Title)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Use WaitGroup to ensure all output is processed
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Create multi-reader pipe
	pr, pw := io.Pipe()

	// Process stdout
	go func() {
		defer wg.Done()
		io.Copy(pw, stdout)
	}()

	// Process stderr
	go func() {
		defer wg.Done()
		io.Copy(pw, stderr)
	}()

	// Close pipe writer after both streams are processed
	go func() {
		wg.Wait()
		pw.Close()
	}()

	return cmd, pr, wg
}
