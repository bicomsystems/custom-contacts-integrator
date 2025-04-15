package tablemanager

import (
	"fmt"
	"strings"
)

type TableDeletedContacts struct {
	*DBTable
}

func (t *TableDeletedContacts) GetTable() *DBTable {
	return t.DBTable
}

func (t *TableDeletedContacts) GetCreateTable() string {
	var query strings.Builder

	query.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", t.TableName))
	query.WriteString("(")
	query.WriteString(" contact_id VARCHAR(255) PRIMARY KEY  NOT NULL,")
	query.WriteString(" timestamp int unsigned")
	query.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")

	return query.String()
}
