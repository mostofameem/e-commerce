package db

func init() {
	// initialize resource type repository
	initCartTypeRepo()

	// initialize audit log repository
	initProductTypeRepo()

	initUserTypeRepo()
}
