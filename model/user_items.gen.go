package model

import "time"

const TableNameUserItem = "user_items"

// UserItem mapped from table <user_items>
type UserItem struct {
	UserId    string    `gorm:"column:user_id;primaryKey;unique;not null"`
	ItemId    string    `gorm:"column:item_id;primaryKey;unique;not null"`
	User      User      `gorm:"foreignKey:UserId"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

// TableName UserItem's table name
func (*UserItem) TableName() string {
	return TableNameUserItem
}
