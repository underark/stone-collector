// Package store defines methods related to the shared db connection pool
package store

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/underark/stone-collector/internal/models"
)

type Store struct {
	pool *pgxpool.Pool
}

func New() (Store, error) {
	s, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error creating connection pool: %s\n", err.Error())
	}

	store := Store{s}
	return store, nil
}

func (s *Store) GetTx() (pgx.Tx, error) {
	tx, err := s.pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func GetTicksForUpdate(tx pgx.Tx, userID int) (string, error) {
	rows, err := tx.Query(context.Background(), "SELECT last_tick::text FROM users WHERE id = $1 FOR UPDATE;", userID)
	if err != nil {
		return "", nil
	}

	var s string
	if rows.Next() {
		rows.Scan(s)
	}

	return s, nil
}

func UpdateLastTicks(tx pgx.Tx, userID int, newTicks string) error {
	_, err := tx.Exec(context.Background(), "UPDATE users SET last_tick = $2 WHERE id = $1", userID, newTicks)
	if err != nil {
		return err
	}

	return nil
}

func UpdateStones(tx pgx.Tx, userID int, drops []models.Drop) error {
	for _, d := range drops {
		_, err := tx.Exec(context.Background(), "UPDATE stones SET amount = amount + $3 WHERE owner_id = $1 AND material = $2;", userID, d.Material, d.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) GetTotalStones(userID int) (models.State, error) {
	rows, err := s.pool.Query(context.Background(), "SELECT sum(amount) AS stones FROM stones WHERE owner_id = $1;", userID)
	defer func() {
		rows.Conn().Close(context.Background())
		rows.Close()
	}()

	if err != nil {
		return models.State{}, err
	}

	state, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.State])
	if err != nil {
		return models.State{}, err
	}
	return state, nil
}
