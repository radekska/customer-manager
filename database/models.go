package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID              string     `gorm:"primaryKey"                                 json:"id"`
	FirstName       string     `                                                  json:"first_name"       validate:"required"`
	LastName        string     `                                                  json:"last_name"        validate:"required"`
	TelephoneNumber string     `gorm:"uniqueIndex:uniqueTelephoneNumber;size:256" json:"telephone_number" validate:"required"`
	CreatedAt       time.Time  `                                                  json:"created_at"`
	UpdatedAt       time.Time  `                                                  json:"updated_at"`
	Purchases       []Purchase `gorm:"foreignKey:CustomerID;"                     json:"-"`
	Repairs         []Repair   `gorm:"foreignKey:CustomerID;"                     json:"-"`
}

func (u *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

type Purchase struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	FrameModel   string    `                  json:"frame_model"`
	LensType     string    `                  json:"lens_type"`
	LensPower    string    `                  json:"lens_power"`
	PD           string    `                  json:"pd"`
	CustomerID   string    `gorm:"size:256"   json:"customer_id"`
	PurchaseType string    `                  json:"purchase_type"`
	PurchasedAt  time.Time `gorm:"type:date"  json:"purchased_at"`
	CreatedAt    time.Time `                  json:"created_at"`
	UpdatedAt    time.Time `                  json:"updated_at"`
}

func (p *Purchase) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}

type Repair struct {
	ID          string    `gorm:"primaryKey"  json:"id"`
	Description string    `                   json:"description"`
	Cost        float64   `gorm:"precision:2" json:"cost"`
	CustomerID  string    `gorm:"size:256"    json:"customer_id"`
	CreatedAt   time.Time `                   json:"created_at"`
}

func (r *Repair) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	return
}
