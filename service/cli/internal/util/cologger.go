package colog

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

func Fatal(format string, args ...interface{}) {
	log.Fatalln(color.RedString(format, args...))
}

func Error(format string, args ...interface{}) {
	fmt.Println(color.RedString(format, args...))
}

func Warn(format string, args ...interface{}) {
	fmt.Println(color.YellowString(format, args...))
}

func Success(format string, args ...interface{}) {
	fmt.Println(color.GreenString(format, args...))
}

func Info(format string, args ...interface{}) {
	fmt.Println(color.WhiteString(format, args...))
}
