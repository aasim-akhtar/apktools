package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("main.go")
	startServer()
}

func startServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/fileupload", fileUpload).Methods("POST")
	r.HandleFunc("/api/apktool", rest_apktool).Methods("POST")
	r.HandleFunc("/api/dummyapktool",dummyApkTool).Methods("POST")

	fmt.Println("SERVER STARTED AT PORT 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func dummyApkTool(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("tools/Voice_Recorder_v54.1_apkpure.com"); !os.IsNotExist(err) {
		// path/to/whatever exists
		fmt.Println("Deleting existing Decoded files")
		exec.Command("cmd.exe", "/c", "rmdir", "/q", "/s", "D:/Learning/Cyfinoid/GO Codes/apktools/tools/Voice_Recorder_v54.1_apkpure.com").Output()

	}
	// java -jar apktool.jar d ../../../apk/Voice_Recorder_v54.1_apkpure.com.apk

	path := "tools"
	vr := "Voice_Recorder_v54.1_apkpure.com.apk"
	cmd := "java -jar apktool.jar d ../apk/" + vr

	cmdStruct := exec.Command("cmd.exe", "/c", cmd)
	cmdStruct.Dir = path
	fmt.Println(cmdStruct.Args)
	cmdStruct.Stdout = w
	err := cmdStruct.Start()

	if err != nil {
		fmt.Println("Unable to start apktool", err)
		return
	}

	// fmt.Println(string(out))

	err = cmdStruct.Wait()
	if err != nil {
		fmt.Println("apktool completion error", err)
	}
	fmt.Println("Reached End of Command")
	fmt.Fprintf(w,"Task Completed Sucessfully")
}

func fileUpload(w http.ResponseWriter, r *http.Request) {
	fPath := "apk"
	fmt.Fprintf(w, "Uploadig File\n")
	// store received file
	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("apk")
	if err != nil {
		fmt.Fprintf(w, "Error Retrieving file %s", err)
		return
	}
	fmt.Println("File: ", handler.Filename, "\nFile Size:", handler.Size, "\nMIME Header", handler.Header)

	tempFile, err := ioutil.TempFile(fPath, "file*.apk")

	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("Error Reading file", err)
		return
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "File Uploaded sucessfully\n")

	// rest_apktool(tempFile)
}

func rest_apktool(w http.ResponseWriter, r *http.Request) {

	file,err := uploadHandler(w,r)
	if err != nil {
		fmt.Fprintf(w,"File creation failed")
		return
	}

	fmt.Println(file.Name())
	fmt.Fprintf(w,"File Uploaded Sucessfully\n")

	if runtime.GOOS == "linux" {
		if !isApk("apk", file.Name()) {
			fmt.Println(file.Name(), "is not an apk file!")
			os.RemoveAll(filepath.Join("apk", file.Name()))
			return
		}		
	}

	apktool(file,w)
	// w.Header().Set("Content-Type","application/zip")
	// w.Write(archive(file.Name()))
	fmt.Fprintf(w,"Task Completed Sucessfully\n")

}

func uploadHandler (w http.ResponseWriter, r *http.Request) (*os.File , error){
	// filepath to store apk
	fPath := "apk"
	os.Chdir(fPath)
	fmt.Println(os.Getwd())
	fmt.Fprintf(w, "Uploadig File\n")
	// store received file
	// Parse Multipart-from, set Max filesize to 50MB
	r.ParseMultipartForm(50 << 20)

	// Get Multipart.File and handler for key "apk"
	file, handler, err := r.FormFile("apk")
	if err != nil {
		fmt.Fprintf(w, "Error Retrieving file %s", err)
		return nil, err
	}
	// fmt.Println("File: ", handler.Filename, "\nFile Size:", handler.Size, "\nMIME Header", handler.Header)

	// Convert Multipart.File to []byte
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error Reading file", err)
		return nil, err
	}

	// Detect Content Type
	// contentType := http.DetectContentType(fileBytes)
	// fmt.Fprintf(w,contentType)

	// fPath = filepath.Join(fPath, handler.Filename)
	// @TODO
	// checkFile()
	f, err := os.OpenFile(handler.Filename, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(w,"Error creating file %s",err)
		return nil, err
	}
	defer f.Close()

	_,err = f.Write(fileBytes)

	if err != nil {
		fmt.Fprintf(w,"Error Writing to file %s",err)
		return nil, err
	}
	os.Chdir("..")
	return f, err
} 

// deprecated commmand example
// java -jar apktool.jar d ../../../apk/Voice_Recorder_v54.1_apkpure.com.apk

// apktool
func apktool(f *os.File,w http.ResponseWriter) {
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
	err := checkFolder(dir, SRC_DIR)
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
		return
	}
	fmt.Println(os.Getwd())

	// fmt.Println(string(out))

	err = cmdStruct.Wait()
	if err != nil {
		fmt.Println("apktool completion error", err)
	}
	fmt.Println("Reached End of Command")

}

func archive (f string) ([]byte) {
	fmt.Println("Adding file to archive")
	fmt.Println(f)
	p :=  f + "_src/"
	cmdStruct := exec.Command("tools/7z.exe","a",p)
	cmdStruct.Stdout = os.Stdout
	cmdStruct.Stderr = os.Stderr
	err := cmdStruct.Start()
	if err != nil {
		fmt.Println("Error adding to archive",err)
	}
	err = cmdStruct.Wait()
	if err != nil {
		fmt.Println(err)
	}
	return nil

}

func checkFolder (path string,f string) error{
	c := filepath.Join(path,f)
	fmt.Println("Filepath:",c)
	if _, err := os.Stat(c); !os.IsNotExist(err) {
		// path/to/whatever exists
		fmt.Println("Deleting existing Decoded files")
		if runtime.GOOS == "windows"{
			_,err = exec.Command("cmd.exe", "/c", "rmdir", "/q", "/s", c).Output()
			return err
		}else{
			_,err = exec.Command("rm","-rf",c).Output()
			return err
		}
	}
	return nil
	// err := os.Link(src, dst)
    // if err != nil {
    //     return err
    // }

    // return os.Remove(src)
}

func isApk(path string, f string) bool {
// @TODO
	// file path/filename | grep Zip
	cmd, err := exec.Command("file", strings.Split(filepath.Join(path, f)+"| grep -q Zip", " ")...).Output()
	if err != nil {
		fmt.Println("Error validating apk",err)
}

	fmt.Println(string(cmd))
	return strings.Contains(string(cmd),"zip")
}
