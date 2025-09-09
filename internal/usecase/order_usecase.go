package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
	"github.com/ojihalawa/daily-coffee-api.git/internal/repository"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validator       *utils.Validator
	OrderRepository *repository.OrderRepository
}

func NewOrderUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator,
	orderRepository *repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		DB:              db,
		Log:             logger,
		Validator:       validator,
		OrderRepository: orderRepository,
	}
}

func (c *OrderUseCase) Create(ctx context.Context, request *model.CreateOrderRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(c.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	// ✅ Generate invoice number (contoh INV-20250908-0001)
	var count int64
	if err := tx.Model(&entity.Order{}).Count(&count).Error; err != nil {
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}
	invoiceNumber := fmt.Sprintf("INV-%s-%04d", utils.TodayString(), count+1)

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

	// ✅ Buat entity order
	order := &entity.Order{
		UserID:        utils.MustParseUUID(request.CustomerID),
		InvoiceNumber: invoiceNumber,
		Status:        entity.OrderStatusPending,
		Amount:        totalAmount,
		PaymentMethod: request.PaymentMethod,
		OrderItems:    orderItems,
	}

	// (opsional) Integrasi ke Midtrans Core API
	// transaction := coreapi.ChargeReq{ ... }
	// resp, err := c.Midtrans.ChargeTransaction(&transaction)
	// if err == nil {
	//     order.TransactionID = resp.TransactionID
	//     tx.Save(order)
	//     // simpan PaymentLog
	// }

	if err := c.OrderRepository.Create(tx, order); err != nil {
		c.Log.Warnf("Failed create order to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}
