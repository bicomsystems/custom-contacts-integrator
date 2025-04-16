## General

Mysql needs to be installed on local system or wherever the app is running and to everything work out of the box set password for 'root' user for mysql to 'Test123!'.

Depending on OS, inside this repository you can found two binaries. For MAC you use ./csm-mac and for linux ./csm-linux

## Importing Contacts from CSV into MySQL

To insert data into the `contacts` table inside the `csm` database, use the following command:

```bash

mysql --local-infile=1 -u root -p -D csm \
-e "LOAD DATA LOCAL INFILE '/Users/ajnur/test/csm/bcm_contacts.csv' INTO TABLE contacts \
FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
```