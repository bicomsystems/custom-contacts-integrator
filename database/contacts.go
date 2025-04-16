package database

import (
	"csm/database/tablemanager"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(database *DBDriver) *Database {
	return &Database{db: database.DB}
}

func (database *Database) GetContacts(page, limit int) ([]Contact, bool, error) {
	query := fmt.Sprintf("SELECT id, name, surname, company, person_type AS type, config FROM csm.%s LIMIT %d,%d", tablemanager.DB_TABLE_CONTACTS, limit*(page-1), limit+1)
	log.Println(query)
	rows, err := database.db.Query(query)
	if err != nil {
		log.Printf("Query failed. Error %v", err)
		return nil, false, err
	}
	defer rows.Close()

	contacts := make([]Contact, 0)
	hasMore := false
	i := 0
	for rows.Next() {
		i++
		if i > limit {
			hasMore = true
			break
		}

		c := Contact{}
		var config json.RawMessage
		if err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Company, &c.Type, &config); err != nil {
			log.Printf("Scan failed. Error %v", err)
			return nil, false, err
		}

		err = c.updateFieldsFromJsonConfig(config)
		if err != nil {
			return nil, false, err
		}

		contacts = append(contacts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	return contacts, hasMore, nil
}

func (database *Database) GetUpdatedOrDeletedContactsSinceLastSync(page, limit, timestamp int) ([]Contact, []string, bool, error) {
	query := fmt.Sprintf("SELECT id, name, surname, company, person_type AS type, config FROM csm.%s WHERE updated_date_time >= %d LIMIT %d,%d", tablemanager.DB_TABLE_CONTACTS, timestamp, limit*(page-1), limit+1)
	log.Println(query)
	rows, err := database.db.Query(query)
	if err != nil {
		log.Printf("Query failed. Error %v", err)
		return nil, nil, false, err
	}
	defer rows.Close()

	updatedContacts := make([]Contact, 0)
	deletedContactIDs := make([]string, 0)
	hasMore := false
	i := 0
	for rows.Next() {
		i++
		if i > limit {
			hasMore = true
			break
		}

		c := Contact{}
		var config json.RawMessage
		if err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Company, &c.Type, &config); err != nil {
			log.Printf("Scan failed. Error %v", err)
			return nil, nil, false, err
		}

		err = c.updateFieldsFromJsonConfig(config)
		if err != nil {
			return nil, nil, false, err
		}

		updatedContacts = append(updatedContacts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, false, err
	}

	q := fmt.Sprintf("SELECT contact_id FROM csm.%s LIMIT %d,%d", tablemanager.DB_TABLE_DELETED_CONTACTS, limit*(page-1), limit+1)
	log.Println(q)
	records, err := database.db.Query(q)
	if err != nil {
		log.Printf("Query failed. Error %v", err)
		return nil, nil, false, err
	}
	defer records.Close()

	i = 0
	hasMoreDeleted := false
	for records.Next() {
		i++
		if i > limit {
			hasMore = true
			hasMoreDeleted = true
			break
		}

		var contactID string
		if err := records.Scan(&contactID); err != nil {
			log.Printf("Scan failed. Error %v", err)
			return nil, nil, false, err
		}

		deletedContactIDs = append(deletedContactIDs, contactID)
	}

	if err := records.Err(); err != nil {
		return nil, nil, false, err
	}

	if !hasMoreDeleted {
		q = fmt.Sprintf("DELETE FROM csm.%s", tablemanager.DB_TABLE_DELETED_CONTACTS)
		log.Println(q)
		_, err = database.db.Exec(q)
		if err != nil {
			log.Printf("Delete contacts failed. Error %v", err)
			return nil, nil, false, err
		}
	}

	return updatedContacts, deletedContactIDs, hasMore, nil
}

func (c *Contact) updateFieldsFromJsonConfig(bytes json.RawMessage) error {
	if len(bytes) == 0 {
		return nil
	}

	config := Config{}
	err := json.Unmarshal(bytes, &config)
	if err == nil {
		c.Phones = config.Phones
		c.Emails = config.Emails
	}

	return err
}
