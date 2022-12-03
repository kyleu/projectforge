# Types

This is a module for [Project Forge](https://projectforge.dev). It provides classes for representing persistent user records 

https://github.com/kyleu/projectforge/tree/master/module/user

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- A user profile is provided in the session, new users are saved in the database

No implementation of a User is provided. Either provide your own, or have Project Forge generate one by making a new file:

`./.projectforge/export/models/user.json`

```json
{
  "name": "user",
  "package": "user",
  "description": "A user of the system",
  "icon": "profile",
  "columns": [
    {
      "name": "id",
      "type": "uuid",
      "pk": true,
      "search": true
    },
    {
      "name": "name",
      "type": "string",
      "search": true,
      "tags": [
        "title"
      ]
    },
    {
      "name": "created",
      "type": "timestamp",
      "sqlDefault": "now()",
      "tags": [
        "created"
      ]
    },
    {
      "name": "updated",
      "type": "timestamp",
      "nullable": true,
      "sqlDefault": "now()",
      "tags": [
        "updated"
      ]
    }
  ]
}
```

Feel free to customize this model for your purposes
