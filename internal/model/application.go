package model

type Application struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index"`
	CompanyName string `gorm:"size:255"`
	Position    string `gorm:"size:255"`
	Status      string `gorm:"type:enum('applied', 'interview', 'offer', 'rejected')"`
	Notes       string `gorm:"type:text"`
	User        User   `gorm:"foreignKey:UserID"`
}
