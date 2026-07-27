# sandbox-daemon

[![Go Reference](https://pkg.go.dev/badge/github.com/gongt/sandbox-daemon.svg)](https://pkg.go.dev/github.com/gongt/sandbox-daemon)

一个全功能的容器内init。

*主要用于e2b这种轻量虚拟机，docker中有更好的选择*

## 主要功能

### PID 1常见功能

1. 启动指定主程序
1. 转发 SIGTERM、SIGINT 信号
1. 处理 SIGCHLD 信号，回收僵尸进程
1. 转发退出码给父进程

### stdout、stderr 转发

1. 同时转发输出和错误，同时支持二者混合，三条通道互相独立。
1. 支持多种输出，且可以复制到多个目标:
	1. 文件（管道）
	2. 通过http(s)、websocket、tcp、udp发送JSON

### 额外初始化

1. 创建命名空间、挂载文件系统
1. 主程序执行前运行指定的初始化命令或脚本
1. 等待指定的文件、端口就绪

### 重启与热更新

1. 重启主程序
1. 动态更新主程序
1. 动态更新本init程序自身

# 主程序

与完整的init不同，sandbox-daemon只关注单一程序的启动，此程序可以是其他init程序。但更常见的情况是启动一个web服务、claude code等。

主程序退出后，sandbox-daemon会退出，任何情况下都不会尝试重启之类的操作。除非是通过`sb stop`命令，主程序停止后不退出。

# 使用方法

1. 启动守护进程

通常在构建template的最终运行阶段执行

会准备好运行环境、日志系统

如果有-program，则还会立即启动主程序

```bash
sb init [-dir /var/run/sandbox-daemon] [-program /app/my-program.yaml]
```

2. 启动、停止

读取配置文件中的主程序配置，让init进程启动它，可选等待启动完成。

如果已有主程序，则报错；除非设置了-replace参数。

* -log: 桥接输出到当前终端，方便调试
* -wait: 等待主程序启动完成

```bash
sb start [-dir /var/run/sandbox-daemon] [-replace] [-log] [-wait]
```

停止主程序，如果未运行，直接返回成功

```bash
sb stop [-dir /var/run/sandbox-daemon]
```

判断主程序是否存在，如果设置-signal，则向主程序发送指定信号

```bash
sb kill [-dir /var/run/sandbox-daemon] [-signal SIUHUP]
```

3. 桥接stdio

启动一个交互进程，连接到主程序的stdin、stdout、stderr，用于调试

按Ctrl+C退出，不会终止主程序

```bash
sb attach [-dir /var/run/sandbox-daemon] [-ro]
```

4. 在主程序空间运行其他程序

包括mount namespace、cwd、环境变量等信息

不要求主程序必须运行
* 如果主程序**已**在运行，则主程序退出**前**此程序会收到SIGTERM
* 如果主程序**未**在运行，则主程序退出**后**此程序会收到SIGTERM
停止动作超时会被KILL

* -pipe: 隐含-wait，桥接当前终端的stdin、stdout、stderr。不设置（默认）则: 无stdin; stdout、stderr输出到和主程序相同的地方
* -wait: 等待程序退出，返回退出码
* -nons: 使用根namespace，而不是主程序的
* -noenv: 不应用主程序的环境变量设置

```bash
sb exec [-dir /var/run/sandbox-daemon] [-wait] [-pipe] [-nons] [-noenv]
```

# 主程序配置文件

[./example.yaml](./example.yaml)是一个示例配置文件，包含了所有可用的配置项。
