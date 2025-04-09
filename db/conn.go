package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := os.Getenv("Database URL")

	var err error
	Conn, err = pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("gagal koneksi DB: %w", err)
	}
	fmt.Println("Koneksi ke DB berhasil")
	return nil
}
