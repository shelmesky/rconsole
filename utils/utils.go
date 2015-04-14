package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
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

// 计算字符换的MD5
func MD5(text string) string {
	hashMD5 := md5.New()
	io.WriteString(hashMD5, text)
	return fmt.Sprintf("%x", hashMD5.Sum(nil))
}

// 获取随机id字符串
func MakeRandomID() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()

	md5_nano := MD5(strconv.FormatInt(nano, 10))
	md5_rand := MD5(strconv.FormatInt(rndNum, 10))
	return MD5(md5_nano + md5_rand)
}
