package model

import "gorm.io/gorm"

func getRecord(query *gorm.DB, destination interface{}) error {
	switch err := query.First(destination).Error; err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		destination = nil
		return ErrNotFound
	default:
		destination = nil
		return err
	}
}
