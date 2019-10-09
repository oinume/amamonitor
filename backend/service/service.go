package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
		rate, err := strconv.ParseFloat(gi.DiscountRate, 10)
		if err != nil {
			return fmt.Errorf("failed to parse DiscountRate: %v", err)
		}
		giftItem := model.GiftItem{
			FetchResultID:  fetchResult.ID,
			SalesPrice:     gi.SalesPrice,
			CataloguePrice: gi.CatalogPrice,
			DiscountRate:   rate,
			CreatedAt:      createdAt,
			UpdatedAt:      createdAt,
		}
		if err := giftItem.Insert(db); err != nil {
			return err
		}
	}
	return nil
}
