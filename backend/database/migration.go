package database

import "github.com/OAthooh/BiasharaTrack.git/models"

func (d *DB) Migrate() error {
	err := d.DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return err
	}
	return nil
}
