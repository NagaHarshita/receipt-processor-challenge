// models/receipt.go
package models

type Receipt struct {
    Retailer      string  `json:"retailer" validate:"required"`
    PurchaseDate  string  `json:"purchaseDate" validate:"required,date"`
    PurchaseTime  string  `json:"purchaseTime" validate:"required,time"`
    Items         []Item  `json:"items" validate:"required,min=1,dive"`
    Total         string  `json:"total" validate:"required,regexp=^\\d+\\.\\d{2}$"`
}
