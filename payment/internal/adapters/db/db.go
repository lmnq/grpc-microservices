package db

import (
	"fmt"

	"github.com/lmnq/grpc-microservices/payment/internal/application/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderID    int64
	TotalPrice float32
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	err := db.AutoMigrate(&Payment{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) Save(payment *domain.Payment) error {
	paymentModel := Payment{
		CustomerID: payment.CustomerID,
		Status:     payment.Status,
		OrderID:    payment.OrderID,
		TotalPrice: payment.TotalPrice,
	}
	res := a.db.Create(&paymentModel)
	if res.Error == nil {
		payment.ID = int64(paymentModel.ID)
	}
	return res.Error
}
