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
	fmt.Fprintf(w,"File Uploaded Sucessfully")
	apktool(file,w)
	// w.Header().Set("Content-Type","application/zip")
	// w.Write(archive(file.Name()))
	fmt.Fprintf(w,"Task Completed Sucessfully")

}

func uploadHandler (w http.ResponseWriter, r *http.Request) (*os.File , error){
	fPath := "apk/"
	fmt.Fprintf(w, "Uploadig File\n")
	// store received file
	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("apk")
	if err != nil {
		fmt.Fprintf(w, "Error Retrieving file %s", err)
		return nil, err
	}
	// fmt.Println("File: ", handler.Filename, "\nFile Size:", handler.Size, "\nMIME Header", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("Error Reading file", err)
		return nil, err
	}

	// Detect Content Type
	// contentType := http.DetectContentType(fileBytes)
	// fmt.Fprintf(w,contentType)

	fPath += handler.Filename
	f,err := os.OpenFile(fPath, os.O_RDONLY | os.O_CREATE,0644)
	if err != nil {
		fmt.Fprintf(w,"Error creating file %s",err)
		return nil, err
	}
	// defer f.Close()

	_,err = f.Write(fileBytes)

	if err != nil {
		fmt.Fprintf(w,"Error Writing to file %s",err)
		return nil, err
	}

	return f, err
} 


func apktool(f *os.File,w http.ResponseWriter) {
	c := "tools/"+f.Name()+"_src"
	fmt.Println(c)
	if _, err := os.Stat(c); !os.IsNotExist(err) {
		// path/to/whatever exists
		fmt.Println("Deleting existing Decoded files")
		exec.Command("cmd.exe", "/c", "rmdir", "/q", "/s", "D:/Learning/Cyfinoid/GO Codes/apktools/tools/apk/Voice_Recorder_v54.1_apkpure.com.apk_src").Output()

	}
	// java -jar apktool.jar d ../../../apk/Voice_Recorder_v54.1_apkpure.com.apk

	path := "tools"
	// vr := Voice_Recorder_v54.1_apkpure.com.apk
	SRC_PATH := f.Name() + "_src"
	cmd := "java -jar apktool.jar d ../" + f.Name() + " -o " + SRC_PATH

	cmdStruct := exec.Command("cmd.exe", "/c", cmd)
	cmdStruct.Dir = path
	fmt.Println(cmdStruct.Args)
	cmdStruct.Stdout = os.Stdout
	cmdStruct.Stderr = os.Stderr

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