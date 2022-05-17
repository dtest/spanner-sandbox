package models

import (
	"context"
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
)

type GameItem struct {
	ItemUUID      string    `json:"itemUUID"`
	ItemName      string    `json:"item_name"`
	ItemValue     big.Rat   `json:"item_value"`
	AvailableTime time.Time `json:"available_time"`
	Duration      int64     `json:"duration"`
}

func generateUUID() string {
	return uuid.NewString()
}

func (i *GameItem) Create(ctx context.Context, client spanner.Client) error {
	// Initialize game values
	i.ItemUUID = generateUUID()

	if i.AvailableTime.IsZero() {
		i.AvailableTime = time.Now()
	}

	// insert into spanner
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT game_items (itemUUID, item_name, item_value, available_time, duration) VALUES
					(@itemUUID, @itemName, @itemValue, @availableTime, @duration)
			`,
			Params: map[string]interface{}{
				"itemUUID":      i.ItemUUID,
				"itemName":      i.ItemName,
				"itemValue":     i.ItemValue,
				"availableTime": i.AvailableTime,
				"duration":      i.Duration,
			},
		}

		_, err := txn.Update(ctx, stmt)
		return err
	})

	if err != nil {
		return err
	}

	// return empty error on success
	return nil
	// return errors.New("Testing")
}

func GetItemByUUID(ctx context.Context, client spanner.Client, itemUUID string) (GameItem, error) {
	row, err := client.Single().ReadRow(ctx, "game_items",
		spanner.Key{itemUUID}, []string{"item_name", "item_value", "available_time", "duration"})
	if err != nil {
		return GameItem{}, err
	}

	item := GameItem{}
	err = row.ToStruct(&item)

	if err != nil {
		return GameItem{}, err
	}
	return item, nil
}
