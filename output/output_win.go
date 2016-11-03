// +build windows

package output

import (
	"fmt"
	"github.com/mattn/go-colorable"
)

func Printf(format string, v ...interface{}) {
	if _, err := colorOut.Write([]byte(fmt.Sprintf(format, v...))); err != nil {
		panic(err)
	}
}

func Print(v ...interface{}) {
	if _, err := colorOut.Write([]byte(fmt.Sprint(v...))); err != nil {
		panic(err)
	}
}

var colorOut = colorable.NewColorableStdout()
