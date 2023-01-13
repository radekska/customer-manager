package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	ID              string     `gorm:"primaryKey"                 json:"-"`
	FirstName       string     `gorm:"uniqueIndex:uniquecustomer" json:"first_name"`
	LastName        string     `gorm:"uniqueIndex:uniquecustomer" json:"last_name"`
	TelephoneNumber string     `gorm:"uniqueIndex:uniquecustomer" json:"telephone_number"`
	CreatedAt       time.Time  `                                  json:"created_at"`
	UpdatedAt       time.Time  `                                  json:"updated_at"`
	Purchases       []Purchase `gorm:"foreignKey:CustomerID;"     json:"-"`
	Repairs         []Repair   `gorm:"foreignKey:CustomerID;"     json:"-"`
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
	Cost        float64 `gorm:"precision:2"`
	CustomerID  string
	Customer    Customer
	CreatedAt   time.Time
}

func (r *Repair) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	return
}
