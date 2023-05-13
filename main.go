package main

import (
	"github.com/zhangyiming748/AVmerger"
	"github.com/zhangyiming748/goini"
	"github.com/zhangyiming748/pretty"
	"github.com/zhangyiming748/processAudio"
	"github.com/zhangyiming748/processImage"
	"github.com/zhangyiming748/processVideo"
	"github.com/zhangyiming748/resizeVideo"
	"github.com/zhangyiming748/rotateVideo"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"strings"
	"time"
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
func startOn(t string) {
	for true {
		now := time.Now().Local().Format("15")
		if t == now {
			return
		} else {
			logger.Warn("still alive", slog.Any("time", now), slog.String("target", t))
			time.Sleep(30 * time.Minute)
		}
	}
}
func main() {

	os.Setenv("QUIET", "True")
	if len(os.Args) > 1 {
		slog.Info("使用自定义配置文件", slog.String("配置文件路径", os.Args[1]))
		conf = goini.SetConfig(os.Args[1])
	} else {
		slog.Info("使用默认配置文件")
		conf = goini.SetConfig(configPath)
	}
	level, _ := conf.GetValue("log", "level")
	mission, _ := conf.GetValue("main", "mission")
	config := conf.ReadList()
	pretty.P(config)
	setLevel(level)
	err := os.Setenv("LEVEL", level)
	if err != nil {
		logger.Error("设置日志输出环境变量失败")
		return
	}
	var (
		root      string
		pattern   string
		threads   string
		direction string
	)
	staterOn, _ := conf.GetValue("StartAt", "time")
	startOn(staterOn)
	switch mission {
	case "video":
		pattern, _ = conf.GetValue("pattern", "video")
		root, _ = conf.GetValue("root", "video")
		threads, _ = conf.GetValue("thread", "threads")
		logger.Info("开始视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		processVideo.ProcessAllVideos(root, pattern, threads, false)
	case "audio":
		pattern, _ = conf.GetValue("pattern", "audio")
		pattern = strings.Join([]string{pattern, strings.ToUpper(pattern)}, ";")
		root, _ = conf.GetValue("root", "audio")
		logger.Info("开始音频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		processAudio.ProcessAllAudios(root, pattern)
	case "image":
		pattern, _ = conf.GetValue("pattern", "image")
		root, _ = conf.GetValue("root", "image")
		threads, _ = conf.GetValue("thread", "threads")
		logger.Info("开始图片处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		processImage.ProcessAllImages(root, pattern, threads)
	case "rotate":
		pattern, _ = conf.GetValue("pattern", "video")
		root, _ = conf.GetValue("root", "video")
		threads, _ = conf.GetValue("thread", "threads")
		direction, _ = conf.GetValue("rotate", "direction")
		logger.Info("开始旋转视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads), slog.String("方向", direction))
		rotateVideo.Rotate(root, pattern, direction, threads)
	case "resize":
		pattern, _ = conf.GetValue("pattern", "video")
		root, _ = conf.GetValue("root", "video")
		threads, _ = conf.GetValue("thread", "threads")
		logger.Info("开始缩小视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		resizeVideo.ResizeAllVideos(root, pattern, threads)
	case "avmerger":
		root, _ = conf.GetValue("bilibili", "root")
		logger.Info("开始合并哔哩哔哩进程", slog.String("根目录", root))
		AVmerger.AllIn(root)
	case "speedUp":
		root, _ = conf.GetValue("root", "speedUp")
		pattern, _ = conf.GetValue("pattern", "speedUp")
		processAudio.SpeedUpAudios(root, pattern, processAudio.AudioBook)
		logger.Info("开始有声小说加速处理", slog.String("根目录", root))
	case "gif":
		root, _ = conf.GetValue("root", "gif")
		pattern, _ = conf.GetValue("pattern", "gif")
		threads, _ = conf.GetValue("thread", "threads")
		processImage.ProcessAllImagesLikeGif(root, pattern, threads)
	}
}
