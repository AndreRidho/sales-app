package types

type Transaction struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	CustomerID  uint         `gorm:"not null" json:"customer_id"`
	Customer    Customer     `gorm:"constraint:OnDelete:CASCADE;" json:"customer"`
	Redemptions []Redemption `gorm:"foreignKey:TransactionID" json:"redemptions"`
	TotalPoints int          `gorm:"not null" json:"total_points"`
	CreatedAt   int64        `gorm:"autoCreateTime" json:"created_at"`
}
