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

func salvarLocal(bid string) {
	content := fmt.Sprintf("Dólar: %s\n", bid)

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Printf("Cotação salva em cotacao.txt: Dólar: %s\n", bid)
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

	salvarLocal(cotacao.Bid)
}
