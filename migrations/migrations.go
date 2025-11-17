package migrations

import "github.com/go-gormigrate/gormigrate/v2"

type Migration = gormigrate.Migration

func All() []*Migration {
	return []*Migration{
		createNextOfKin(),
		// add more migrations here
	}
}
