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

func (s *Service) CreateFetchResultGiftItems(
	ctx context.Context,
	db model.XODB,
	giftItems []*fetcher.GiftItem,
	createdAt time.Time,
) error {
	fetchResult := &model.FetchResult{
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	if err := fetchResult.Insert(db); err != nil {
		return err
	}
	for _, gi := range giftItems {
		giftItem := model.GiftItem{
			FetchResultID:  fetchResult.ID,
			SalesPrice:     gi.SalesPrice,
			CataloguePrice: gi.CatalogPrice,
			DiscountRatio:  0,
			CreatedAt:      createdAt,
			UpdatedAt:      createdAt,
		}
		if err := giftItem.Insert(db); err != nil {
			return err
		}
	}
	return nil
}
