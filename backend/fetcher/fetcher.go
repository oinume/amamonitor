package fetcher

type Result struct {}

type Client interface {
	Fetch(url string) (*Result, error)
}

func NewClientFromURL(url string) (Client, error) {
	return nil, nil
}
