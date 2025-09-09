package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
)

type CreateOrderRequest struct {
	InvoiceNumber string `json:"invoice_number"`
	TotalAmount   int64  `json:"total_amount"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	PaymentMethod string `json:"payment_method"`
	ShippingAddr  string `json:"shipping_address"`
}

func GetMidtransAuthHeader(serverKey string) string {
	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	return fmt.Sprintf("Basic %s", auth)
}

func ExtractMidtransURLs(resp *coreapi.ChargeResponse) (qrURL, deeplinkURL string) {
	for _, action := range resp.Actions {
		switch action.Name {
		case "generate-qr-code":
			qrURL = action.URL
		case "deeplink-redirect":
			deeplinkURL = action.URL
		}
	}
	return
}

// BuildChargeReq mapping payload user → coreapi.ChargeReq
func BuildChargeReq(req *model.CreateOrderRequest, invoiceNumber string, totalAmount int64) (*coreapi.ChargeReq, error) {
	transaction := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  invoiceNumber,
			GrossAmt: totalAmount,
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		},
	}

	switch req.PaymentMethod {
	case "gopay", "qris", "shopeepay":
		transaction.PaymentType = coreapi.CoreapiPaymentType(req.PaymentMethod)

	case "bca", "bni", "bri", "permata":
		transaction.PaymentType = coreapi.PaymentTypeBankTransfer

		var bank midtrans.Bank
		switch req.PaymentMethod {
		case "bca":
			bank = midtrans.BankBca
		case "bni":
			bank = midtrans.BankBni
		case "bri":
			bank = midtrans.BankBri
		case "permata":
			bank = midtrans.BankPermata
		}

		transaction.BankTransfer = &coreapi.BankTransferDetails{
			Bank: bank,
		}

	default:
		return nil, fmt.Errorf("unsupported payment method: %s", req.PaymentMethod)
	}

	return transaction, nil
}
