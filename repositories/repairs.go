package repositories

import (
	"customer-manager/database"
	"errors"
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
	db *gorm.DB
}

func (d *DBRepairRepository) Create(customer *database.Customer, repair *database.Repair) (error, *database.Repair) {
	return d.db.Model(customer).Association("Repairs").Append(repair), repair
}

func (d *DBRepairRepository) DeleteByID(repairID string) error {
	return d.db.Delete(&database.Repair{ID: repairID}).Error
}
