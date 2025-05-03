# DBD_QTE_go

## 项目简介

`DBD_QTE_go` 是一个使用 Go 语言开发的项目，使用go语言模拟了黎明杀机游戏的QTE检定，调用了 [DG-Lab-Coyote-Game-Hub](https://github.com/hyperzlib/DG-Lab-Coyote-Game-Hub) 提供的接口。它的功能是通过获取游戏的 `clientId`，然后运行游戏相关的逻辑，并通过构建一个 Go 项目来启动程序。

## 安装依赖

首先，确保你已经安装了 Go 语言环境。如果没有安装 Go，请访问 [Go 官方网站](https://golang.org/dl/)进行安装。

### 1. 克隆项目

```bash
git clone https://github.com/H1W0XXX/DBD_QTE_go.git
cd DBD_QTE_go
```
2. 安装第三方库
使用以下命令安装项目所依赖的第三方库：

```bash
go get github.com/gen2brain/raylib-go/raylib
```
此命令会自动安装项目所需的库，包括 raylib-go，这是一个用于创建图形和处理图像的库。

配置游戏常量
在 main.go 文件中，你需要修改一些游戏常量来适应你的需求。打开 main.go 文件并找到相关配置区域，修改常量值来适应你的环境。例如：

```go
// 在 main.go 中找到游戏常量部分，自己修改游戏参数
const (
	PointerSpeed    = 120
	NormalQTEAngle  = 30
	PerfectQTEAngle = 10

	PerfectQTEProgressBonus  = 3
	FailureProgressRetreat   = 1
	FailureProgressPauseTime = 300 * time.Millisecond
	QTETipSoundDuration      = 500 * time.Millisecond
	frameRate                = 120
	rotationSpeedPerFrame    = PointerSpeed / float32(frameRate)
	LongestWaitingTime       = time.Second * 15
	ShortestWaitingTime      = time.Second * 5
)

```
获取 clientId
接着启动 Coyote game Hub 后，在浏览器中打开 http://127.0.0.1:8920/#/。

在网页上获取 clientId。

编译和启动
1. 编译项目
使用以下命令来编译 Go 项目：

```bash
go build -o DBD_QTE_go.exe
```
2. 启动程序
在终端中运行编译后的程序：

```bash
./DBD_QTE_go.exe
```
启动后，程序会提示你输入刚才获取的 clientId。请输入并按回车。
