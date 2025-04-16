package tablemanager

import (
	"fmt"
	"strings"
)

type TableContacts struct {
	*DBTable
}

func (t *TableContacts) GetTable() *DBTable {
	return t.DBTable
}

func (t *TableContacts) GetCreateTable() string {
	var query strings.Builder

	query.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", t.TableName))
	query.WriteString("(")
	query.WriteString(" id VARCHAR(255) NOT NULL  PRIMARY KEY,")
	query.WriteString(" name VARCHAR(255) NOT NULL DEFAULT '',")
	query.WriteString(" surname VARCHAR(255) NOT NULL DEFAULT '',")
	query.WriteString(" company VARCHAR(255) NOT NULL DEFAULT '',")
	query.WriteString(" person_type ENUM('customer', 'lead') NOT NULL DEFAULT 'customer',")
	query.WriteString(" updated_date_time int unsigned,")
	query.WriteString(" config JSON NOT NULL")
	query.WriteString(") ENGINE=InnoDB")

	return query.String()
}

func (t *TableContacts) GetAlterTable(rev float64) ([]string, float64, error) {
	return nil, 0, nil
}
