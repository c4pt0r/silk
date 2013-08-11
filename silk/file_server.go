package silk

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type FileServer struct {
	port int
	path string
}

var fs = &FileServer{
	port: 17160,
	path: "/tmp/",
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte("ERROR"))
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.Write([]byte("ERROR"))
		return
	}
	err = ioutil.WriteFile(fs.path+handler.Filename, data, 0777)
	if err != nil {
		w.Write([]byte("ERROR"))
		return
	}
	w.Write([]byte("OK"))
}

func ListenAndServe() error {
	http.HandleFunc("/upload", upload)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(fs.path))))
	err := http.ListenAndServe(":"+strconv.Itoa(fs.port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return err
	}
	return nil
}
