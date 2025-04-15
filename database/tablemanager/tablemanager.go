package tablemanager

import (
	"database/sql"
)

const (
	DB_TABLE_CONTACTS         = "contacts"
	DB_TABLE_DELETED_CONTACTS = "deleted_contacts"
)

type DBTable struct {
	TableName string
	Db        *sql.DB
}

type DBTableProvider interface {
	GetTable() *DBTable
	GetCreateTable() string
}

func GetListOfTables(db *sql.DB) []DBTableProvider {
	return []DBTableProvider{
		&TableContacts{DBTable: &DBTable{Db: db, TableName: DB_TABLE_CONTACTS}},
		&TableDeletedContacts{DBTable: &DBTable{Db: db, TableName: DB_TABLE_DELETED_CONTACTS}},
	}
}
