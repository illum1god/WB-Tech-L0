package models

import (
	"time"
)

type OrderInput struct {
	Order Order `json:"order" binding:"required,dive,required"`
}

type Order struct {
	OrderUID          string    `json:"order_uid" binding:"required,uuid"`
	TrackNumber       string    `json:"track_number" binding:"required"`
	Entry             string    `json:"entry" binding:"required"`
	Delivery          Delivery  `json:"delivery" binding:"required"`
	Payment           Payment   `json:"payment" binding:"required"`
	Items             []Item    `json:"items" binding:"required,dive,required"`
	Locale            string    `json:"locale" binding:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" binding:"required"`
	DeliveryService   string    `json:"delivery_service" binding:"required"`
	Shardkey          string    `json:"shardkey" binding:"required"`
	SmID              int       `json:"sm_id" binding:"required"`
	DateCreated       time.Time `json:"date_created" binding:"required"`
	OofShard          string    `json:"oof_shard" binding:"required"`
}

type Delivery struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	ZIP     string `json:"zip" binding:"required"`
	City    string `json:"city" binding:"required"`
	Address string `json:"address" binding:"required"`
	Region  string `json:"region" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" binding:"required"`
	RequestID    string `json:"request_id" binding:"required"`
	Currency     string `json:"currency" binding:"required"`
	Provider     string `json:"provider" binding:"required"`
	Amount       int    `json:"amount" binding:"required"`
	PaymentDT    int64  `json:"payment_dt" binding:"required"`
	Bank         string `json:"bank" binding:"required"`
	DeliveryCost int    `json:"delivery_cost" binding:"required"`
	GoodsTotal   int    `json:"goods_total" binding:"required"`
	CustomFee    int    `json:"custom_fee" binding:"required"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" binding:"required"`
	TrackNumber string `json:"track_number" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	RID         string `json:"rid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Sale        int    `json:"sale" binding:"required"`
	Size        string `json:"size" binding:"required"`
	TotalPrice  int    `json:"total_price" binding:"required"`
	NmID        int    `json:"nm_id" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	Status      int    `json:"status" binding:"required"`
}
