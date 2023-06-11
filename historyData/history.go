package historyData

import (
	smartapigo "github.com/TredingInGo/smartapi"
)

type History interface {
	GetCandle(params smartapigo.CandleParams) ([]smartapigo.CandleResponse, error)
}

type history struct {
	client *smartapigo.Client
}

func New(client *smartapigo.Client) History {
	return history{client: client}
}
func (h history) GetCandle(params smartapigo.CandleParams) ([]smartapigo.CandleResponse, error) {
	return h.client.GetCandleData(params)
}
