package domain

type User struct {
	ID        int
	Age       int
	FirstName string
	LastName  string
	Email     string
}

type Payment struct {
	PaymentID   int
	UserID      int
	Amount      int
	PaymentName string
}
