// Package model contains the types for schema 'amamonitor'.
package model

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// FetchResult represents a row from 'amamonitor.fetch_result'.
type FetchResult struct {
	ID        uint      `json:"id"`         // id
	CreatedAt time.Time `json:"created_at"` // created_at
	UpdatedAt time.Time `json:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the FetchResult exists in the database.
func (fr *FetchResult) Exists() bool {
	return fr._exists
}

// Deleted provides information if the FetchResult has been deleted from the database.
func (fr *FetchResult) Deleted() bool {
	return fr._deleted
}

// Insert inserts the FetchResult to the database.
func (fr *FetchResult) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if fr._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO amamonitor.fetch_result (` +
		`created_at, updated_at` +
		`) VALUES (` +
		`?, ?` +
		`)`

	// run query
	XOLog(sqlstr, fr.CreatedAt, fr.UpdatedAt)
	res, err := db.Exec(sqlstr, fr.CreatedAt, fr.UpdatedAt)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	fr.ID = uint(id)
	fr._exists = true

	return nil
}

// Update updates the FetchResult in the database.
func (fr *FetchResult) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !fr._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if fr._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE amamonitor.fetch_result SET ` +
		`created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, fr.CreatedAt, fr.UpdatedAt, fr.ID)
	_, err = db.Exec(sqlstr, fr.CreatedAt, fr.UpdatedAt, fr.ID)
	return err
}

// Save saves the FetchResult to the database.
func (fr *FetchResult) Save(db XODB) error {
	if fr.Exists() {
		return fr.Update(db)
	}

	return fr.Insert(db)
}

// Delete deletes the FetchResult from the database.
func (fr *FetchResult) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !fr._exists {
		return nil
	}

	// if deleted, bail
	if fr._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM amamonitor.fetch_result WHERE id = ?`

	// run query
	XOLog(sqlstr, fr.ID)
	_, err = db.Exec(sqlstr, fr.ID)
	if err != nil {
		return err
	}

	// set deleted
	fr._deleted = true

	return nil
}

// FetchResultByID retrieves a row from 'amamonitor.fetch_result' as a FetchResult.
//
// Generated from index 'fetch_result_id_pkey'.
func FetchResultByID(db XODB, id uint) (*FetchResult, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, created_at, updated_at ` +
		`FROM amamonitor.fetch_result ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	fr := FetchResult{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&fr.ID, &fr.CreatedAt, &fr.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &fr, nil
}
