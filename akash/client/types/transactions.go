package types

import "errors"

type Transaction struct {
	Height string           `json:"height"`
	Logs   []TransactionLog `json:"logs"`
	RawLog string           `json:"raw_log"`
}

type TransactionLog struct {
	Events []TransactionEvent `json:"events"`
}

type TransactionEvent struct {
	Type       string                     `json:"type"`
	Attributes TransactionEventAttributes `json:"attributes"`
}

type TransactionEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TransactionEventAttributes []TransactionEventAttribute

func (a TransactionEventAttributes) Get(key string) (string, error) {
	for i, attr := range a {
		if attr.Key == key {
			return a[i].Value, nil
		}
	}

	return "", errors.New("attribute not found")
}
