package models

import (
	"context"
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
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

// Helper function to read rows from Spanner.
func readRows(iter *spanner.RowIterator) ([]spanner.Row, error) {
	var rows []spanner.Row
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		rows = append(rows, *row)
	}

	return rows, nil
}

// Get list of item UUIDs
// TODO: Currently limits to 10k by default.
func GetItemUUIDs(ctx context.Context, client spanner.Client) ([]string, error) {
	ro := client.ReadOnlyTransaction()
	stmt := spanner.Statement{SQL: `SELECT itemUUID FROM game_items LIMIT 10000`}
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	itemRows, err := readRows(iter)
	if err != nil {
		return nil, err
	}

	var itemUUIDs []string

	for _, row := range itemRows {
		var iUUID string
		if err := row.Columns(&iUUID); err != nil {
			return nil, err
		}

		itemUUIDs = append(itemUUIDs, iUUID)
	}

	return itemUUIDs, nil
}

// Retrieve an item price
func GetItemPrice(ctx context.Context, txn *spanner.ReadWriteTransaction, itemUUID string) (big.Rat, error) {
	var price big.Rat

	row, err := txn.ReadRow(ctx, "game_items", spanner.Key{itemUUID}, []string{"item_value"})
	if err != nil {
		return price, err
	}

	err = row.Columns(&price)
	if err != nil {
		return price, err
	}

	return price, nil
}

func (i *GameItem) Create(ctx context.Context, client spanner.Client) error {
	// Initialize item values
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
