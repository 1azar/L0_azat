package storage

import "errors"

const (
	ORDERS_TABLE      = "orders"
	ORDER_ITEMS_TABLE = "order_items"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)
