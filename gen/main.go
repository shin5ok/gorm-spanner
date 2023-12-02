package main

import (
	"fmt"
	"os"

	"gorm.io/gen"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	_ "github.com/googleapis/go-sql-spanner"
	"gorm.io/gorm"
)

var (
	projectId  = os.Getenv("PROJECT_ID")
	instanceId = os.Getenv("INSTANCE_ID")
	databaseId = os.Getenv("DATABASE_ID")
	dsn        = fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceId, databaseId)
)

func main() {

	db, _ := gorm.Open(spannergorm.New(spannergorm.Config{DriverName: "spanner", DSN: dsn}), &gorm.Config{})

	g := gen.NewGenerator(gen.Config{
		OutPath: "../model",
		Mode:    gen.WithDefaultQuery,
	})

	g.UseDB(db)
	g.GenerateAllTable()
	g.Execute()
}
