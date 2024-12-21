package main

import (
	"fmt"
	"os"
)

func DebugMsg(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
}
