# TestAll

# 配置文件说明

```ini
[main]
mission = video // 运行程序的类别
# mission = audio
# mission = image
# mission = rotate
# mission = resize
[log]
level = Debug // 输出日志的等级
# level=Info
# level=Warn
# level=Error
[pattern]
image = jpeg;jpg;JPG;png;webp
audio = mp3;m4a;flac;MP3;wma;wav
video = webm;mkv;m4v;MP4;mp4;mov;MOV;avi;wmv;ts;TS;rmvb
[root]
image = /Users/zen/Downloads // 资源所在目录
audio = /Volumes/swap/Back
video = /Volumes/swap/Back
[thread]
threads = 2 // 运行所需线程数
[rotate]
direction = ToRight // 视频旋转方向
# direction = ToLeft
```