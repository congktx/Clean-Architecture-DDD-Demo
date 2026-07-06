package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	budgetInfra "src/src/budget/infrastructure"
	budgetPres "src/src/budget/presentation"
	budgetUse "src/src/budget/usecase"
	"src/src/config"
	sharedDomain "src/src/shared/domain"
	walletInfra "src/src/wallet/infrastructure"
	walletPres "src/src/wallet/presentation"
	walletUse "src/src/wallet/usecase"

	_ "github.com/lib/pq"
)

type consoleEventDispatcher struct{}

func (c *consoleEventDispatcher) Dispatch(events []sharedDomain.DomainEvent) error {
	for _, e := range events {
		log.Printf("[EVENT] Dispatching Event: %s | Timestamp: %d\n",
			e.EventName(), e.OccurredOn())
	}
	return nil
}

func main() {
	log.Println("Starting Clean Architecture DDD Demo Server...")
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: Failed to ping DB (is PostgreSQL running?). Error: %v", err)
	} else {
		log.Println("Connected to PostgreSQL successfully.")
	}

	if err := config.InitDB(db); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	dispatcher := &consoleEventDispatcher{}

	walletRepo := walletInfra.NewPostgresWalletRepository(db)
	walletUsecase := walletUse.NewWalletUsecase(walletRepo, dispatcher)
	walletHandler := walletPres.NewWalletHandler(walletUsecase)

	budgetRepo := budgetInfra.NewPostgresBudgetRepository(db)
	budgetUsecase := budgetUse.NewBudgetUsecase(budgetRepo, dispatcher)
	budgetHandler := budgetPres.NewBudgetHandler(budgetUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/wallets", walletHandler.CreateWallet)
	mux.HandleFunc("/wallets/expenses", walletHandler.RecordExpense)

	mux.HandleFunc("/budgets", budgetHandler.CreateBudget)
	mux.HandleFunc("/budgets/allocations", budgetHandler.AllocateFunds)

	// Swagger UI Endpoints
	mux.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		// Serve from root if run from root, or ../ if run from src/
		if _, err := os.Stat("swagger.yaml"); err == nil {
			http.ServeFile(w, r, "swagger.yaml")
		} else {
			http.ServeFile(w, r, "../swagger.yaml")
		}
	})

	mux.HandleFunc("/api-docs/", func(w http.ResponseWriter, r *http.Request) {
		html :=
			`<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="utf-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
			<title>Swagger UI</title>
			<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
			</head>
			<body>
			<div id="swagger-ui"></div>
			<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
			<script>
			window.onload = () => {
				window.ui = SwaggerUIBundle({
				url: '/swagger.yaml',
				dom_id: '#swagger-ui',
				});
			};
			</script>
			</body>
			</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	serverPort := ":8080"
	log.Printf("HTTP Server is listening on port %s...", serverPort)
	if err := http.ListenAndServe(serverPort, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
