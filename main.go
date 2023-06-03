package main

import (
	"fmt"
	"github.com/zhangyiming748/AVmerger"
	"github.com/zhangyiming748/goini"
	"github.com/zhangyiming748/processAudio"
	"github.com/zhangyiming748/processImage"
	"github.com/zhangyiming748/processVideo"
	"github.com/zhangyiming748/sendEmailAlert"
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
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(logf, os.Stdout), &opt))
	slog.SetDefault(logger)
}
func startOn(t string) {
	for true {
		now := time.Now().Local().Format("15")
		if t == now {
			return
		} else {
			slog.Warn("still alive", slog.Any("time", now), slog.String("target", t))
			time.Sleep(30 * time.Minute)
		}
	}
}
func main() {
	start := time.Now()
	if len(os.Args) > 1 {
		slog.Info("使用自定义配置文件", slog.String("配置文件路径", os.Args[1]))
		conf = goini.SetConfig(os.Args[1])
	} else {
		slog.Info("使用默认配置文件")
		conf = goini.SetConfig(configPath)
	}
	level, _ := conf.GetValue("log", "level")
	mission, _ := conf.GetValue("main", "mission")
	//config := conf.ReadList()
	//pretty.P(config)

	setLevel(level)
	err := os.Setenv("LEVEL", level)
	if err != nil {
		slog.Error("设置日志输出环境变量失败")
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
	if quiet, _ := conf.GetValue("alert", "quiet"); quiet == "yes" {
		os.Setenv("QUIET", "True")
		slog.Info("静音模式")
	}
	switch mission {
	case "i&v":
		{
			pattern, _ = conf.GetValue("pattern", "video")
			pattern = strings.Join([]string{pattern, strings.ToUpper(pattern)}, ";")
			root, _ = conf.GetValue("root", "video")
			threads, _ = conf.GetValue("thread", "threads")
			slog.Info("开始视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
			processVideo.ConvAllVideos2H265(root, pattern, threads)
		}

		{
			pattern, _ = conf.GetValue("pattern", "image")
			root, _ = conf.GetValue("root", "image")
			threads, _ = conf.GetValue("thread", "threads")
			pattern = strings.Join([]string{pattern, strings.ToUpper(pattern)}, ";")
			slog.Info("开始图片处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
			processImage.ProcessAllImages(root, pattern, threads)
		}
	case "video":
		pattern, _ = conf.GetValue("pattern", "video")
		root, _ = conf.GetValue("root", "video")
		threads, _ = conf.GetValue("thread", "threads")
		slog.Info("开始视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		processVideo.ConvVideos2H265(root, pattern, threads)
	case "audio":
		pattern, _ = conf.GetValue("pattern", "audio")
		pattern = strings.Join([]string{pattern, strings.ToUpper(pattern)}, ";")
		root, _ = conf.GetValue("root", "audio")
		slog.Info("开始音频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		processAudio.ConvAllAudios(root, pattern)
	case "image":
		pattern, _ = conf.GetValue("pattern", "image")
		root, _ = conf.GetValue("root", "image")
		threads, _ = conf.GetValue("thread", "threads")
		slog.Info("开始图片处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		processImage.ProcessAllImages(root, pattern, threads)
	case "rotate":
		pattern, _ = conf.GetValue("pattern", "video")
		root, _ = conf.GetValue("root", "video")
		threads, _ = conf.GetValue("thread", "threads")
		direction, _ = conf.GetValue("rotate", "direction")
		slog.Info("开始旋转视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads), slog.String("方向", direction))
		processVideo.Rotate(root, pattern, direction, threads)
	case "resize":
		pattern, _ = conf.GetValue("pattern", "video")
		root, _ = conf.GetValue("root", "video")
		threads, _ = conf.GetValue("thread", "threads")
		slog.Info("开始缩小视频处理进程", slog.String("根目录", root), slog.String("pattern", pattern), slog.String("进程数", threads))
		processVideo.ResizeAllVideos(root, pattern, threads)
	case "avmerger":
		root, _ = conf.GetValue("bilibili", "root")
		slog.Info("开始合并哔哩哔哩进程", slog.String("根目录", root))
		AVmerger.AllInH265(root)
	case "speedUp":
		root, _ = conf.GetValue("root", "speedUp")
		pattern, _ = conf.GetValue("pattern", "speedUp")
		processAudio.SpeedUpAudios(root, pattern, processAudio.AudioBook)
		slog.Info("开始有声小说加速处理", slog.String("根目录", root))
	case "gif":
		root, _ = conf.GetValue("root", "gif")
		pattern, _ = conf.GetValue("pattern", "gif")
		threads, _ = conf.GetValue("thread", "threads")
		processImage.ProcessAllImagesLikeGif(root, pattern, threads)
	}
	end := time.Now()
	if email, _ := conf.GetValue("alert", "email"); email == "yes" {
		slog.Info("发送任务完成邮件")
		sendEmail(start, end)
	}
}
func sendEmail(start, end time.Time) {
	i := new(sendEmailAlert.Info)
	if username, err := conf.GetValue("email", "username"); err == nil {
		i.SetUsername(username)
	}
	if password, err := conf.GetValue("email", "password"); err == nil {
		i.SetPassword(password)
	}
	if tos, err := conf.GetValue("email", "tos"); err == nil {
		i.SetTo(strings.Split(tos, ";"))
	}
	if from, err := conf.GetValue("email", "from"); err == nil {
		i.SetFrom(from)
	}
	i.SetHost(sendEmailAlert.NetEase.SMTP)
	i.SetPort(sendEmailAlert.NetEase.SMTPProt)
	subject, _ := conf.GetValue("main", "mission")
	i.SetSubject(strings.Join([]string{"AllInOne", subject, "任务完成"}, ""))
	text := strings.Join([]string{start.Format("任务开始时间 2006年01月02日 15:04:05"), end.Format("任务结束时间 2006年01月02日 15:04:05"), fmt.Sprintf("任务用时%.3f分\n", end.Sub(start).Minutes())}, "<br>")
	i.SetText(text)
	sendEmailAlert.Send(i)
}
