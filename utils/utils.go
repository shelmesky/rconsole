package utils

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Println(args ...interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	funcname := function.Name()
	msg := fmt.Sprintf(generateFmtStr(len(args)), args...)
	log.Printf("[%s : %d] %s", funcname, line, msg)
}

func Printf(format string, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	funcname := function.Name()
	msg := fmt.Sprintf(format, args...)
	log.Printf("[%s : %d] %s", funcname, line, msg)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
