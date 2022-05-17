CREATE TABLE game_items
(
  itemUUID STRING(36) NOT NULL,
  item_name STRING(MAX) NOT NULL,
  item_value NUMERIC NOT NULL,
  available_time TIMESTAMP NOT NULL,
  duration int64
)PRIMARY KEY (itemUUID);

CREATE TABLE player_items
(
  playerUUID STRING(36) NOT NULL,
  itemUUID STRING(36) NOT NULL,
  price NUMERIC NOT NULL,
  aquire_time TIMESTAMP NOT NULL,
  expires_time TIMESTAMP NOT NULL,
  visible BOOL NOT NULL DEFAULT(true),
  FOREIGN KEY (itemUUID) REFERENCES game_items (itemUUID)
) PRIMARY KEY (playerUUID, itemUUID),
    INTERLEAVE IN PARENT players ON DELETE CASCADE;
