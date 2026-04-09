// Package game defines logic for the game's execution
package game

import (
	"context"

	"github.com/underark/stone-collector/internal/models/drops"
	"github.com/underark/stone-collector/internal/service/store"
	"github.com/underark/stone-collector/internal/service/ticks"
)

type GameService struct {
	s store.Store
}

func New(store store.Store) GameService {
	return GameService{
		store,
	}
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
