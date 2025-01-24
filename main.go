package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./site_c")))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	http.HandleFunc("/data.txt", dataHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":20059", nil))
}

