package main

import (
	"fmt"
	"github.com/fatih/color"
)

type Logger struct {}

func (l Logger) Error(msg string) {
	fmt.Println(color.RedString(msg))
}

func (l Logger) Success(msg string) {
	fmt.Println(color.GreenString(msg))
}

func (l Logger) SuccessF(format string, args ...interface{}) {
	d := color.New(color.FgGreen)
	d.Printf(format, args...)
}
