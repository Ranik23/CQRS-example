package postgres

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"order-service/internal/repository/storage"
	"order-service/internal/repository/storage/errs"
	"order-service/pkg/logger"
	"order-service/pkg/txmanager"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type postgresOrderStorage struct {
	txmanager txmanager.TxManager
	sq squirrel.StatementBuilderType
	logger.Logger
}

// DeleteOrder implements storage.OrderStorage.
func (p *postgresOrderStorage) DeleteOrder(ctx context.Context, orderID int64) error {
	tx := p.txmanager.Tx(ctx)

	query, args, err := p.sq.
		Delete("orders").
		Where("id = ?", orderID).
		ToSql()
	if err != nil {
		p.Logger.Errorw("Failed to build delete order query", "error", err)
		return err
	}

	cmd, err := tx.Exec(ctx, query, args...)
	if err != nil {
		p.Logger.Errorw("Failed to execute delete order query", "error", err, "query", query, "args", args)
		return err
	}

	rowsAffected := cmd.RowsAffected()
	if rowsAffected == 0 {
		p.Logger.Warnw("No order found to delete", "orderID", orderID)
		return errs.ErrNoOrderFound
	}

	return nil
}

// GetOrderByID implements storage.OrderStorage.
func (p *postgresOrderStorage) GetOrderByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	tx := p.txmanager.Tx(ctx)

	query, args, err := p.sq.
		Select("id, user_id").
		From("orders").
		Where("id = ?", orderID).
		ToSql()
	if err != nil {
		p.Logger.Errorw("Failed to build get order by ID query", "error", err)
		return nil, err
	}

	row := tx.QueryRow(ctx, query, args...)
	var order domain.Order
	if err := row.Scan(&order.ID, &order.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.Logger.Warnw("No order found with the given ID", "orderID", orderID)
			return nil, errs.ErrNoOrderFound
		}
		p.Logger.Errorw("Failed to scan order", "error", err)
		return nil, err
	}

	return &order, nil
}

// GetOrdersByUserID implements storage.OrderStorage.
func (p *postgresOrderStorage) GetOrdersByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	tx := p.txmanager.Tx(ctx)

	query, args, err := p.sq.
		Select("id, user_id").
		From("orders").
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		p.Logger.Errorw("Failed to build get orders by user ID query", "error", err)
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		p.Logger.Errorw("Failed to execute get orders by user ID query", "error", err, "query", query, "args", args)
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.UserID); err != nil {
			p.Logger.Errorw("Failed to scan order row", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// SaveOrder implements storage.OrderStorage.
func (p *postgresOrderStorage) SaveOrder(ctx context.Context, order *domain.Order) error {
	tx := p.txmanager.Tx(ctx)

	query, args, err := p.sq.
		Insert("orders").
		Columns("user_id").
		Values(order.UserID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		p.Logger.Errorw("Failed to build insert order query", "error", err)
		return err
	}

	row := tx.QueryRow(ctx, query, args...)
	if err := row.Scan(&order.ID); err != nil {
		p.Logger.Errorw("Failed to scan inserted order ID", "error", err, "query", query, "args", args)
		return err
	}

	for _, item := range order.Items {
		itemQuery, itemArgs, itemErr := p.sq.
			Insert("order_items").
			Columns("order_id", "name", "price").
			Values(order.ID, item.Name, item.Price).
			ToSql()
		if itemErr != nil {
			p.Logger.Errorw("Failed to build insert order item query", "error", itemErr)
			return itemErr
		}

		if _, err := tx.Exec(ctx, itemQuery, itemArgs...); err != nil {
			p.Logger.Errorw("Failed to execute insert order item query", "error", err, "query", itemQuery, "args", itemArgs)
			return err
		}
	}

	return nil
}

func NewPostgresOrderStorage(txmanager txmanager.TxManager) storage.OrderStorage {
	return &postgresOrderStorage{
		txmanager: txmanager,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
