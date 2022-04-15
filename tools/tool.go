package tools

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aasimakhtar/apktools/filehandler"
)

// apktool
func Apktool(f *os.File,w http.ResponseWriter) error {
	// @TODO f.Name() already contains path eg: "apk/myUploadedFile.apk".
	// Path where apk files are stored. 
	// apk_path := "apk"

	// Path to where the tools are stored.
	// path := "tools"

	dir := filepath.Join("Decompiled Files")
	// os.MkdirAll(dir,0444)
	// Constructing folder name to store apktool output
	SRC_DIR := f.Name() + "_src"

	// Deletes if folder already exists, apktool fails if the folder exists
	err := filehandler.CheckFolder(dir, SRC_DIR)
	if err != nil {
		fmt.Println("Error deleting folder", err)
	}
	fmt.Println(f.Name())

	// CMD 1
	// cmd := "java -jar apktool.jar d ../" + f.Name() + " -o " + SRC_DIR
	// cmd := "java"
	args := "d " + filepath.Join("..", "apk", f.Name()) + " -o " + SRC_DIR
	// argsSlice := strings.Split(args," ")
	// CMD 2 , can also run with this instead of CMD 1
	// cmd := "apktool.bat d ../" + f.Name() + " -o " + SRC_DIR
	
	var cmdStruct *exec.Cmd
	if runtime.GOOS == "windows" {
		// @TODO fix windows cmd
		cmdStruct = exec.Command("apktool",strings.Split(args, " ")...)
	}
	if runtime.GOOS == "linux" {
		cmdStruct = exec.Command("apktool", strings.Split(args, " ")...)
	 
	}

	// In case of CMD 1, without the cmdStruct.Dir = path, cmdStruct.Wait() returns: "Error: Unable to access jarfile apktool.jar"
	// In case of CMD 2, without the cmdStruct.Dir = path, cmdStruct.Stderr [afaik] returns: "Input file (../apk\myUploadedFile.apk) was not found or was not readable."
	
	/*CMD 2 command COULD ALSO be ran by doing the following changes:
From cmd remove "../" and set cmdStruct.Dir= apk_path i.e. to "apk"
From cmd remove f.Name() and hardcode the filename, this is because f.Name() returns "apk/myUploadedFile.apk" instead of "myUploadedFile.apk"*/
	cmdStruct.Dir = dir
	// cmdStruct.Path = path
	fmt.Println(cmdStruct.Args)
	fmt.Println(os.Getwd())
	// Connecting Output and Error to commandline
	cmdStruct.Stdout = os.Stdout
	cmdStruct.Stderr = os.Stderr

	err = cmdStruct.Start()
	if err != nil {
		fmt.Println("Unable to start apktool", err)
		return err
	}
	fmt.Println(os.Getwd())

	// fmt.Println(string(out))

	err = cmdStruct.Wait()
	if err != nil {
		fmt.Println("apktool completion error", err)
		return err
	}
	fmt.Println("Reached End of Command")
	return nil
}
