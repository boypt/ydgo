# YouDao Console Version

Simple wrapper for Youdao online translate (Chinese <-> English) service [API](http://fanyi.youdao.com/openapi?path=data-mode), as an alternative to the StarDict Console Version(sdcv).

This is a GO portation from the original [ydcv in python](https://github.com/felixonmars/ydcv)

This project is a practise during the study of Go (newbie work in 100 lines). If you have any idea improving this project, pull request is welcomed.

# Config

User need to apply for his own service API key from [有道智云](https://ai.youdao.com).
Keys writted to a config file `.ydgo`, put at `HOME` dir, or `%USERPROFILE%` for windows.

Example:
```
$ cat ~/.ydgo
YDAPPID="123456"
YDAPPSEC="abcd1234"
```

## Environment:
 * Any OS that GO supports (windows, linux, freebsd, darwin ...)
