package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"sqlTr/internal/handlers"
	"sqlTr/internal/repository"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your pass"
	dbname   = "your db name"
)

func main() {
	// Подключение к базе данных
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging db: %v", err)
	}
	fmt.Println("Connected to database")

	// Создаем репозиторий
	userRepo := repository.NewUserRepository(db)

	// Создаем обработчики
	h := handlers.NewHandlers(userRepo)

	// Настраиваем маршруты
	http.HandleFunc("/", h.Index)
	http.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.CreateUserPage(w, r)
		} else if r.Method == http.MethodPost {
			h.CreateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/add-payment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.AddPaymentPage(w, r)
		} else if r.Method == http.MethodPost {
			h.AddPayment(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/search", h.SearchUser)

	// Обслуживание статических файлов
	staticDir := filepath.Join(".", "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// Если папки static нет, создаем ее
		os.Mkdir(staticDir, 0755)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Запускаем сервер
	serverAddr := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
