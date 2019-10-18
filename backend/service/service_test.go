package service

import (
	"context"
	"database/sql"
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
					fetcher.NewGiftItem(fetcher.AmatenProvider, "93.0", 1000, 930),
					fetcher.NewGiftItem(fetcher.AmatenProvider, "89.5", 7390, 6621),
				},
				createdAt: createdAt,
			},
			wantErr: false,
		},
	}

	dbURL := model.ReplaceToTestDBURL(config.DefaultVars.XODBURL())
	db, err := dburl.Open(dbURL)
	if err != nil {
		t.Fatalf("dburl.Open failed: %v\n", err)
	}
	defer func() { _ = db.Close() }()
	s := New(db)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			var (
				gotFetchResult *model.FetchResult
				gotGiftItems   []*model.GiftItem
			)
			err := s.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
				var err error
				gotFetchResult, gotGiftItems, err = s.CreateFetchResultGiftItems(
					ctx, db, tt.args.giftItems, tt.args.createdAt,
				)
				if err != nil {
					return err
				}
				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateFetchResultGiftItems() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got, want := gotFetchResult.CreatedAt, tt.args.createdAt; got != want {
				t.Errorf("unexpected createdAt: got=%v, want=%v", got, want)
			}
			if got, want := len(gotGiftItems), len(tt.args.giftItems); got != want {
				t.Errorf("unexpected giftItems length: got=%v, want=%v", got, want)
			}
		})
	}
}
