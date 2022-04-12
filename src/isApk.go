package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	// "strings"
)


func isApk (path string,f string) error {
// @TODO
	// file path/filename | grep 
	out,_ := exec.Command("ls",path).Output()
	fmt.Println(string(out))
	cmd, err := exec.Command("file", filepath.Join(path, f)+" | grep Zip").Output()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(cmd))
	return err

}

func main() {
	err := isApk(filepath.Join("..","apk"),"filexyz.apk")
	fmt.Println(err)
	fmt.Println(os.Getwd())
	os.Chdir("..")
	fmt.Println(os.Getwd())

}