package base

// RootCommand is the root command of all commands
// RootCommand变量包含一个存储Command类型的数组，其中存储了所有注册了的命令
var RootCommand *Command

func init() {
	RootCommand = &Command{
		UsageLine: CommandEnv.Exec,
		Long:      "The root command",
	}
}

// RegisterCommand register a command to RootCommand
func RegisterCommand(cmd *Command) {
	RootCommand.Commands = append(RootCommand.Commands, cmd)
}
