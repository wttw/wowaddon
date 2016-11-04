// +build !windows

package output

import "fmt"

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Print(v ...interface{}) {
	fmt.Print(v...)
}
