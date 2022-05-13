package models

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"time"

	spanner "cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	iterator "google.golang.org/api/iterator"
)

type PlayerStats struct {
	Games_played spanner.NullInt64 `json:"games_played"`
	Games_won    spanner.NullInt64 `json:"games_won"`
}

type Player struct {
	PlayerUUID      string `json:"playerUUID"`
	Player_name     string `json:"player_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"` // not stored
	Password_hash   []byte `json:"password_hash"`
	created         time.Time
	updated         time.Time
	Stats           spanner.NullJSON `json:"stats"`
	Account_balance big.Rat          `json:"account_balance"`
	last_login      time.Time
	is_logged_in    bool
	valid_email     bool
	Current_game    string `json:"current_game"`
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

// TODO check for valid domains, and not allow local domains?
func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// TODO complexity validation
func hashPassword(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

func validatePassword(pwd string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(pwd))
}

func generateUUID() string {
	return uuid.NewString()
}

func (p *Player) AddPlayer(ctx context.Context, client spanner.Client) error {
	// Ensure email is valid
	if err := validateEmail(p.Email); err != true {
		return fmt.Errorf("New player has invalid email '%s'", p.Email)
	}

	// take supplied password+salt, hash. Store in user_password
	passHash, err := hashPassword(p.Password)

	if err != nil {
		return errors.New("Unable to hash password")
	}

	p.Password_hash = passHash

	// Generate UUIDv4
	p.PlayerUUID = generateUUID()

	// Initialize player stats
	emptyStats := spanner.NullJSON{Value: PlayerStats{
		Games_played: spanner.NullInt64{Int64: 0, Valid: true},
		Games_won:    spanner.NullInt64{Int64: 0, Valid: true},
	}, Valid: true}

	// insert into spanner
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT players (playerUUID, player_name, email, password_hash, created, stats) VALUES
					(@playerUUID, @playerName, @email, @passwordHash, CURRENT_TIMESTAMP(), @pStats)
			`,
			Params: map[string]interface{}{
				"playerUUID":   p.PlayerUUID,
				"playerName":   p.Player_name,
				"email":        p.Email,
				"passwordHash": p.Password_hash,
				"pStats":       emptyStats,
			},
		}

		_, err := txn.Update(ctx, stmt)
		return err
	})

	// todo: Handle 'AlreadyExists' errors
	if err != nil {
		return err
	}

	// return player object if successful, else return nil
	return nil
}

// TODO: Currently limits to 10k by default. This shouldn't be exposed to public API usage
func GetPlayerUUIDs(ctx context.Context, client spanner.Client) ([]string, error) {

	ro := client.ReadOnlyTransaction()
	stmt := spanner.Statement{SQL: `SELECT playerUUID FROM players LIMIT 10000`}
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	playerRows, err := readRows(iter)
	if err != nil {
		return nil, err
	}

	var playerUUIDs []string

	for _, row := range playerRows {
		var pUUID string
		if err := row.Columns(&pUUID); err != nil {
			return nil, err
		}

		playerUUIDs = append(playerUUIDs, pUUID)
	}

	return playerUUIDs, nil
}

func GetPlayerByUUID(ctx context.Context, client spanner.Client, uuid string) (Player, error) {
	row, err := client.Single().ReadRow(ctx, "players",
		spanner.Key{uuid}, []string{"email"})
	if err != nil {
		return Player{}, err
	}

	player := Player{}
	err = row.ToStruct(&player)

	if err != nil {
		return Player{}, err
	}
	return player, nil
}

// Getting player by login information
// Uses player name and password. Should return an error if no player was found
// func GetPlayerByLogin(name string, password string) (Player, error) {

// }

// Retrieves only the playerUUID and stats
func (p *Player) GetPlayerStats(ctx context.Context, client spanner.Client) error {
	row, err := client.Single().ReadRow(ctx, "players",
		spanner.Key{p.PlayerUUID}, []string{"playerUUID", "stats"})

	if err != nil {
		return err
	}

	player := Player{}
	err = row.ToStruct(&player)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
