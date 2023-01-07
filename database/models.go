package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	ID              string `gorm:"primaryKey"`
	FirstName       string
	LastName        string
	TelephoneNumber string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Purchases       []Purchase `gorm:"foreignKey:CustomerID"`
}

func (u *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

type Purchase struct {
	ID         string `gorm:"primaryKey"`
	FrameModel string
	LensType   string
	LensPower  string
	PD         string
	CustomerID string
	Customer   Customer
	CreatedAt  time.Time
}

func (p *Purchase) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}

type Repair struct {
	ID          string `gorm:"primaryKey"`
	Description string
	Cost        float32
	CustomerID  string
	Customer    Customer `gorm:"foreignKey:CustomerID"`
	CreatedAt   time.Time
}

func (r *Repair) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	return
}
