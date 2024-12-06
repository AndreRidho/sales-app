package types

type Voucher struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	BrandID      uint   `gorm:"not null" json:"brand_id"`
	Name         string `gorm:"not null" json:"name"`
	CostInPoints int    `gorm:"not null" json:"cost_in_points"`
	Brand        Brand  `gorm:"constraint:OnDelete:CASCADE;" json:"brand"`
}
