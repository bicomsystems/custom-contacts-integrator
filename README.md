## General

For building this app you need to install go 1.21.1 or higher and you should just rung go build command. You will get csm executable file if everything is ok and you can run like ./csm .

## Importing Contacts from CSV into MySQL

To insert data into the `contacts` table inside the `csm` database, use the following command:

```bash

mysql --local-infile=1 -u root -p -D csm \
-e "LOAD DATA LOCAL INFILE '/Users/ajnur/test/csm/bcm_contacts.csv' INTO TABLE contacts \
FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n';"
```