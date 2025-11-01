package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kadyrbaev2005/go-practice5/internal/models"
)

type Handler struct{
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler{
	return &Handler{DB: db}
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	start := time.Now()

	category := r.URL.Query().Get("category")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	sort := r.URL.Query().Get("sort")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	query := `
		SELECT p.id, p.name, c.name AS category, p.price
		FROM products p
		JOIN categories c ON p.category_id = c.id
	`
	conditions := []string{}
	args := []any{}
	argID := 1

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("c.name = $%d", argID))
		args = append(args, category)
		argID++
	}

	if minPrice != "" {
		val, err := strconv.Atoi(minPrice)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("p.price >= $%d", argID))
			args = append(args, val)
			argID++
		}
	}

	if maxPrice != "" {
		val, err := strconv.Atoi(maxPrice)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("p.price <= $%d", argID))
			args = append(args, val)
			argID++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	switch sort {
	case "price_asc":
		query += " ORDER BY p.price ASC"
	case "price_desc":
		query += " ORDER BY p.price DESC"
	default:
		query += " ORDER BY p.id"
	}

	if limit != "" {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		val, _ := strconv.Atoi(limit)
		args = append(args, val)
		argID++
	}
	if offset != "" {
		query += fmt.Sprintf(" OFFSET $%d", argID)
		val, _ := strconv.Atoi(offset)
		args = append(args, val)
		argID++
	}

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
			log.Fatal(err)
			return
		}
		products = append(products, p)
	}

	elapsed := time.Since(start)
	log.Printf("Query took %v", elapsed)
	w.Header().Set("X-Query-Time", elapsed.String())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}