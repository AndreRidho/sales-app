package types

type Customer struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Email string `gorm:"unique;not null" json:"email"`
}
