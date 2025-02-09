package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"loan-service/internal/pkg/emailService"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5/middleware"

	"loan-service/internal/events"
	"loan-service/internal/loan"
	"loan-service/internal/middleware/auth"
	logmw "loan-service/internal/middleware/logger"
)

func main() {
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS environment variable is not set")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "loan_events"
	}

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	publisher := events.NewKafkaPublisher([]string{kafkaBrokers}, kafkaTopic)

	emailService := emailService.NewSendGridEmailService("testing-sendgrid")

	loanRepo := loan.NewLoanRepository(db)
	loanService := loan.NewLoanService(loanRepo, publisher, emailService)
	loanHandler := loan.NewLoanHandler(loanService)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(logmw.LoggerMiddleware)
	r.Use(auth.AuthMiddleware)

	r.Route("/loans", func(r chi.Router) {
		r.Post("/", loanHandler.CreateLoan)
		r.Get("/{id}", loanHandler.GetLoan)
		r.Post("/approve", loanHandler.ApproveLoan)
		r.Post("/invest", loanHandler.InvestLoan)
		r.Post("/disburse", loanHandler.DisburseLoan)
	})

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
