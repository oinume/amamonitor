-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE gift_item
  ADD COLUMN `provider` ENUM('amaten') NOT NULL DEFAULT 'amaten' AFTER `fetch_result_id`;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE gift_item
  DROP COLUMN provider;
