package fetcher

type GiftCard struct {
	discountRate string
	salesPrice   uint
}

func (gc *GiftCard) SalesPrice() uint {
	return gc.salesPrice
}

func (gc *GiftCard) DiscountRate() string {
	return gc.discountRate
}

type Client interface {
	Fetch(url string) ([]*GiftCard, error)
}

func NewClientFromURL(url string) (Client, error) {
	return nil, nil
}
