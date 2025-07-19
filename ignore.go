package main

import (
	"fmt"
	"os"
)

// ignoreFprint ignores the integer and runs ignore(e) on the error
var ignoreFprint = func(_ int, e error) {
	ignore(e)
}

// ignore will only fmt.Fprintln to os.Stderr if e (error) is not nil
var ignore = func(e error) {
	if e != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: "+e.Error())
	}
}
