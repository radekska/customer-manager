package repositories

import (
	"customer-manager/database"
	"fmt"
	"gorm.io/gorm"
)


type RepairNotFoundError struct {
	RepairID string
}

func (r *RepairNotFoundError) Error() string {
	return fmt.Sprintf("repair with ID '%s' does not exist", r.RepairID)
}

type DBRepairRepository struct {
	DB *gorm.DB
}

func (d *DBRepairRepository) GetAll(customerID string) (error, []database.Repair) {
	var repairs []database.Repair
 	result := d.DB.Where("customer_id = ?", customerID).Order("created_at desc").Find(&repairs)
	return result.Error, repairs
}


func (d *DBRepairRepository) Create(customer *database.Customer, repair *database.Repair) (error, *database.Repair) {
	return d.DB.Model(customer).Association("Repairs").Append(repair), repair
}

func (d *DBRepairRepository) DeleteByID(repairID string) error {
	return d.DB.Delete(&database.Repair{ID: repairID}).Error
}
