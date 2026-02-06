package services

import (
	"test1/models"
	"test1/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetDailyReport() (*models.DailyReport, error) {
	today := time.Now()
	transactions, err := s.repo.GetTransactionsByDate(today)
	if err != nil {
		return nil, err
	}

	totalSales := 0
	for _, t := range transactions {
		totalSales += t.TotalAmount
	}

	return &models.DailyReport{
		Date:              today.Format("2006-01-02"),
		TotalSales:        totalSales,
		TotalTransactions: len(transactions),
		Transactions:      transactions,
	}, nil
}