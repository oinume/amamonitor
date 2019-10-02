package service

import (
	"context"
	"testing"
	"time"

	"github.com/oinume/amamonitor/backend/fetcher"
	"github.com/oinume/amamonitor/backend/model"
	"github.com/xo/dburl"
)

func Test_Service_CreateFetchResultGiftItems(t *testing.T) {
	type args struct {
		ctx       context.Context
		db        model.XODB
		giftItems []*fetcher.GiftItem
		createdAt time.Time
	}
	tests := map[string]struct {
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		"": {},
	}

	// TODO: config
	db, err := dburl.Open("")
	if err != nil {
		t.Fatal(err)
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := &Service{
				db: db,
			}
			// TODO: model.Transaction
			err := s.CreateFetchResultGiftItems(context.Background(), db, tt.args.giftItems, tt.args.createdAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFetchResultGiftItems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
