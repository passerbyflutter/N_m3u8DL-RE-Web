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

	// Initialize command and output processor
	cmd, reader, wg := GenerateCmd(Param{
		Url:   task.Url,
		Title: task.Title,
	})
	task.Cmd = cmd
	task.Status = Downloading
	task.StartTime = time.Now().UTC()

	// Start command
	if err := cmd.Start(); err != nil {
		task.Status = Error
		task.FinishTime = time.Now().UTC()
		return
	}

	// Configure scanner for output processing
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	saveNameRegex, _ := regexp.Compile(`Save Name: ([\w\W]+)`)
	progressRegex, _ := regexp.Compile(`(\d{0,3})%$`)
	errorRegex, _ := regexp.Compile(`(?i)(error|fail|exception)`)

	// Track error status
	hasError := false

	// Process output in new goroutine
	outputDone := make(chan bool)
	go func() {
		defer close(outputDone)
		for scanner.Scan() {
			msg := scanner.Text()

			// Check for error messages
			if errorRegex.MatchString(msg) {
				hasError = true
			}

			if task.Title == "" {
				if match := saveNameRegex.FindStringSubmatch(msg); match != nil {
					task.Title = match[1]
				}
			}

			if match := progressRegex.FindStringSubmatch(msg); match != nil {
				if progress, err := strconv.ParseUint(match[1], 10, 8); err == nil {
					task.Progress = uint8(progress)
				}
			}
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			hasError = true
		}
	}()

	// Wait for command completion
	cmdErr := cmd.Wait()

	// Wait for output processing to complete
	<-outputDone
	wg.Wait()

	// Set final status based on execution results
	task.FinishTime = time.Now().UTC()
	if cmdErr != nil || hasError {
		task.Status = Error
	} else {
		task.Status = Finished
		task.Progress = 100
	}
}
