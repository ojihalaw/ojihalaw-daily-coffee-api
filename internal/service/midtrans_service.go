package service

import (
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransService struct {
	client *coreapi.Client
}

func NewMidtransService(client *coreapi.Client) *MidtransService {
	return &MidtransService{client: client}
}

func (m *MidtransService) Charge(req *coreapi.ChargeReq) (*coreapi.ChargeResponse, error) {
	resp, err := m.client.ChargeTransaction(req)
	if err != nil {
		// err dari midtrans-go memang bisa nil pointer → bungkus jadi error biasa
		return nil, err
	}
	return resp, nil
}
