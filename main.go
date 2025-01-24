package main

import (
    "html/template"
    "net/http"
    "os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func main() {
    port := "20059" // Устанавливаем порт на 20059
    mux := http.NewServeMux()
    mux.HandleFunc("/", indexHandler)

    // Запускаем сервер
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        os.Exit(1)
    }
}
