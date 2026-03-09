// internal/handlers/handlers.go
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"sqlTr/internal/domain"
	generator "sqlTr/internal/functions"
	"sqlTr/internal/repository"
	"strconv"
)

type Handlers struct {
	userRepo  *repository.UserRepository
	templates map[string]*template.Template
}

// Общая структура для всех шаблонов с поддержкой сообщений об успехе
type TemplateData struct {
	Title       string
	Success     string
	Error       string
	Users       []domain.User
	User        *domain.User
	UserID      string
	Found       bool
	SearchQuery string
}

func NewHandlers(userRepo *repository.UserRepository) *Handlers {
	h := &Handlers{
		userRepo:  userRepo,
		templates: make(map[string]*template.Template),
	}

	h.loadTemplates()
	return h
}

func (h *Handlers) loadTemplates() {
	// Парсим все шаблоны
	templateNames := []string{"index", "create_user", "add_payment", "search"}

	for _, name := range templateNames {
		tmpl, err := template.ParseFiles(
			"templates/base.html",
			"templates/"+name+".html",
		)
		if err != nil {
			log.Printf("Error parsing template %s: %v", name, err)
			continue
		}
		h.templates[name] = tmpl
	}
}

func (h *Handlers) renderTemplate(w http.ResponseWriter, name string, data TemplateData) {
	tmpl, ok := h.templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Главная страница
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	// Получаем список последних пользователей
	users, err := h.userRepo.List(10)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		users = []domain.User{}
	}

	// Проверяем наличие сообщения об успехе в query параметрах
	success := ""
	if r.URL.Query().Get("success") == "user_created" {
		success = "Пользователь успешно создан! ID: " + r.URL.Query().Get("id")
	} else if r.URL.Query().Get("success") == "payment_added" {
		success = "Платеж успешно добавлен! ID: " + r.URL.Query().Get("id")
	}

	data := TemplateData{
		Title:   "Главная страница",
		Users:   users,
		Success: success,
	}

	h.renderTemplate(w, "index", data)
}

// Страница создания пользователя
func (h *Handlers) CreateUserPage(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Создание пользователя",
	}
	h.renderTemplate(w, "create_user", data)
}

// Обработка создания пользователя
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Получаем данные из формы
	age, _ := strconv.Atoi(r.FormValue("age"))
	email := r.FormValue("email")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")

	// Если поля не заполнены, генерируем случайные значения
	if age == 0 {
		age = generator.GenerateAge()
	}
	if email == "" {
		email = generator.GenerateEmail()
	}
	if firstName == "" {
		firstName = generator.GenerateName()
	}
	if lastName == "" {
		lastName = generator.GenerateName()
	}

	user := &domain.User{
		Age:       age,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	userID, err := h.userRepo.Create(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		data := TemplateData{
			Title: "Создание пользователя",
			Error: "Ошибка при создании пользователя: " + err.Error(),
		}
		h.renderTemplate(w, "create_user", data)
		return
	}

	// Перенаправляем на главную с сообщением об успехе
	http.Redirect(w, r, "/?success=user_created&id="+strconv.Itoa(userID), http.StatusSeeOther)
}

// Страница добавления платежа
func (h *Handlers) AddPaymentPage(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title:  "Добавление платежа",
		UserID: r.URL.Query().Get("user_id"),
	}
	h.renderTemplate(w, "add_payment", data)
}

// Обработка добавления платежа
func (h *Handlers) AddPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(r.FormValue("user_id"))
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	paymentName := r.FormValue("payment_name")

	// Если сумма не указана, генерируем случайную
	if amount == 0 {
		amount = generator.GenerateNumber()
	}
	// Если название платежа не указано, генерируем случайное
	if paymentName == "" {
		paymentName = generator.GeneratePaymentName()
	}

	payment := &domain.Payment{
		UserID:      userID,
		Amount:      amount,
		PaymentName: paymentName,
	}

	paymentID, err := h.userRepo.CreatePayment(payment)
	if err != nil {
		log.Printf("Error creating payment: %v", err)
		data := TemplateData{
			Title:  "Добавление платежа",
			Error:  "Ошибка при создании платежа: " + err.Error(),
			UserID: strconv.Itoa(userID),
		}
		h.renderTemplate(w, "add_payment", data)
		return
	}

	http.Redirect(w, r, "/?success=payment_added&id="+strconv.Itoa(paymentID), http.StatusSeeOther)
}

// Поиск пользователя
func (h *Handlers) SearchUser(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	if firstName == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := h.userRepo.FindByFirstName(firstName)
	if err != nil {
		log.Printf("Error searching user: %v", err)
		data := TemplateData{
			Title:       "Результаты поиска",
			Error:       "Ошибка при поиске: " + err.Error(),
			SearchQuery: firstName,
		}
		h.renderTemplate(w, "search", data)
		return
	}

	data := TemplateData{
		Title:       "Результаты поиска",
		User:        user,
		Found:       user != nil,
		SearchQuery: firstName,
	}

	h.renderTemplate(w, "search", data)
}
