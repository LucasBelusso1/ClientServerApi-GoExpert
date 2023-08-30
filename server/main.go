package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Exchange struct {
	ID       string
	Name     string
	Exchange string
}

type DollarExchange struct {
	Usdbrl struct {
		code       string `json:"code"`
		codein     string `json:"codein"`
		name       string `json:"name"`
		high       string `json:"high"`
		low        string `json:"low"`
		varBid     string `json:"varBid"`
		pctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		ask        string `json:"ask"`
		timestamp  string `json:"timestamp"`
		createDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", getDollarExchange)
	http.ListenAndServe(":8080", nil)
}

func getDollarExchange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

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

	var data DollarExchange
	err = json.Unmarshal(body, &data)

	if err != nil {
		panic(err)
	}

	err = persistOnDatabase(data)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data.Usdbrl)
}

func persistOnDatabase(data DollarExchange) error {
	db, err := sql.Open("sqlite3", "exchange.db")

	if err != nil {
		return err
	}
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS exchange (
			id VARCHAR(255),
			name VARCHAR(255),
			exchangeValue DECIMAL(10,2),
			PRIMARY KEY (id)
		);
	`)

	statement, _ := db.Prepare(`INSERT INTO exchange (id, name, exchangeValue) VALUES (?, ?, ?)`)

	context, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	_, err = statement.ExecContext(context, uuid.New(), data.Usdbrl.name, data.Usdbrl.Bid)

	if err != nil {
		panic(err)
	}

	rows, _ := db.Query(`SELECT * FROM exchange`)

	var exchanges []Exchange

	for rows.Next() {
		var exchange Exchange
		err = rows.Scan(&exchange.ID, &exchange.Name, &exchange.Exchange)

		if err != nil {
			return err
		}

		exchanges = append(exchanges, exchange)
	}

	fmt.Println(exchanges)

	return nil
}