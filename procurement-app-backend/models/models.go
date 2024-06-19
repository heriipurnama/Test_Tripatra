package models

import "time"

type User struct {
	UserID   string `json:"userId" bson:"userId"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

type Item struct {
	ItemID   string  `json:"itemId" bson:"itemId"`
	Name     string  `json:"name" bson:"name"`
	Quantity int     `json:"quantity" bson:"quantity"`
	Price    float64 `json:"price" bson:"price"`
}

type PurchaseOrder struct {
	OrderID     string    `json:"orderId" bson:"orderId"`
	UserID      string    `json:"userId" bson:"userId"`
	Items       []Item    `json:"items" bson:"items"`
	TotalAmount float64   `json:"totalAmount" bson:"totalAmount"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type Report struct {
	ReportID    string    `json:"reportId" bson:"reportId"`
	OrderID     string    `json:"orderId" bson:"orderId"`
	GeneratedAt time.Time `json:"generatedAt" bson:"generatedAt"`
	ReportData  string    `json:"reportData" bson:"reportData"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
