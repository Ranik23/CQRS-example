package usecase_test

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"testing"

// 	"order-service/internal/domain"
// 	mockstorage "order-service/internal/repository/storage/mock"
// 	"order-service/internal/usecase"
// 	"order-service/internal/usecase/errs"
// 	"order-service/pkg/logger"
// 	mocktx "order-service/pkg/txmanager/mock"

// 	"github.com/stretchr/testify/assert"
// 	gomock "go.uber.org/mock/gomock"
// )



// func TestCreateOrderUseCase_Execute(t *testing.T) {
// 	ctx := context.Background()
// 	order := domain.Order{
// 		ID:     1,
// 		UserID: 123,
// 	}
// 	item := domain.OrderItem{
// 		Name:  "item1",
// 		Price: 100,
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		orderStorage := mockstorage.NewMockOrderStorage(ctrl)
// 		outboxStorage := mockstorage.NewMockOutboxStorage(ctrl)
// 		tx := mocktx.NewMockTxManager(ctrl)
// 		log, err := logger.NewLogger()
// 		assert.NoError(t, err)
		
// 		tx.EXPECT().Run(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
// 			return f(ctx)
// 		})

// 		orderStorage.EXPECT().SaveOrder(ctx, gomock.Any()).Return(nil)

// 		// проверим сериализуемость
// 		expectedBytes, _ := json.Marshal(order) // Changed gomock.Any() to order
// 		outboxStorage.EXPECT().CreateOutboxMessage(ctx, expectedBytes).Return(nil)


// 		useCase := usecase.NewCreateOrderUseCase(orderStorage, outboxStorage, tx, log)
// 		err = useCase.Execute(ctx, order, item)

// 		assert.NoError(t, err)
// 	})

// 	t.Run("no items", func(t *testing.T) {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		logger, err := logger.NewLogger()
// 		assert.NoError(t, err)

// 		useCase := usecase.NewCreateOrderUseCase(
// 			mockstorage.NewMockOrderStorage(ctrl),
// 			mockstorage.NewMockOutboxStorage(ctrl),
// 			mocktx.NewMockTxManager(ctrl),
// 			logger,
// 		)

// 		err = useCase.Execute(ctx, order)
// 		assert.Equal(t, errs.ErrNoItemsInOrder, err)
// 	})

// 	t.Run("fail save order", func(t *testing.T) {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		orderStorage := mockstorage.NewMockOrderStorage(ctrl)
// 		outboxStorage := mockstorage.NewMockOutboxStorage(ctrl)
// 		tx := mocktx.NewMockTxManager(ctrl)
// 		log, err := logger.NewLogger()
// 		assert.NoError(t, err)

// 		tx.EXPECT().Run(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
// 			return f(ctx)
// 		})

// 		orderStorage.EXPECT().SaveOrder(ctx, gomock.Any()).Return(errors.New("db error"))

// 		useCase := usecase.NewCreateOrderUseCase(orderStorage, outboxStorage, tx, log)
// 		err = useCase.Execute(ctx, order, item)

// 		assert.EqualError(t, err, "db error")
// 	})


// 	t.Run("fail outbox", func(t *testing.T) {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		orderStorage := mockstorage.NewMockOrderStorage(ctrl)
// 		outboxStorage := mockstorage.NewMockOutboxStorage(ctrl)
// 		tx := mocktx.NewMockTxManager(ctrl)
// 		log, err := logger.NewLogger()
// 		assert.NoError(t, err)

// 		tx.EXPECT().Run(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
// 			return f(ctx)
// 		})

// 		orderStorage.EXPECT().SaveOrder(ctx, gomock.Any()).Return(nil)
// 		outboxStorage.EXPECT().CreateOutboxMessage(ctx, gomock.Any()).Return(errors.New("outbox error"))

// 		useCase := usecase.NewCreateOrderUseCase(orderStorage, outboxStorage, tx, log)
// 		err = useCase.Execute(ctx, order, item)

// 		assert.EqualError(t, err, "outbox error")
// 	})
// }
