package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/ojihalawa/daily-coffee-api.git/internal/service"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validator       *utils.Validator
	OrderRepository *repository.OrderRepository
	Midtrans        *service.MidtransService
}

func NewOrderUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator,
	orderRepository *repository.OrderRepository, midtrans *service.MidtransService) *OrderUseCase {
	return &OrderUseCase{
		DB:              db,
		Log:             logger,
		Validator:       validator,
		OrderRepository: orderRepository,
		Midtrans:        midtrans,
	}
}

func (o *OrderUseCase) Create(ctx context.Context, request *model.CreateOrderRequest) error {
	tx := o.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := o.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(o.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	todayCount, _ := o.OrderRepository.GetTodayOrderCount(ctx, tx)
	invoiceNumber := utils.GenerateInvoice(todayCount + 1)

	// ✅ Hitung total
	var totalAmount int64
	var orderItems []entity.OrderItem
	for _, item := range request.Items {
		subtotal := item.Price * int64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, entity.OrderItem{
			ProductID:   utils.MustParseUUID(item.ProductID),
			ProductName: item.Name,
			Qty:         item.Quantity,
			Price:       item.Price,
			Subtotal:    subtotal,
		})
	}

	chargeReq, err := utils.BuildChargeReq(request, invoiceNumber, totalAmount)
	if err != nil {
		return fmt.Errorf("build charge req failed: %w", err)
	}

	resp, err := o.Midtrans.Charge(chargeReq)
	if err != nil {
		// cek apakah error dari midtrans
		if midErr, ok := err.(*midtrans.Error); ok && midErr != nil {
			o.Log.Errorf("Midtrans error: StatusCode=%v, Message=%v", midErr.StatusCode, midErr.Message)
			return fmt.Errorf("%w: midtrans error: %s", utils.ErrPayment, midErr.Message)

		}

		// fallback kalau bukan *midtrans.Error
		return fmt.Errorf("%w: midtrans error: %v", utils.ErrPayment, err)
	}
	_, deeplink := utils.ExtractMidtransURLs(resp)

	amountFloat, err := strconv.ParseFloat(resp.GrossAmount, 64)
	if err != nil {
		o.Log.Warnf("Failed to parse GrossAmount: %v", err)
	}
	amount := int64(amountFloat)

	var expiredAt *time.Time
	if resp.ExpiryTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", resp.ExpiryTime)
		if err == nil {
			expiredAt = &t
		}
	}

	// ✅ Buat entity order
	order := &entity.Order{
		UserID:        utils.MustParseUUID(request.CustomerID),
		InvoiceNumber: resp.OrderID,           // ex: INV-20250909-0003
		Status:        resp.TransactionStatus, // pending / settlement / cancel
		Amount:        amount,
		PaymentMethod: request.PaymentMethod, // ex: "e-wallet"
		PaymentType:   resp.PaymentType,      // ex: "gopay"
		TransactionID: resp.TransactionID,
		RedirectURL:   deeplink, // kalau pakai coreapi
		ExpiredAt:     expiredAt,
		OrderItems:    orderItems,
		Notes:         request.Notes,           // simpan catatan order
		ShippingAddr:  request.ShippingAddress, // simpan alamat pengiriman
	}

	if err := o.OrderRepository.Create(tx, order); err != nil {
		o.Log.Warnf("Failed create order to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		o.Log.Warnf("Failed commit transaction : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}
