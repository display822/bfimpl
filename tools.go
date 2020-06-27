// +build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func changeFileModuleName(path string) {
	fmt.Println("scan", path, "...")
	fbuf, err := ioutil.ReadFile(path)
	if err != nil {
		checkError(err)
	}
	rd := bufio.NewReader(bytes.NewBuffer(fbuf))
	wc := bytes.NewBuffer(nil)
	wbuf := bufio.NewWriter(wc)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		line = strings.Replace(line, oldprojectname, newprojectname, -1)
		_, err = wbuf.WriteString(line)
		checkError(err)
	}
	checkError(wbuf.Flush())
	//fmt.Println(string(wc.Bytes()))
	checkError(ioutil.WriteFile(path, wc.Bytes(), 644))
}

func WalkDir(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(info.Name(), ".go") {
		return nil
	}
	changeFileModuleName(path)
	return nil
}

var oldprojectname string
var newprojectname string

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage 1: go run tools.go [your_app_name]")
		fmt.Println("usage 2: go run tools.go [old_app_name] [your_app_new_name]")
		return
	}
	if len(os.Args) == 3 {
		oldprojectname = os.Args[1]
		newprojectname = os.Args[2]
	} else {
		oldprojectname = "bfimpl"
		newprojectname = os.Args[1]
	}
	filepath.Walk(".", WalkDir)
	changeFileModuleName("go.mod")
}
