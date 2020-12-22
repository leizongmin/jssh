package jsshcmd

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/leizongmin/jssh/internal/pkginfo"
	"os"
)

func printAuthorInfo() {
	fmt.Printf("Welcome to %s %s\n", pkginfo.Name, pkginfo.LongVersion)
	fmt.Println("Author:  leizongmin@gmail.com")
	fmt.Println("Project: https://github.com/leizongmin/jssh")
	fmt.Println()
}

func printUsage(code int) {
	printAuthorInfo()
	fmt.Println("Example usage:")
	fmt.Printf("  %s script_file.js [arg1] [arg2] [...]                 Run script file\n", pkginfo.Name)
	fmt.Printf("  %s http://example.com/file.js [arg1] [arg2] [...]     Run script from URL\n", pkginfo.Name)
	fmt.Printf("  %s exec \"script\" [arg1] [arg2] [...]                  Run script from argument\n", pkginfo.Name)
	fmt.Printf("  %s eval \"script\" [arg1] [arg2] [...]                  Run script from argument and print the result\n", pkginfo.Name)
	fmt.Printf("  %s build script_file.js [target_file]                 Create self-contained binary file\n", pkginfo.Name)
	fmt.Printf("  %s repl                                               Start REPL\n", pkginfo.Name)
	fmt.Printf("  %s help                                               Show usage\n", pkginfo.Name)
	fmt.Printf("  %s version                                            Show version\n", pkginfo.Name)
	fmt.Println()
	os.Exit(code)
}

func printExitMessage(message string, code int, usage bool) {
	fmt.Println(color.FgRed.Render(message))
	if usage {
		fmt.Println()
		printUsage(code)
	} else {
		os.Exit(code)
	}
}
