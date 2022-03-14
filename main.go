package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("main.go")
	startServer()
	apktool()
}

func startServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/apktool", fileUpload).Methods("POST")

	fmt.Println("SERVER STARTED AT PORT 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func apktool() {
	if _, err := os.Stat("tools/Voice_Recorder_v54.1_apkpure.com"); !os.IsNotExist(err) {
		// path/to/whatever exists
		fmt.Println("Deleting existing Decoded files")
		exec.Command("cmd.exe", "/c", "rmdir", "/q", "/s", "D:/Learning/Cyfinoid/GO Codes/apktools/tools/Voice_Recorder_v54.1_apkpure.com").Output()

	}
	// java -jar apktool.jar d ../../../apk/Voice_Recorder_v54.1_apkpure.com.apk

	path := "tools"
	cmd := "java -jar apktool.jar d ../../../apk/Voice_Recorder_v54.1_apkpure.com.apk"

	cmdStruct := exec.Command("cmd.exe", "/c", cmd)
	cmdStruct.Dir = path

	cmdStruct.Stdout = os.Stdout
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

	fmt.Fprintf(w, "File Uploaded sucessfully")

	rest_apktool(tempFile)
}

func rest_apktool(f *os.File) {

	fmt.Println("File Name:", f.Name())
}
