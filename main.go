package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var dataMutex sync.RWMutex // Мьютекс для управления доступом к data.txt

func main() {
	port := "20059" // Устанавливаем порт на 20059
	mux := http.NewServeMux()

	// Обслуживаем статические файлы из папки "site_c"
	fs := http.FileServer(http.Dir("./site_c"))
	mux.Handle("/", fs)

	// Обслуживаем загруженные изображения из папки "uploads"
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Обработчик для data.txt с использованием мьютекса
	mux.HandleFunc("/data.txt", dataHandler)

	// Обработчик для загрузки файлов
	mux.HandleFunc("/upload", uploadHandler)

	// Обработчик для удаления изображений
	mux.HandleFunc("/delete", deleteHandler)

	// Запускаем сервер
	log.Println("Сервер запущен на порту:" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		os.Exit(1)
	}
}

// Обработчик для загрузки файлов
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // до 10 МБ
		if err != nil {
			http.Error(w, "Ошибка при обработке данных формы", http.StatusBadRequest)
			return
		}

		// Получаем позицию
		x := r.FormValue("x")
		y := r.FormValue("y")

		// Получаем загруженный файл
		file, handler, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Ошибка при загрузке файла", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Проверяем тип файла
		buffer := make([]byte, 512)
		bytesRead, err := file.Read(buffer)
		if err != nil {
			http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
			return
		}
		contentType := http.DetectContentType(buffer[:bytesRead])
		if !strings.HasPrefix(contentType, "image/") {
			http.Error(w, "Можно загружать только изображения", http.StatusBadRequest)
			return
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			http.Error(w, "Ошибка при обработке файла", http.StatusInternalServerError)
			return
		}

		// Проверяем расширение файла
		allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
		fileExtension := strings.ToLower(filepath.Ext(handler.Filename))
		isValidExtension := false
		for _, ext := range allowedExtensions {
			if ext == fileExtension {
				isValidExtension = true
				break
			}
		}
		if !isValidExtension {
			http.Error(w, "Недопустимое расширение файла", http.StatusBadRequest)
			return
		}

		// Сохраняем файл
		uploadsDir := "./uploads"
		err = os.MkdirAll(uploadsDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Не удалось создать каталог для загрузок", http.StatusInternalServerError)
			return
		}

		filename := filepath.Base(handler.Filename)
		dstPath := filepath.Join(uploadsDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
			return
		}

		// Сохраняем данные о позиции и имени файла в data.txt
		positionData := x + "," + y + "," + filename

		dataMutex.Lock()
		defer dataMutex.Unlock()

		dataFilePath := "./data.txt"
		f, err := os.OpenFile(dataFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			http.Error(w, "Не удалось открыть файл данных", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		_, err = f.WriteString(positionData + "\n")
		if err != nil {
			http.Error(w, "Не удалось записать данные позиции", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Файл успешно загружен"))
	} else {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}

// Обработчик для data.txt с использованием мьютекса
func dataHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	http.ServeFile(w, r, "./data.txt")
}

// Обработчик для удаления изображений
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		filename := r.FormValue("filename")
		if filename == "" {
			http.Error(w, "Параметр 'filename' отсутствует", http.StatusBadRequest)
			return
		}

		dataMutex.Lock()
		defer dataMutex.Unlock()

		// Удаляем изображение
		err := os.Remove(filepath.Join("./uploads", filename))
		if err != nil {
			http.Error(w, "Не удалось удалить файл", http.StatusInternalServerError)
			return
		}

		// Обновляем data.txt
		dataBytes, err := os.ReadFile("./data.txt")
		if err != nil {
			http.Error(w, "Не удалось прочитать файл данных", http.StatusInternalServerError)
			return
		}
		lines := strings.Split(string(dataBytes), "\n")
		var newLines []string
		for _, line := range lines {
			if !strings.Contains(line, filename) && line != "" {
				newLines = append(newLines, line)
			}
		}
		err = os.WriteFile("./data.txt", []byte(strings.Join(newLines, "\n")), 0644)
		if err != nil {
			http.Error(w, "Не удалось записать файл данных", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Файл успешно удалён"))
	} else {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}
