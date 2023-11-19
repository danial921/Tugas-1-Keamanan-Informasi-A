package utils

import (
	"fmt"
	"os"
	"runtime"
)

// CheckError checks if an error is not nil and exits the program if it is.
func CheckError(err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("Error at %s:%d - %s\n", file, line, err)
		} else {
			fmt.Printf("Error - %s\n", err)
		}
		os.Exit(1)
	}
}
