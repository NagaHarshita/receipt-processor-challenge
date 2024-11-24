// models/item.go
package models

type Item struct {
    ShortDescription string `json:"shortDescription" validate:"required"`
    Price            string `json:"price" validate:"required,regexp=^\\d+\\.\\d{2}$"`
}
