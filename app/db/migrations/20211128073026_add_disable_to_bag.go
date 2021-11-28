package migrations

import (
	"cuboid-challenge/app/models"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "20211128073026",
		Migrate: func(transaction *gorm.DB) error {
			fmt.Println("Running migration add_disable_to_bag")
			type Bag struct {
				models.Model
				Title   string
				Volume  uint
				Disable bool
			}

			return transaction.AutoMigrate(&Bag{})
		},
		Rollback: func(transaction *gorm.DB) error {
			fmt.Println("Rollback migration add_disable_to_bag")
			type Bag struct{}

			return transaction.Migrator().DropTable(&Bag{})
		},
	})
}
