package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os/exec"
	"time"
	"ytmp3api/cache"
	convert "ytmp3api/covert"
)

func RunGinServer() {
	router := gin.New()
	pprof.Register(router)

	//not need auth
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/convert", Convert)

	if err := router.Run(":8989"); err != nil {
		panic(err)
	}
}

type ConvertResponse struct {
	Path     string `json:"path,omitempty"`
	ErrorMsg string `json:"error_msg,omitempty"`
}

func Convert(ctx *gin.Context) {
	resp := new(ConvertResponse)

	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ytURL := ctx.Query("url")
	videoID, _ := convert.ExtractVideoID(ytURL)

	if ytURL == "" || videoID == "" {
		resp.ErrorMsg = convert.ApiConvertInvalidParams
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	//64kbps-128kbps-256kbps-320kbps
	quality := ctx.Query("quality")
	if quality == "" {
		quality = "128k"
	}

	//directory := "/Users/rockey-lyy/ad-tencent/ytmp3api/musicsource/"
	directory := "/root/ytmp3api/musicsource/download/"
	fileName := fmt.Sprintf("%s-%s", videoID, quality)
	downLoadName := "http://154.82.111.99/" + fileName + ".mp3"
	filePath := directory + fileName

	//缓存检查，确保缓存中的文件路径和目录下的一致
	_, ok := cache.GetCacheManger().Get(fileName)
	if ok {
		resp.Path = downLoadName
		ctx.JSON(http.StatusOK, resp)
		return
	}

	//yt-dlp --extract-audio --audio-format mp3 --audio-quality 320k https://www.youtube.com/watch\?v\=YudHcBIxlYw -o /Users/rockey-lyy/ad-tencent/320kid.mp3
	cmdArray := []string{"--extract-audio", "--audio-format", "mp3", "--audio-quality", quality, ytURL, "-o", filePath}
	cmd := exec.CommandContext(c, "yt-dlp", cmdArray...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("CombinedOutput:", string(out))
		resp.ErrorMsg = convert.ApiConvertCMDExecFail
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	fmt.Println("out:", string(out))

	//TODO 检查返回值的正确性
	/*out:
	[youtube] Extracting URL: https://www.youtube.com/watch?v=YudHcBIxlYw
	[youtube] YudHcBIxlYw: Downloading webpage
	[youtube] YudHcBIxlYw: Downloading android player API JSON
	[info] YudHcBIxlYw: Downloading 1 format(s): 251
	[dashsegments] Total fragments: 1
	[download] Destination: yang
	[download] 100% of    3.10MiB in 00:00:00 at 62.39MiB/s
	[ExtractAudio] Destination: yang.mp3
	Deleting original file yang (pass -k to keep)
	*/

	//应该判断文件夹下是否有对应名称的文件生成
	filePath += ".mp3"
	isExist, _ := convert.PathExists(filePath)
	if !isExist {
		resp.ErrorMsg = convert.ApiConvertDownloadFail
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	//TODO 增加host，返回完整路径
	resp.Path = downLoadName
	ctx.JSON(http.StatusOK, resp)

	//fileName入缓存
	cache.GetCacheManger().Add(fileName, "")

	return
}
