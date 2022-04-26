package main

import (
	"github.com/v2fly/v2ray-core/v5/main/commands"
	"github.com/v2fly/v2ray-core/v5/main/commands/base"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
)

func main() {
	//变更RootCommand命令的说明
	base.RootCommand.Long = "A unified platform for anti-censorship."
	base.RegisterCommand(commands.CmdRun)     //注册v2ray run命令
	base.RegisterCommand(commands.CmdVersion) //注册v2ray version命令
	base.RegisterCommand(commands.CmdTest)    //注册v2ray test命令
	base.SortLessFunc = runIsTheFirst         //用runIsTheFirst覆盖SortLessFunc的默认实现，此方法用来给定各种命令的排序方式
	base.SortCommands()
	base.Execute()
}

func runIsTheFirst(i, j *base.Command) bool {
	left := i.Name()
	right := j.Name()
	if left == "run" {
		return true
	}
	if right == "run" {
		return false
	}
	return left < right
}
