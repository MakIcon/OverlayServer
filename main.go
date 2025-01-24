package main

import (
	"html/template"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Указываем путь к файлу in.html в папке site_c
	tmpl, err := template.ParseFiles("site_c/in.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	port := "20059" // Устанавливаем порт 20059
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	// Запускаем сервер
	println("Сервер запущен на http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		os.Exit(1)
	}
}
