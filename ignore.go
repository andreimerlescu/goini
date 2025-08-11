package main

import (
	"fmt"
	"os"
)

// ignore will only fmt.Fprintln to os.Stderr if e (error) is not nil
var ignore = func(e error) {
	if e != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: "+e.Error())
	}
}
