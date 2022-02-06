package models

func MigrationTargets() []interface{} {
	return []interface{}{
		&Post{},
		&User{},
		&Tag{},
	}
}
