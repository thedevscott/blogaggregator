# Requirements
Go 1.23+

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