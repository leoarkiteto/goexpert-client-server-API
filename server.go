package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

type AwesomeAPIResponse struct {
	USDBRL USDBRL `json:"USDBRL"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./cotacoes.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bid TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table (timeout or db error): %v", err)
	}

	http.HandleFunc("/cotacao", handleCotacao)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Printf("Error creating API request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error calling external API or timeout exceeded: %v", err)
		http.Error(w, "Failed to fetch exchange rate", http.StatusGatewayTimeout)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var apiResp AwesomeAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("Error parsing API JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	salvarCotacao(apiResp.USDBRL.Bid)

	respJSON := USDBRL{Bid: apiResp.USDBRL.Bid}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respJSON)
}

func salvarCotacao(bid string) {
	ctxDB := context.Background()
	ctxDB, cancelDB := context.WithTimeout(ctxDB, 10*time.Millisecond)
	defer cancelDB()

	_, err := db.ExecContext(ctxDB, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	if err != nil {
		log.Printf("Error persisting to database or timeout exceeded: %v", err)
	}
}
