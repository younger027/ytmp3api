package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"ytmp3api/covert"
)

func RunGinServer() {
	router := gin.New()
	pprof.Register(router)

	//not need auth
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/convert", Convert)

	if err := router.Run(":8888"); err != nil {
		panic(err)
	}
}

func Convert(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ytURL := ctx.Query("url")
	if ytURL == "" {
		ctx.JSON(http.StatusBadRequest, "fail:empty youtube url")
		return
	}

	videoID, err := covert.ExtractVideoID(ytURL)
	if videoID == "" {
		ctx.JSON(http.StatusBadRequest, "fail:get youtube id from url")
		return
	}

	quality := ctx.Query("quality")
	if quality == "" {
		quality = "128k"
	}

	filePath := "/Users/rockey-lyy/ad-tencent/ytmp3api/musicsource/" + strings.Split(ytURL, "v=")[1] + "-" + quality
	//yt-dlp --extract-audio --audio-format mp3 --audio-quality 320k https://www.youtube.com/watch\?v\=YudHcBIxlYw -o /Users/rockey-lyy/ad-tencent/320kid.mp3
	cmdArray := []string{"--extract-audio", "--audio-format", "mp3", "--audio-quality", quality, ytURL, "-o", filePath}
	cmd := exec.CommandContext(c, "yt-dlp", cmdArray...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("CombinedOutput:", string(out))
		ctx.JSON(http.StatusInternalServerError, "success convert")
		return
	}

	fmt.Println("out:", string(out))
	ctx.JSON(http.StatusOK, "success convert")
	return
}
