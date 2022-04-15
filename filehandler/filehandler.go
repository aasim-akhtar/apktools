package filehandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) (*os.File, error) {
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
	f, err := os.OpenFile(handler.Filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(w, "Error creating file %s", err)
		return nil, err
	}
	defer f.Close()

	_, err = f.Write(fileBytes)

	if err != nil {
		fmt.Fprintf(w, "Error Writing to file %s", err)
		return nil, err
	}
	os.Chdir("..")
	return f, err
}


func tempfileUpload(w http.ResponseWriter, r *http.Request) (string ,error) {
	fPath := "apk"
	fmt.Fprintf(w, "Uploadig File\n")
	// store received file
	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("apk")
	if err != nil {
		fmt.Fprintf(w, "Error Retrieving file %s", err)
		return "", err
	}
	fmt.Println("File: ", handler.Filename, "\nFile Size:", handler.Size, "\nMIME Header", handler.Header)

	tempFile, err := ioutil.TempFile(fPath, "file*.apk")

	if err != nil {
		fmt.Println("Error creating file", err)
		return "",err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("Error Reading file", err)
		return "",err
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "File Uploaded sucessfully\n")

	return tempFile.Name() , nil
}

func Archive(f string) []byte {
	fmt.Println("Adding file to archive")
	fmt.Println(f)
	p := f + "_src/"
	cmdStruct := exec.Command("tools/7z.exe", "a", p)
	cmdStruct.Stdout = os.Stdout
	cmdStruct.Stderr = os.Stderr
	err := cmdStruct.Start()
	if err != nil {
		fmt.Println("Error adding to archive", err)
	}
	err = cmdStruct.Wait()
	if err != nil {
		fmt.Println(err)
	}
	return nil

}

func CheckFolder(path string, f string) error {
	c := filepath.Join(path, f)
	fmt.Println("Filepath:", c)
	if _, err := os.Stat(c); !os.IsNotExist(err) {
		// path/to/whatever exists
		fmt.Println("Deleting existing Decoded files")
		if runtime.GOOS == "windows" {
			_, err = exec.Command("cmd.exe", "/c", "rmdir", "/q", "/s", c).Output()
			return err
		} else {
			_, err = exec.Command("rm", "-rf", c).Output()
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

func IsApk(path string, f string) bool {
	// @TODO
	// file path/filename | grep Zip
	cmd, err := exec.Command("file", strings.Split(filepath.Join(path, f)+"| grep -q Zip", " ")...).Output()
	if err != nil {
		fmt.Println("Error validating apk", err)
	}

	fmt.Println(string(cmd))
	return strings.Contains(string(cmd), "zip")
}
