// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package base defines shared basic pieces of the commands,
// in particular logging and the Command structure.
package base

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

// A Command is an implementation of a v2ray command
// like v2ray run or v2ray version.
// Command结构体是一个v2ray命令的实现
// 比如v2ray run，或者v2ray version。
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The words between the first word (the "executable name") and the first flag or argument in the line are taken to be the command name.
	//
	// UsageLine supports go template syntax. It's recommended to use "{{.Exec}}" instead of hardcoding name
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	//
	// Note: Short does not support go template syntax.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	//
	// Long supports go template syntax. It's recommended to use "{{.Exec}}", "{{.LongName}}" instead of hardcoding strings
	// Long变量为 在go help <this-command>命令中的输出字符串
	// Long变量支持go模版语法，推荐使用{{.Exec}}或者{{.LongName}}的写法，而不要使用硬编码
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool

	// Commands lists the available commands and help topics.
	// The order here is the order in which they are printed by 'go help'.
	// Note that subcommands are in general best avoided.
	Commands []*Command
}

// LongName returns the command's long name: all the words in the usage line between first word (e.g. "v2ray") and a flag or argument,
func (c *Command) LongName() string {
	name := c.UsageLine
	if i := strings.Index(name, " ["); i >= 0 {
		name = strings.TrimSpace(name[:i])
	}
	if i := strings.Index(name, " "); i >= 0 {
		name = name[i+1:]
	} else {
		name = ""
	}
	return strings.TrimSpace(name)
}

// Name returns the command's short name: the last word in the usage line before a flag or argument.
func (c *Command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return strings.TrimSpace(name)
}

// Usage prints usage of the Command
func (c *Command) Usage() {
	buildCommandText(c)
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "Run '%s help %s' for details.\n", CommandEnv.Exec, c.LongName())
	SetExitStatus(2)
	Exit()
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// Exit exits with code set with SetExitStatus()
func Exit() {
	os.Exit(exitStatus)
}

// Fatalf logs error and exit with code 1
func Fatalf(format string, args ...interface{}) {
	Errorf(format, args...)
	Exit()
}

// Errorf logs error and set exit status to 1, but not exit
func Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
	SetExitStatus(1)
}

// ExitIfErrors exits if current status is not zero
func ExitIfErrors() {
	if exitStatus != 0 {
		Exit()
	}
}

var (
	exitStatus = 0
	exitMu     sync.Mutex
)

// SetExitStatus set exit status code
func SetExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

// GetExitStatus get exit status code
func GetExitStatus() int {
	return exitStatus
}
