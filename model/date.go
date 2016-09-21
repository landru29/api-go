package model

import "time"

// Dater is the standard dater for database objects
type Dater struct {
	CreatedAt  int64 `bson:"createdAt"     json:"createdAt"`
	ModifiedAt int64 `bson:"modifiedAt"    json:"modifiedAt"`
}

// UpdateDates updates the dater
func (obj Dater) UpdateDates() {
	now := time.Now().Unix()
	if obj.CreatedAt == 0 {
		obj.CreatedAt = now
	}
	obj.ModifiedAt = now
}
