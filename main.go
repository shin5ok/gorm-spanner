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
	where := flag.String("where", "", "filter conditions with where like %")
	migrate := flag.Bool("migrate", false, "migrate schema")
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

	// not work
	if *migrate {
		err := db.AutoMigrate(&model.User{})
		if err != nil {
			panic(err)
		}
	}

	if *debug {
		db.Logger = db.Logger.LogMode(logger.Info)
	}

	db.Transaction(func(tx *gorm.DB) error {
		var newUser = model.User{}

		newUser.UserId = uuid.New().String()
		newUser.UserName = "new-user"
		newUser.CreatedAt = time.Now()
		newUser.UpdatedAt = time.Now()

		if err := tx.Create(&newUser).Error; err != nil {
			fmt.Println(err)
			return err
		}

		sampleItemId := `0241e827-cc8d-4d62-b999-650e03e2ef72`
		if err := tx.Create(&model.UserItem{
			UserId: newUser.UserId,
			ItemId: sampleItemId,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	fmt.Println("users listing ------------------")
	var users []model.User

	if err := db.Debug().Find(&users).Error; err != nil {
		panic(err)
	}

	for n, user := range users {
		fmt.Println(n, user.UserId, user.UserName)
	}

	fmt.Println("user_items listing ------------------")
	var userItems []model.UserItem

	if err := db.Debug().Find(&userItems).Error; err != nil {
		panic(err)
	}

	for n, userItem := range userItems {
		fmt.Println(n, userItem.UserId, userItem.ItemId)
	}

	/*
			spanner> SELECT users.name, users.user_id,user_items.item_id FROM `users` inner join user_items on users.user_id = user_items.user_id;
		+-----------+--------------------------------------+--------------------------------------+
		| name      | user_id                              | item_id                              |
		+-----------+--------------------------------------+--------------------------------------+
		| test-user | 02355a04-bb91-45cf-98c2-522801344668 | d169f397-ba3f-413b-bc3c-a465576ef06e |
		| new-user  | 03b457bc-7481-4548-9dbe-237d83d644cf | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 0b25687c-25d2-42bd-8bd0-10b128d0b6a5 | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 0fe95ba8-d35b-4f2c-a4f3-e2366a932b63 | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 1dea0d0b-fc15-4a01-aff3-f80f6b962d37 | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 1e13c272-520e-4b42-b3aa-986f8b521bb1 | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 3aa89523-4856-4f91-bfce-6832029e31dc | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 41d2f02f-fa94-4e62-b866-0162ed2da362 | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 519e84b7-d5b6-4180-a42d-cc8a927a0e7c | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | 9bfaf1d4-da26-4567-84c0-8e612a84f75c | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | a0546f02-4f20-4d87-a038-25575b8456fc | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | abb713d6-a381-471c-9931-1b6a45f118bd | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | b88a05d2-4ffc-4aad-ad62-f12bceefad16 | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		| new-user  | c718cebf-ca53-4f85-b714-a7c4f41c242a | 0241e827-cc8d-4d62-b999-650e03e2ef72 |
		+-----------+--------------------------------------+--------------------------------------+
	*/

	/* []struct for JOIN, which are mixed from multiple struct */
	var userWithItems []struct {
		model.UserItem
		model.User
	}

	db.Debug().Table("users").
		Select("users.*, user_items.*").
		Joins("inner join user_items on users.user_id = user_items.user_id").
		Scan(&userWithItems)

	for n, userWithItem := range userWithItems {
		fmt.Println(n, userWithItem.User.UserId, userWithItem.User.UserName, userWithItem.ItemId, userWithItem.UserItem.CreatedAt)
	}

	fmt.Println("items listing ------------------")
	var items []model.Item

	var whereStr string
	if *where != "" {
		whereStr = fmt.Sprintf("%s%%", *where)
	}

	if err := db.Where("item_id like ?", whereStr).Find(&items).Error; err != nil {
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
	if result = db.Debug().Where("item_name = ?", "new").Find(&items); result.Error != nil {
		panic(result.Error)
	}
	fmt.Println("raw affected:", result.RowsAffected)
	// fmt.Println(items)

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
