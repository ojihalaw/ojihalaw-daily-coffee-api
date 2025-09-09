package service

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransService struct {
	client *coreapi.Client
}

func NewMidtransService(client *coreapi.Client) *MidtransService {
	return &MidtransService{client: client}
}

func (m *MidtransService) Charge(orderID string, amount int64) (string, error) {
	req := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
	}

	resp, err := m.client.ChargeTransaction(req)
	if err != nil {
		return "", err
	}
	return resp.TransactionID, nil
}
