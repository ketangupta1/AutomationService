package historyData

import (
	smartapigo "github.com/TredingInGo/smartapi"
)

type history struct {
	client *smartapigo.Client
}

func New(client *smartapigo.Client) history {
	return history{client: client}
}
func (h history) GetCandle(params smartapigo.CandleParams) ([]smartapigo.CandleResponse, error) {
	return h.client.GetCandleData(params)
}
