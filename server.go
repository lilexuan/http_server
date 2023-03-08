package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func httpServer() {
	http.HandleFunc("/", printAllHandler)
	http.HandleFunc("/upload", uploadFileHandler)
	http.HandleFunc("/download/", downloadFIleHandler)
	http.ListenAndServe(":8888", nil)
}

func main() {
	httpServer()
}

func printAllHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Received a GET request\n")
	case "POST":
		fmt.Fprintf(w, "Received a POST request\n")
		r.ParseForm()
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POSTmethods are supported.\n")
	}
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only post is allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	filename := r.FormValue("filename")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// TODO: 改成读取配置文件
	uploadDir := "/tmp/upload/"
	mkdirIfNotExist(uploadDir)

	f, err := os.OpenFile(uploadDir+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully!", filename)
}

func downloadFIleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	downloadDir := "/tmp/download/"
	mkdirIfNotExist(downloadDir)
	info := strings.Split(r.RequestURI, "/")
	filename := info[len(info)-1]
	filepath := downloadDir + filename

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, filename+"is not exist!", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "read file error!", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "download error!", http.StatusInternalServerError)
		return
	}
}

func mkdirIfNotExist(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.MkdirAll(dirName, 0755)
	}
}
