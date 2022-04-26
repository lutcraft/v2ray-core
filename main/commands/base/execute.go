package base

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// copied from "github.com/golang/go/main.go"

// Execute excute the commands
func Execute() {
	//解析命令参数
	flag.Parse()
	args := flag.Args()
	//如果参数小于1（为0）
	if len(args) < 1 {
		PrintUsage(os.Stderr, RootCommand)
		return
	}
	//获取命令参数
	cmdName := args[0] // for error messages
	if args[0] == "help" {
		//如果第一个参数为help就打印 对应命令的help:比如v2ray help run
		Help(os.Stdout, args[1:])
		return
	}

BigCmdLoop:
	for bigCmd := RootCommand; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				// test sub commands
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					PrintUsage(os.Stderr, bigCmd)
					SetExitStatus(2)
					Exit()
				}
				if args[0] == "help" {
					// Accept 'go mod help' and 'go mod help foo' for 'go help mod' and 'go help mod foo'.
					Help(os.Stdout, append(strings.Split(cmdName, " "), args[1:]...))
					return
				}
				cmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}

			buildCommandText(cmd)
			cmd.Run(cmd, args)
			Exit()
			return
		}
		helpArg := ""
		if i := strings.LastIndex(cmdName, " "); i >= 0 {
			helpArg = " " + cmdName[:i]
		}
		fmt.Fprintf(os.Stderr, "%s %s: unknown command\nRun '%s help%s' for usage.\n", CommandEnv.Exec, cmdName, CommandEnv.Exec, helpArg)
		SetExitStatus(2)
		Exit()
	}
}

// SortCommands sorts the first level sub commands
func SortCommands() {
	//用给定的方法排序RootCommand.Commands数组，方法由 SortLessFunc 给定
	sort.Slice(RootCommand.Commands, func(i, j int) bool {
		return SortLessFunc(RootCommand.Commands[i], RootCommand.Commands[j])
	})
}

// SortLessFunc used for sort commands list, can be override from outside
var SortLessFunc = func(i, j *Command) bool {
	return i.Name() < j.Name()
}
