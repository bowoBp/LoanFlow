package domians

import "time"

type (
	User struct {
		ID        string    `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
		UserName  string    `gorm:"column:user_name" json:"userName"`
		Password  string    `gorm:"password"`
		CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
		UpdatedAt time.Time `gorm:"column:update_at" json:"updatedAt"`
	}
)
