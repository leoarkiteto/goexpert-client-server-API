package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://localhost:8080/cotacao",
		nil,
	)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error fetching cotacao or timeout exceeded: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var cotacao CotacaoResponse
	if err := json.Unmarshal(body, &cotacao); err != nil {
		log.Fatalf("Error parsing JSON response: %v", err)
	}

	content := fmt.Sprintf("Dólar: %s\n", cotacao.Bid)
	if err := os.WriteFile("cotacao.txt", []byte(content), 0o644); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Printf("Cotação salva em cotacao.txt: Dólar: %s\n", cotacao.Bid)
}
