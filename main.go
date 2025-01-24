package main

import (
	"log"
	"net/http"
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

// Обработчик для загрузки файлов
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Upload handler called"))
	} else {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}

// Обработчик для data.txt
func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data handler called"))
}

// Обработчик для удаления изображений
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delete handler called"))
	} else {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}
