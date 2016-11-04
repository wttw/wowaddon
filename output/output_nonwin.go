// +build !windows

package output

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Print(v ...interface{}) {
	fmt.Print(v...)
}
