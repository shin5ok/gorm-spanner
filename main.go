package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	spannergorm "github.com/googleapis/go-gorm-spanner"
	_ "github.com/googleapis/go-sql-spanner"
	"github.com/shin5ok/gorm-spanner/model"
	"gorm.io/gorm"
)

var (
	projectId  = os.Getenv("PROJECT_ID")
	instanceId = os.Getenv("INSTANCE_ID")
	databaseId = os.Getenv("DATABASE_ID")
	dsn        = fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceId, databaseId)
)

func main() {
	dsnString := "projects/" + projectId + "/instances/" + instanceId + "/databases/" + databaseId
	fmt.Println("DSN: ", dsnString)

	db, err := gorm.Open(spannergorm.New(spannergorm.Config{DriverName: "spanner", DSN: dsn}),
		&gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}

	var items []model.Item

	if err := db.Where("item_id like ?", "0%").Find(&items).Error; err != nil {
		panic(err)
	}

	for n, item := range items {
		fmt.Println(n, item.ItemId, item.ItemName)
	}

	var newItem = model.Item{}

	newItem.ItemId = uuid.New().String()
	newItem.ItemName = "new"
	newItem.Price = 64000
	newItem.CreatedAt = time.Now()

	db.Create(&newItem)

	if err := db.Where("item_name = ?", "new").Find(&items).Error; err != nil {
		panic(err)
	}
	fmt.Println(items)

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("item_name = ?", "new256").Delete(&items).Error; err != nil {
			return err
		}

		var newItem = model.Item{}

		newItem.ItemId = uuid.New().String()
		newItem.ItemName = "new256"
		newItem.Price = 128000
		newItem.CreatedAt = time.Now()

		if err := tx.Debug().Create(&newItem).Error; err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err := db.Where("item_name = ?", "new256").Find(&items).Error; err != nil {
		panic(err)
	}
	fmt.Println(items)
}
