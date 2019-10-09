package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/oinume/amamonitor/backend/fetcher"
	"github.com/oinume/amamonitor/backend/model"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

// Transaction is a shortcut of model.Transaction
func (s *Service) Transaction(
	ctx context.Context,
	txOptions *sql.TxOptions,
	f func(ctx context.Context, tx *sql.Tx) error,
) error {
	tx, err := s.db.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	if err := f(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateFetchResultGiftItems(
	ctx context.Context,
	db model.XODB,
	giftItems []*fetcher.GiftItem,
	createdAt time.Time,
) (*model.FetchResult, []*model.GiftItem, error) {
	fetchResult := &model.FetchResult{
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	if err := fetchResult.Insert(db); err != nil {
		return nil, nil, err
	}
	items := make([]*model.GiftItem, len(giftItems))
	for i, gi := range giftItems {
		giftItem := &model.GiftItem{
			FetchResultID:  fetchResult.ID,
			SalesPrice:     gi.SalesPrice,
			CataloguePrice: gi.CatalogPrice,
			DiscountRatio:  0,
			CreatedAt:      createdAt,
			UpdatedAt:      createdAt,
		}
		if err := giftItem.Insert(db); err != nil {
			return nil, nil, err
		}
		items[i] = giftItem
	}
	return fetchResult, items, nil
}
