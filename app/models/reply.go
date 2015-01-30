package models

import (
	"time"
)

type Reply struct {
	BaseModel
	UserId    int32  `sql:"not null"`
	Body      string `sql:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
