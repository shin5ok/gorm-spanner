package main

import (
	"fmt"
	"os"

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

	db.Where("item_id like ?", "0%").Find(&items)

	fmt.Println(items)

}
