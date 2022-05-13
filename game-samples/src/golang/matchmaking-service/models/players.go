package models

import spanner "cloud.google.com/go/spanner"

type PlayerStats struct {
	Games_played int `json:"games_played"`
	Games_won    int `json:"games_won"`
}

type Player struct {
	PlayerUUID   string           `json:"playerUUID"`
	Stats        spanner.NullJSON `json:"stats"`
	Current_game string           `json:"current_game"`
}
