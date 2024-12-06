package types

type Redemption struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	TransactionID uint        `gorm:"not null" json:"transaction_id"`
	Transaction   Transaction `gorm:"constraint:OnDelete:CASCADE;" json:"transaction"`
	VoucherID     uint        `gorm:"not null" json:"voucher_id"`
	Voucher       Voucher     `gorm:"constraint:OnDelete:CASCADE;" json:"voucher"`
	Quantity      int         `gorm:"not null" json:"quantity"`
	Points        int         `gorm:"not null" json:"points"`
	CreatedAt     int64       `gorm:"autoCreateTime" json:"created_at"`
}
