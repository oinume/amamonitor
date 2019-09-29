-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS fetch_result (
    id int unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`id`)
) ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS gift_item (
    id int unsigned NOT NULL AUTO_INCREMENT,
    fetch_result_id int unsigned NOT NULL,
    sales_price int unsigned NOT NULL,
    catalogue_price int unsigned NOT NULL,
    discount_ratio decimal(3,1) NOT NULL,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`id`)
) ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS gift_item;
DROP TABLE IF EXISTS fetch_result;
