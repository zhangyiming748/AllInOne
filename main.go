package main

import (
	"github.com/zhangyiming748/AVmerger"
	"github.com/zhangyiming748/goini"
	"github.com/zhangyiming748/processAudio"
	"github.com/zhangyiming748/processImage"
	"github.com/zhangyiming748/processVideo"
	"github.com/zhangyiming748/resizeVideo"
	"github.com/zhangyiming748/rotateVideo"
	"golang.org/x/exp/slog"
	"io"
	"os"
)

const (
	configPath = "./settings.ini"
)

var (
	conf *goini.Config
)
var (
	logger *slog.Logger
)

// todo 全部修改为显式传递日志等级
func setLevel(level string) {
	var opt slog.HandlerOptions
	switch level {
	case "Debug":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}
	case "Info":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelInfo, // slog 默认日志级别是 info
		}
	case "Warn":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelWarn, // slog 默认日志级别是 info
		}
	case "Err":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelError, // slog 默认日志级别是 info
		}
	default:
		slog.Warn("需要正确设置环境变量 Debug,Info,Warn or Err")
		slog.Info("默认使用Debug等级")
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}
	}
	file := "AllInOne.log"
	logf, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	logger = slog.New(opt.NewJSONHandler(io.MultiWriter(logf, os.Stdout)))
}
func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("程序运行过程中有错误产生", slog.Any("错误内容", err))
		}
	}()
	conf = goini.SetConfig(configPath)
	var (
		mission   string
		root      string
		pattern   string
		threads   string
		direction string
	)
	level, _ := conf.GetValue("log", "level")
	setLevel(level)
	if len(os.Args) > 1 {
		slog.Warn("没有指定功能")
		return
	} else {
		mission = os.Args[1]
		switch mission {
		case "video":
			pattern, _ = conf.GetValue("pattern", "video")
			root, _ = conf.GetValue("root", "video")
			threads, _ = conf.GetValue("thread", "")
			logger.Info("开始视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
			processVideo.ProcessAllVideos(root, pattern, threads, false)
		case "audio":
			pattern, _ = conf.GetValue("pattern", "audio")
			root, _ = conf.GetValue("root", "audio")
			threads, _ = conf.GetValue("thread", "threads")
			logger.Info("开始音频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
			processAudio.ProcessAllAudios(root, pattern)
		case "image":
			pattern, _ = conf.GetValue("pattern", "image")
			root, _ = conf.GetValue("root", "image")
			threads, _ = conf.GetValue("thread", "threads")
			logger.Info("开始图片处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
			processImage.ProcessAllImages(root, pattern, threads)
		case "avmerger":
			AVmerger.AllIn(root)
		case "rotate":
			pattern, _ = conf.GetValue("pattern", "rotate")
			root, _ = conf.GetValue("root", "rotate")
			threads, _ = conf.GetValue("thread", "threads")
			direction, _ = conf.GetValue("rotate", "direction")
			logger.Info("开始旋转视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads), slog.String("方向", direction))
			rotateVideo.Rotate(root, pattern, direction, threads)
		case "resize":
			pattern, _ = conf.GetValue("pattern", "resize")
			root, _ = conf.GetValue("root", "resize")
			threads, _ = conf.GetValue("thread", "threads")
			logger.Info("开始缩小视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
			resizeVideo.ResizeAllVideos(root, pattern, threads)
		}
	}
}
