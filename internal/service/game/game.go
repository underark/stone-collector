// Package game defines logic for the game's execution
package game

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/underark/stone-collector/internal/models"
	"github.com/underark/stone-collector/internal/models/drops"
	"github.com/underark/stone-collector/internal/service/store"
	"github.com/underark/stone-collector/internal/service/ticks"
)

type GameService struct {
	s store.Store
}

func New(store store.Store) *GameService {
	return &GameService{
		store,
	}
}

func (g *GameService) InsertNewUser() (string, error) {
	tx, err := g.s.GetTx()
	defer tx.Rollback(context.Background())
	if err != nil {
		return "", err
	}

	id := makeSessionID()

	user, err := store.InsertNewUser(tx, id)
	if err != nil {
		return "", err
	}

	drops := drops.Droppable()

	err = store.InsertNewUserStones(tx, drops, user)
	if err != nil {
		return "", err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return "", err
	}
	return id, nil
}

func (g *GameService) ExecuteTrade(userID int, tradeID int) error {
	tx, err := g.s.GetTx()
	defer tx.Rollback(context.Background())
	if err != nil {
		return err
	}

	trade, err := store.GetTradesForUpdate(tx, tradeID)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = store.TryTrade(tx, userID, trade)
	if err != nil {
		return err
	}

	err = store.CloseTrade(tx, tradeID)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) ProcessTicks(userID int) error {
	tx, err := g.s.GetTx()
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	last, err := store.GetTicksForUpdate(tx, userID)
	if err != nil {
		return err
	}

	t, err := ticks.TicksSince(last)
	if err != nil {
		return err
	}

	now, err := ticks.ConsumeTicks(last, t)
	if err != nil {
		return err
	}

	err = store.UpdateLastTicks(tx, userID, now)
	if err != nil {
		return err
	}

	d := drops.Drops(t)

	err = store.UpdateStones(tx, userID, d)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) GetUserState(userID int) (models.State, error) {
	state, err := g.s.GetTotalStones(userID)
	if err != nil {
		return models.State{}, err
	}

	stones, err := g.s.GetStonesByType(userID)
	if err != nil {
		return models.State{}, err
	}

	state.Stones = stones

	return state, nil
}

func (g *GameService) GetTrades(userID int) ([]models.Trade, error) {
	trades, err := g.s.GetTrades(userID)
	if err != nil {
		return make([]models.Trade, 0), err
	}

	return trades, nil
}

func (g *GameService) CreateTrade(ownerID int, offerMat string, offerAm int, reqMat string, reqAm int) error {
	err := g.s.CreateTrade(ownerID, offerMat, offerAm, reqMat, reqAm)
	if err != nil {
		return err
	}

	return nil
}

func makeSessionID() string {
	b := make([]byte, 12)
	rand.Read(b)
	val := base64.RawStdEncoding.EncodeToString(b)
	return val
}
