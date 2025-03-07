package lib

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Param struct {
	Url, Title string
}

var saveDir string

func GenerateCmd(param Param) (*exec.Cmd, io.Reader) {
	if saveDir == "" {
		saveDir = os.Getenv("SAVE_PATH")
		if saveDir == "" {
			saveDir = filepath.Join(".", "download")
		}
	}

	cmd := exec.Command(`F:\Downloads\[tools]\N_m3u8DL-RE\N_m3u8DL-RE.exe`, `--ui-language`, `en-US`, `--auto-select`, "--save-dir", saveDir, "--tmp-dir", saveDir, param.Url)
	if param.Title != "" {
		cmd.Args = append(cmd.Args, `--save-name`, param.Title)
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	return cmd, io.MultiReader(stdout, stderr)
}
