package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ValueDollarExchange struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var data ValueDollarExchange
	err = json.Unmarshal(body, &data)

	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if err != nil {
		file, err = os.Create("cotacao.txt")

		if err != nil {
			panic(err)
		}
	}

	defer file.Close()

	file.Write([]byte(fmt.Sprintf("Dólar: %v\n", data.Bid)))
}