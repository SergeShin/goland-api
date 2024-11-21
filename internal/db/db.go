package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type Item struct {
	Task   string
	Status string
}

type DB struct {
	conn *pgx.Conn
}

func New(user, password, dbname, host string, port int) (*DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect tot the db: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}
	return &DB{conn: conn}, nil
}

func (db *DB) InsertItem(ctx context.Context, item Item) error {
	query := `INSERT INTO todo_items (task, status) VALUES ($1, $2)`
	_, err := db.conn.Exec(ctx, query, item.Task, item.Status)
	return err
}

func (db *DB) GetAllItems(ctx context.Context) ([]Item, error) {
	query := `SELECT task, status FROM todo_items`
	rows, err := db.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Task, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return items, nil
}

func (db *DB) Close(ctx context.Context) {
	err := db.conn.Close(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
}
