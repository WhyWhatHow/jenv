package main

import "github.com/whywhathow/jenv/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	// 执行命令行工具的主入口点
	cmd.Execute()
}
