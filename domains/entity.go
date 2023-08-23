package domains

import (
	"gorm.io/gorm"
	"time"
)

type (
	DBResource struct {
		ShouldUpdate bool `gorm:"-" json:"-"`
		NewlyCreated bool `gorm:"-" json:"-"`
	}
	Entity struct {
		DBResource
		ID        string          `gorm:"primaryKey;unique;column:id;default:uuid_generate_v4()" json:"id"`
		CreatedAt time.Time       `gorm:"column:created_at" json:"createdAt"`
		UpdatedAt time.Time       `gorm:"column:updated_at" json:"updatedAt"`
		DeletedAt *gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deletedAt,omitempty"`
	}
)

func (e *Entity) TableName() string {
	return "entities"
}

func (e *Entity) GetCacheKey() string {
	return e.TableName() + ":" + e.ID
}
