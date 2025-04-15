package http

import "csm/database"

type Database interface {
	GetContacts(page, limit int) ([]database.Contact, bool, error)
	GetUpdatedOrDeletedContactsSinceLastSync(page, limit, timestamp int) ([]database.Contact, []string, bool, error)
}

type Http struct {
	db Database
}
