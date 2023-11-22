package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	spannergorm "github.com/googleapis/go-gorm-spanner"
	_ "github.com/googleapis/go-sql-spanner"
	"github.com/shin5ok/gorm-spanner/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	projectId  = os.Getenv("PROJECT_ID")
	instanceId = os.Getenv("INSTANCE_ID")
	databaseId = os.Getenv("DATABASE_ID")
	dsn        = fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceId, databaseId)
)

func main() {

	debug := flag.Bool("debug", false, "enable debug logging")
	prepared := flag.String("prepared", "0", "enable prepared statement")
	flag.Parse()

	if *debug {
		logger.Default = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(
		spannergorm.New(spannergorm.Config{DriverName: "spanner", DSN: dsn}),
		&gorm.Config{PrepareStmt: true},
	)
	if err != nil {
		panic(err)
	}

	if *debug {
		db.Logger = db.Logger.LogMode(logger.Info)
	}

	var items []model.Item

	var preparedStr string
	if *prepared != "" {
		preparedStr = fmt.Sprintf("%s%%", *prepared)
	}

	if err := db.Where("item_id like ?", preparedStr).Find(&items).Error; err != nil {
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

	var result *gorm.DB
	if result = db.Where("item_name = ?", "new").Find(&items); result.Error != nil {
		panic(result.Error)
	}
	fmt.Println("raw affected:", result.RowsAffected)

	fmt.Println(items)

	result = db.Take(&newItem)
	fmt.Println("Take:", newItem)
	fmt.Println("raw affected:", result.RowsAffected)

	result = db.First(&newItem)
	fmt.Println("First:", newItem)
	fmt.Println("raw affected:", result.RowsAffected)

	result = db.Last(&newItem)
	fmt.Println("Last:", newItem)
	fmt.Println("raw affected:", result.RowsAffected)

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

	for n, item := range items {
		fmt.Println(n, item.ItemId, item.ItemName)
	}
}
