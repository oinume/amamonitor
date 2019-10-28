-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE gift_item
  MODIFY COLUMN `provider` ENUM('amaten', 'giftissue') NOT NULL DEFAULT 'amaten';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE gift_item
  MODIFY COLUMN `provider` ENUM('amaten') NOT NULL DEFAULT 'amaten';
