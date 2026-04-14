// Package store defines methods related to the shared db connection pool
package store

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/underark/stone-collector/internal/models"
)

type Store struct {
	pool *pgxpool.Pool
}

func (s Store) CloseStore() {
	s.pool.Close()
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

func GetTradesForUpdate(tx pgx.Tx, tradeID string) (models.Trade, error) {
	rows, err := tx.Query(context.Background(), "SELECT * FROM trades WHERE id = $1 FOR UPDATE;", tradeID)
	if err != nil {
		return models.Trade{}, err
	}
	trade, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Trade])
	if err != nil {
		return models.Trade{}, err
	}

	return trade, nil
}

func TryTrade(tx pgx.Tx, userID int, trade models.Trade) error {
	r, err := tx.Exec(context.Background(), "UPDATE stones SET amount = amount - $1 WHERE material = $2 AND owner_id = $3 AND amount >= $1;", trade.Amount, trade.Material, trade.OwnerID)
	if err != nil {
		return err
	} else if r.RowsAffected() == 0 {
		return errors.New("error: not enough owner stones")
	}

	r, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount - $1 WHERE material = $2 AND owner_id = $3 AND amount >= $1;", trade.AmountReq, trade.MaterialReq, userID)
	if err != nil {
		return err
	} else if r.RowsAffected() == 0 {
		return errors.New("error: not enough responder stones")
	}

	_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount + $1 WHERE material = $2 AND owner_id = $3;", trade.Amount, trade.Material, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount + $1 WHERE material = $2 AND owner_id = $3;", trade.AmountReq, trade.MaterialReq, trade.OwnerID)
	if err != nil {
		return err
	}

	return nil
}

func InsertNewUser(tx pgx.Tx, session string) (int, error) {
	r, err := tx.Query(context.Background(), "INSERT INTO users (name, last_tick, session_id) VALUES ($1, NOW() at time zone 'utc', $2) RETURNING id;", "newUser", session)
	if err != nil {
		return 0, err
	}

	i, err := pgx.CollectOneRow(r, pgx.RowTo[int])
	if err != nil {
		return 0, err
	}
	return i, nil
}

func InsertNewUserStones(tx pgx.Tx, drops []string, id int) error {
	for _, d := range drops {
		_, err := tx.Exec(context.Background(), "INSERT INTO stones (owner_id, material, amount) VALUES ($1, $2, $3);", id, d, 0)
		if err != nil {
			return err
		}
	}
	return nil
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

func (s Store) GetTrades() ([]models.Trade, error) {
	rows, err := s.pool.Query(context.Background(), "SELECT * FROM trades;")
	if err != nil {
		return make([]models.Trade, 0), err
	}

	trades, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Trade])
	if err != nil {
		return make([]models.Trade, 0), err
	}

	return trades, nil
}

func (s Store) GetTotalStones(userID int) (models.State, error) {
	rows, err := s.pool.Query(context.Background(), "SELECT sum(amount) AS total FROM stones WHERE owner_id = $1;", userID)
	if err != nil {
		return models.State{}, err
	}

	state, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.State])
	if err != nil {
		return models.State{}, err
	}
	return state, nil
}

func (s Store) GetStonesByType(userID int) (stones []models.Inventory, err error) {
	stones = make([]models.Inventory, 0)
	rows, err := s.pool.Query(context.Background(), "SELECT material, amount FROM stones WHERE owner_id = $1;", userID)
	if err != nil {
		return
	}

	drops, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Inventory])
	if err != nil {
		return
	}

	stones = append(stones, drops...)

	return
}

func (s Store) GetUserFromSession(sessionID string) (int, error) {
	rows, err := s.pool.Query(context.Background(), "SELECT id FROM users WHERE session_id = $1;", sessionID)
	if err != nil {
		return 0, err
	}

	i, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, err
	}

	return i, nil
}
