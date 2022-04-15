package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aasimakhtar/apktools/filehandler"
	"github.com/aasimakhtar/apktools/tools"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("main.go")
	startServer()
}

func startServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/apktool", rest_apktool).Methods("POST")

	fmt.Println("SERVER STARTED AT PORT 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}


func rest_apktool(w http.ResponseWriter, r *http.Request) {

	file,err := filehandler.UploadHandler(w,r)
	if err != nil {
		fmt.Fprintf(w,"File creation failed")
		return
	}

	fmt.Println(file.Name())
	fmt.Fprintf(w,"File Uploaded Sucessfully\n")

	if runtime.GOOS == "linux" {
		if ! filehandler.IsApk("apk", file.Name()) {
			fmt.Println(file.Name(), "is not an apk file!")
			os.RemoveAll(filepath.Join("apk", file.Name()))
			return
		}		
	}
	err = tools.Apktool(file,w)
	// w.Header().Set("Content-Type","application/zip")
	// w.Write(archive(file.Name()))
	if err != nil {
		fmt.Fprintf(w,"apktool completion failed\n %s",err)
		return
	}

	fmt.Fprintf(w,"Task Completed Sucessfully\n")
}

