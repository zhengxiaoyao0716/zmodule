# zmodule
## 正的模块模板

***
## Launch line
```
[START]
zmodule.Main() => event.Init() => SYS.init
  |
program.Start() => SYS.start
  |
SYS.stop => go run()

[STOP]
program.STOP()
```
