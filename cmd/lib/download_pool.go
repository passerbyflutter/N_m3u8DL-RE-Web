package lib

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"github.com/panjf2000/ants/v2"
)

type PoolTask struct {
	ID         string
	Url        string
	Title      string
	Progress   uint8
	Status     DownloadStatus
	CreateTime time.Time
	StartTime  time.Time
	FinishTime time.Time
	Cmd        *exec.Cmd
}

type DownloadPool struct {
	pool *ants.Pool
}

func NewDownloadPool(size int) *DownloadPool {
	defer ants.Release()

	pool, _ := ants.NewPool(size)

	return &DownloadPool{
		pool: pool,
	}
}

func (downloadPool *DownloadPool) AdjustSize(size int) {
	downloadPool.pool.Tune(size)
}

func (downloadPool *DownloadPool) AddDownloadTask(url string, title string) *PoolTask {
	task := PoolTask{
		ID:         shortuuid.New(),
		Url:        url,
		Title:      title,
		Status:     Pending,
		CreateTime: time.Now().UTC(),
	}

	go downloadPool.pool.Submit(task.downloadTaskHandler)

	return &task
}

func (task *PoolTask) downloadTaskHandler() {
	if task.Status == Deleted {
		return
	}

	cmd, ioReader := GenerateCmd(Param{
		Url:   task.Url,
		Title: task.Title,
	})
	task.Cmd = cmd
	task.Status = Downloading
	task.StartTime = time.Now().UTC()
	cmd.Start()

	scanner := bufio.NewScanner(ioReader)
	scanner.Split(bufio.ScanLines)
	saveNameRegex, _ := regexp.Compile(`Save Name: ([\w\W]+)`)
	progressRegex, _ := regexp.Compile(`(\d{0,3})%$`)

	for scanner.Scan() {
		msg := scanner.Text()

		if task.Title == "" {
			match := saveNameRegex.MatchString(msg)
			if match {
				task.Title = saveNameRegex.FindStringSubmatch(msg)[1]
			}
		}

		match := progressRegex.MatchString(msg)
		if match {
			progressText := progressRegex.FindStringSubmatch(msg)[1]
			ui64, _ := strconv.ParseUint(progressText, 10, 8)
			task.Progress = uint8(ui64)
		}
	}
	cmd.Wait()
	task.Status = Finished
	task.FinishTime = time.Now().UTC()
	task.Progress = uint8(100)
}
