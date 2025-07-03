# Requirements
Go 1.23+
Postgres 15+

# Config File
Located in the home directory: ~/.gatorconfig.json
Contents:
```json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

internal/config: config internal package used for reading and writing JSON file

# Postgres Database
1. [Install Postgres](https://www.postgresql.org/download/)
2. Verify install 
```bash
psql --version
```
3. Start the Postgres server [Docs](https://www.postgresql.org/docs/current/server-start.html)
* Mac: brew services start postgresql@15
* Linux: sudo service postgresql start

4. Connect to the DB
* Mac: psql postgres
* Linux: sudo -u postgres psql

5. Create the DB & connect
```bash
CREATE DATABASE <name>;

\c <name>

# set user password Linux
ALTER USER postgres PASSWORD 'postgres';
```

# Reminder before DB tests
Run goose up/down commands from inside the sql/schema directory