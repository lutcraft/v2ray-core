package base

import (
	"os"
	"path"
)

// CommandEnvHolder is a struct holds the environment info of commands
// CommandEnvHolder是存储着命令的环境信息的结构，所谓环境信息包含当前命令的二进制文件名和命令块长度
type CommandEnvHolder struct {
	// Excutable name of current binary
	Exec string
	// commands column width of current command
	CommandsWidth int
}

// CommandEnv holds the environment info of commands
var CommandEnv CommandEnvHolder

func init() {
	exec, err := os.Executable()
	if err != nil {
		return
	}
	CommandEnv.Exec = path.Base(exec)
}
