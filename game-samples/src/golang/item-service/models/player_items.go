package models

import (
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
)

type PlayerItem struct {
	PlayerUUID  string           `json:"playerUUID"`
	ItemUUID    string           `json:"itemUUID"`
	Price       big.Rat          `json:"price"`
	AcquireTime time.Time        `json:"acquire_time"`
	ExpiresTime spanner.NullTime `json:"expires_time"`
	Duration    int64            `json:"duration"`
}
