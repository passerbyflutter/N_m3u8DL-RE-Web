package main

import (
	"N_m3u8DL-RE-API/cmd/lib"
	"N_m3u8DL-RE-API/web"
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/devfeel/mapper"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

type taskRequest struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type taskResponse struct {
	ID         string             `json:"id"`
	Url        string             `json:"url"`
	Title      string             `json:"title"`
	Progress   uint8              `json:"progress"`
	Status     lib.DownloadStatus `json:"status"`
	CreateTime time.Time          `json:"createTime"`
	StartTime  time.Time          `json:"startTime"`
	FinishTime time.Time          `json:"finishTime"`
}

var (
	version   string
	buildTime string
	isRelease string = "false"
)

var downloadTasks = make(map[string]*lib.PoolTask)
var downloadPool *lib.DownloadPool

func main() {
	if isRelease == "true" {
		fmt.Printf("Version %s, build at %s.\n\n", version, buildTime)
		gin.SetMode(gin.ReleaseMode)
	}

	downloadPool = lib.NewDownloadPool(getDownloadPoolSize(3))

	router := gin.Default()

	setSpaStaticResource(router)
	setApiHandler(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefulShutdown(srv)
}

func getDownloadPoolSize(defaultValue int) int {
	downloadPoolSizeText := os.Getenv("DOWNLOAD_POOL_SIZE")

	if downloadPoolSizeText != "" {
		if downloadPoolSizeInt64, err := strconv.ParseUint(downloadPoolSizeText, 10, 8); err == nil {
			return int(downloadPoolSizeInt64)
		}
	}

	return defaultValue
}

func setSpaStaticResource(router *gin.Engine) {
	ui, _ := fs.Sub(web.Static, "dist")
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			filePath := c.Request.URL.Path[1:]
			file, err := ui.Open(filePath)
			if err != nil {
				filePath = "index.html"
				file, _ = ui.Open(filePath)
			}
			stat, _ := file.Stat()
			fileBytes := make([]byte, stat.Size())
			bufio.NewReader(file).Read(fileBytes)
			c.Data(http.StatusOK, mime.TypeByExtension(filepath.Ext(filePath)), fileBytes)
		}
	})
}

func setApiHandler(router *gin.Engine) {
	router.GET("/api/tasks", listTasks)
	router.POST("/api/tasks", AddTasks)
	router.DELETE("/api/tasks/:id", DeleteTasks)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}

func listTasks(c *gin.Context) {
	taskResponseList := &[]*taskResponse{}

	mapper.SetEnabledJsonTag(false)
	mapper.MapperSlice(maps.Values(downloadTasks), taskResponseList)

	sort.Slice(*taskResponseList, func(i, j int) bool {
		return (*taskResponseList)[i].CreateTime.Before((*taskResponseList)[j].CreateTime)
	})

	c.JSON(http.StatusOK, taskResponseList)
}

func AddTasks(c *gin.Context) {
	var taskRequest taskRequest

	if err := c.BindJSON(&taskRequest); err != nil {
		return
	}

	poolTask := downloadPool.AddDownloadTask(taskRequest.Url, taskRequest.Title)
	downloadTasks[poolTask.ID] = poolTask

	taskResponse := &taskResponse{}
	mapper.SetEnabledJsonTag(false)
	mapper.Mapper(poolTask, taskResponse)

	c.JSON(http.StatusOK, taskResponse)
}

func DeleteTasks(c *gin.Context) {
	id := c.Param("id")
	if poolTask, ok := downloadTasks[id]; ok {
		if poolTask.Cmd != nil {
			poolTask.Cmd.Process.Signal(os.Kill)
		}
		poolTask.Status = lib.Deleted
		delete(downloadTasks, id)
		c.JSON(http.StatusOK, taskResponse{ID: id})
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": fmt.Sprintf("ID %s not found!", id)})
	}
}
