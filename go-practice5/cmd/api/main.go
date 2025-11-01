package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kadyrbaev2005/go-practice5/internal/handlers"
)

func main(){
	context := context.Background()
	dsn := "postgres://postgres:140205091205@localhost:5432/go-practice5"

	pool, err := pgxpool.New(context, dsn)
	if err != nil{
		log.Fatal(err)
	}
	defer pool.Close()
	h := handlers.NewHandler(pool)

	mux := http.NewServeMux()
	mux.Handle("/products", http.HandlerFunc(h.GetAllProducts))

	log.Fatal(http.ListenAndServe(":8080", mux))
}