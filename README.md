# TestAll

# 配置文件说明

```ini
[main]
mission = i&v
# mission = video
# mission = audio
# mission = image
# mission = rotate
# mission = resize
# mission = avmerger
# mission = speedUp
# mission = gif
[log]
# level = Debug
level = Info
# level = Warn
# level = Error
[pattern]
image = jpeg;jpg;JPG;png;webp
audio = mp3;m4a;flac;wma;wav
speedUp = mp3;m4a;flac;MP3;wma;wav;aac
video = webm;mkv;m4v;MP4;mp4;mov;MOV;avi;wmv;ts;TS;rmvb;wma;WMA;avi;AVI
gif = gif;webm
[root]
image = /Users/zen/Downloads/4x3
audio = /Users/zen/Downloads/original
video = /Users/zen/Downloads/4x3
speedUp = /Users/zen/Downloads/YvDSJYMGBmjdKkPG/elegram
gif = /Users/zen/Downloads/Fugtrup collection/archive/webp
[thread]
threads = 10
[rotate]
direction = ToRight
# direction = ToLeft
[bilibili]
root = /Users/zen/Downloads
[StartAt]
time = 10
[alert]
quiet = yes
email = yes
[email]
username = 18904892728@163.com
password = SMTP授权码
from = 18904892728@163.com
tos = 578779391@qq.com;zhangyiming748@gmail.com
```