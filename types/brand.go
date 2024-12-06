package types

type Brand struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Vouchers    []Voucher `gorm:"foreignKey:BrandID" json:"vouchers"`
}
