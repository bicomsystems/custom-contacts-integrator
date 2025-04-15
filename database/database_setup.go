package database

import (
	"csm/config"
	"csm/database/tablemanager"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBDriver struct {
	DB *sql.DB
}

func (driver *DBDriver) Connect() error {
	dbArgs := map[string]string{"address": config.Conf.MySQL.Address, "port": strconv.Itoa(config.Conf.MySQL.Port), "username": config.Conf.MySQL.User, "password": config.Conf.MySQL.Password}

	var dbparams string
	log.Println("Connecting to MySQL...")

	socket := "unix(/var/run/mysqld/mysqld.sock)"
	if len(dbArgs["address"]) > 0 {
		port, _ := strconv.Atoi(dbArgs["port"])
		socket = fmt.Sprintf("tcp(%s:%d)", dbArgs["address"], port)
	}

	dbparams = fmt.Sprintf("%s:%s@%s/", dbArgs["username"], dbArgs["password"], socket)
	db, err := sql.Open("mysql", dbparams)
	if err != nil {
		log.Println("Failed to connect MySQL")
		return err
	}

	log.Println("MySQL connected")

	db.SetMaxOpenConns(config.Conf.MySQL.MaxOpenConnections)
	db.SetMaxIdleConns(config.Conf.MySQL.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(config.Conf.MySQL.ConnMaxLifeTime) * time.Minute)

	driver.DB = db
	return nil
}

func (driver *DBDriver) Close() {
	if driver.DB != nil {
		log.Println("Closing database connection...")
		driver.DB.Close()
	}
}

func (driver *DBDriver) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := driver.DB.Query(query, args...)
	if err != nil {
		log.Printf("Failed to execute query %s. Error: %s", query, err)
	}
	return rows, err
}

func (driver *DBDriver) QueryRow(query string, args ...any) *sql.Row {
	return driver.DB.QueryRow(query, args...)
}

func (driver *DBDriver) Exec(query string, args ...any) (sql.Result, error) {
	result, err := driver.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Failed to execute query %s. Error: %s", query, err)
	}
	return result, err
}

func (driver *DBDriver) MaintainDatabasesAndTables() error {
	log.Println("Maintain databases and tables running ...")

	log.Println("Creating csm database if it not exists ...")
	if _, err := driver.DB.Exec("CREATE DATABASE IF NOT EXISTS csm"); err != nil {
		log.Printf("Create database failed. Error %v", err)
		return err
	}

	log.Print("Checking tables")

	driver.DB.Exec("USE csm")
	tables := tablemanager.GetListOfTables(driver.DB)
	for _, tp := range tables {
		t := tp.GetTable()

		//First let's create table if table not exists
		qcreate := tp.GetCreateTable()

		if _, err := t.Db.Exec(qcreate); err != nil {
			log.Printf("Create table failed. Error %v", err)
			return err
		}
	}

	q := "SELECT COUNT(*) FROM information_schema.triggers WHERE trigger_name = 'trg_after_delete_table_contacts'"
	var triggerCount int
	err := driver.DB.QueryRow(q).Scan(&triggerCount)
	if err != nil {
		log.Printf("Failed to get trigger count for our trigger trg_after_delete_table_contacts")
		return err
	}

	if triggerCount > 1 {
		log.Printf("Trigger count is bigger than 1 for our trigger trg_after_delete_table_contacts. Exiting...")
		return errors.New("trigger count not ok")
	}

	if triggerCount == 0 {
		//we will add trigger so once any row is deleted from table contacts to insert into deleted_contacts table
		triggerSQL := `
	CREATE TRIGGER trg_after_delete_table_contacts
	AFTER DELETE ON csm.contacts
	FOR EACH ROW
	BEGIN
		INSERT INTO csm.deleted_contacts (contact_id, timestamp)
		VALUES (OLD.id, UNIX_TIMESTAMP());
	END;
	`

		_, err := driver.DB.Exec(triggerSQL)
		if err != nil {
			log.Printf("Failed to create trigger: %v", err)
			return err
		}
	}

	return nil
}
