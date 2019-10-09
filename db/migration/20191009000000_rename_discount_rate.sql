-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE gift_item
  CHANGE discount_ratio `discount_rate` decimal(3,1) NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE gift_item
  CHANGE discount_rate TO `discount_ratio` decimal(3,1) NOT NULL;
