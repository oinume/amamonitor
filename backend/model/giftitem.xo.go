// Package model contains the types for schema 'amamonitor'.
package model

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// GiftItem represents a row from 'amamonitor.gift_item'.
type GiftItem struct {
	ID             uint      `json:"id"`              // id
	FetchResultID  uint      `json:"fetch_result_id"` // fetch_result_id
	SalesPrice     uint      `json:"sales_price"`     // sales_price
	CataloguePrice uint      `json:"catalogue_price"` // catalogue_price
	DiscountRatio  float64   `json:"discount_ratio"`  // discount_ratio
	CreatedAt      time.Time `json:"created_at"`      // created_at
	UpdatedAt      time.Time `json:"updated_at"`      // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the GiftItem exists in the database.
func (gi *GiftItem) Exists() bool {
	return gi._exists
}

// Deleted provides information if the GiftItem has been deleted from the database.
func (gi *GiftItem) Deleted() bool {
	return gi._deleted
}

// Insert inserts the GiftItem to the database.
func (gi *GiftItem) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if gi._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO amamonitor.gift_item (` +
		`fetch_result_id, sales_price, catalogue_price, discount_ratio, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, gi.FetchResultID, gi.SalesPrice, gi.CataloguePrice, gi.DiscountRatio, gi.CreatedAt, gi.UpdatedAt)
	res, err := db.Exec(sqlstr, gi.FetchResultID, gi.SalesPrice, gi.CataloguePrice, gi.DiscountRatio, gi.CreatedAt, gi.UpdatedAt)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	gi.ID = uint(id)
	gi._exists = true

	return nil
}

// Update updates the GiftItem in the database.
func (gi *GiftItem) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !gi._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if gi._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE amamonitor.gift_item SET ` +
		`fetch_result_id = ?, sales_price = ?, catalogue_price = ?, discount_ratio = ?, created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, gi.FetchResultID, gi.SalesPrice, gi.CataloguePrice, gi.DiscountRatio, gi.CreatedAt, gi.UpdatedAt, gi.ID)
	_, err = db.Exec(sqlstr, gi.FetchResultID, gi.SalesPrice, gi.CataloguePrice, gi.DiscountRatio, gi.CreatedAt, gi.UpdatedAt, gi.ID)
	return err
}

// Save saves the GiftItem to the database.
func (gi *GiftItem) Save(db XODB) error {
	if gi.Exists() {
		return gi.Update(db)
	}

	return gi.Insert(db)
}

// Delete deletes the GiftItem from the database.
func (gi *GiftItem) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !gi._exists {
		return nil
	}

	// if deleted, bail
	if gi._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM amamonitor.gift_item WHERE id = ?`

	// run query
	XOLog(sqlstr, gi.ID)
	_, err = db.Exec(sqlstr, gi.ID)
	if err != nil {
		return err
	}

	// set deleted
	gi._deleted = true

	return nil
}

// GiftItemsByFetchResultID retrieves a row from 'amamonitor.gift_item' as a GiftItem.
//
// Generated from index 'fetch_result_id'.
func GiftItemsByFetchResultID(db XODB, fetchResultID uint) ([]*GiftItem, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, fetch_result_id, sales_price, catalogue_price, discount_ratio, created_at, updated_at ` +
		`FROM amamonitor.gift_item ` +
		`WHERE fetch_result_id = ?`

	// run query
	XOLog(sqlstr, fetchResultID)
	q, err := db.Query(sqlstr, fetchResultID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*GiftItem{}
	for q.Next() {
		gi := GiftItem{
			_exists: true,
		}

		// scan
		err = q.Scan(&gi.ID, &gi.FetchResultID, &gi.SalesPrice, &gi.CataloguePrice, &gi.DiscountRatio, &gi.CreatedAt, &gi.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &gi)
	}

	return res, nil
}

// GiftItemByID retrieves a row from 'amamonitor.gift_item' as a GiftItem.
//
// Generated from index 'gift_item_id_pkey'.
func GiftItemByID(db XODB, id uint) (*GiftItem, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, fetch_result_id, sales_price, catalogue_price, discount_ratio, created_at, updated_at ` +
		`FROM amamonitor.gift_item ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	gi := GiftItem{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&gi.ID, &gi.FetchResultID, &gi.SalesPrice, &gi.CataloguePrice, &gi.DiscountRatio, &gi.CreatedAt, &gi.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &gi, nil
}
