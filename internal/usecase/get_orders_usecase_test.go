package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-service/internal/domain"
	"order-service/internal/usecase"
	"order-service/internal/usecase/errs"
	mockstorage "order-service/internal/repository/storage/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetOrdersUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := int64(42)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedOrders := []domain.Order{
			{ID: 1, UserID: userID},
			{ID: 2, UserID: userID},
		}

		mockStorage := mockstorage.NewMockOrderStorage(ctrl)
		mockStorage.
			EXPECT().
			GetOrdersByUserID(ctx, userID).
			Return(expectedOrders, nil)

		useCase := usecase.NewGetOrdersUseCase(mockStorage)
		orders, err := useCase.Execute(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
	})

	t.Run("invalid userID (<= 0)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mockstorage.NewMockOrderStorage(ctrl)
		// НЕ ожидаем вызова GetOrdersByUserID
		mockStorage.
			EXPECT().
			GetOrdersByUserID(gomock.Any(), gomock.Any()).
			Times(0)

		useCase := usecase.NewGetOrdersUseCase(mockStorage)
		orders, err := useCase.Execute(ctx, int64(0))

		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, errs.ErrInvalidUserID, err)
	})

	t.Run("storage error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mockstorage.NewMockOrderStorage(ctrl)
		mockStorage.
			EXPECT().
			GetOrdersByUserID(ctx, userID).
			Return(nil, errors.New("db error"))

		useCase := usecase.NewGetOrdersUseCase(mockStorage)
		orders, err := useCase.Execute(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.EqualError(t, err, "db error")
	})
}
