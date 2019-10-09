package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oinume/amamonitor/backend/config"
	"github.com/oinume/amamonitor/backend/fetcher"
	"github.com/oinume/amamonitor/backend/model"
	"github.com/xo/dburl"
)

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	os.Exit(m.Run())
}

func Test_Service_CreateFetchResultGiftItems(t *testing.T) {
	type args struct {
		giftItems []*fetcher.GiftItem
		createdAt time.Time
	}

	createdAt := time.Date(2019, 9, 1, 15, 55, 20, 0, time.UTC)
	tests := map[string]struct {
		args    args
		wantErr bool
	}{
		"normal": {
			args: args{
				giftItems: []*fetcher.GiftItem{
					{DiscountRate: "93.0", CatalogPrice: 1000, SalesPrice: 930},
					{DiscountRate: "89.5", CatalogPrice: 7390, SalesPrice: 6621},
				},
				createdAt: createdAt,
			},
			wantErr: false,
		},
	}

	dbURL := model.ReplaceToTestDBURL(t, config.DefaultVars.XODBURL())
	fmt.Printf("dbURL = %v\n", dbURL)
	db, err := dburl.Open(dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	s := &Service{db: db}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// TODO: model.Transaction
			err := s.CreateFetchResultGiftItems(context.Background(), db, tt.args.giftItems, tt.args.createdAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFetchResultGiftItems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
