package repositories

import (
	"customer-manager/database"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type PurchaseNotFoundError struct {
	PurchaseID string
}

func (p *PurchaseNotFoundError) Error() string {
	return fmt.Sprintf("purchase with ID '%s' does not exist", p.PurchaseID)
}

type DBPurchaseRepository struct {
	DB *gorm.DB
}

func (d *DBPurchaseRepository) Create(
	customer *database.Customer,
	purchase *database.Purchase,
) (error, *database.Purchase) {
	return d.DB.Model(customer).Association("Purchases").Append(purchase), purchase
}

func (d *DBPurchaseRepository) GetAll(customerID string) (error, []database.Purchase) {
	var purchases []database.Purchase
	result := d.DB.Where("customer_id = ?", customerID).Order("purchased_at desc").Find(&purchases)
	return result.Error, purchases
}

func (d *DBPurchaseRepository) Update(purchase *database.Purchase) (error, *database.Purchase) {
	result := d.DB.Model(purchase).
		Select("FrameModel", "LensType", "LensPower", "PD", "PurchaseType", "PurchasedAt").
		Updates(purchase)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &PurchaseNotFoundError{PurchaseID: purchase.ID}, nil
	}
	return result.Error, purchase
}

func (d *DBPurchaseRepository) DeleteByID(purchaseID string) error {
	result := d.DB.Delete(&database.Purchase{ID: purchaseID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &PurchaseNotFoundError{PurchaseID: purchaseID}
	}
	return nil
}
