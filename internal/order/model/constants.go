package model

type OrderSide string
type OrderType string
type OrderStatus string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"

	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"

	OrderStatusPending       OrderStatus = "PENDING"
	OrderStatusPartialFilled OrderStatus = "PARTIAL_FILLED"
	OrderStatusFilled        OrderStatus = "FILLED"
	OrderStatusCanceled      OrderStatus = "CANCELED"
)
