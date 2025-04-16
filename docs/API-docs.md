# Custom Contact Source API Documentation

This API provides access to contact information, including full listings, delta (incremental) changes, and token-based authentication.

## Authentication

This API uses **Bearer Token** authentication. You must first retrieve a token using **Basic Auth**.

### Get Token

**Endpoint:**

```
GET /token
```

**Auth:**

- **Type:** Basic Auth  
- **Username:** `test_id`  
- **Password:** `test_secret`  

**Response:**
```json
{
  "token": "<JWT_TOKEN>" //expiry should be set for this token for enhanced security
}
```

The returned `token` will be used as a Bearer token in the `Authorization` header for all other requests:

```
Authorization: Bearer <JWT_TOKEN>
```

## Get All Contacts

### Get Contacts List

This request is called on full sync (when all contacts are needed).

**Endpoint:**

```
GET /contacts?limit=2&page=1
```

**Query Parameters:**

| Param  | Type   | Description              |
|--------|--------|--------------------------|
| limit  | int    | Number of contacts to return per page |
| page   | int    | Page number              |


**Headers:**

```
Authorization: Bearer <JWT_TOKEN>
```

**Sample Response:**
```json
{
  "contacts": [
    {
      "id": "100",
      "first_name": "Victor",
      "last_name": "Donnelly",
      "company": "eos",
      "type": "lead",
      "phones": [
        { "number": "+836985101472", "label": "work" },
        ...
      ],
      "emails": [
        { "email": "OuXIVbM@CEYlhly.edu", "label": "business" },
        ...
      ]
    }
  ],
  "has_more": true/false
}
``` 
Required fields are: **first_name** and **last_name**. 

Possible values for phone labels are: `mobile, work, home, fax, other`. If label passed is not some of these values that value will be handled as `other`.

Possible values for email labels are: `private, business, other`. If label passed is not some of these values that value will be handled as `other`. 

Possible values for type: `customer, lead`. If not passed or something else is passed `customer` will be used.

## Get Delta Contacts

### Get Contacts Delta

This endpoint returns updated and deleted contacts since the given `timestamp`.

**Endpoint:**

```
GET /contacts/delta?limit=10&page=1&timestamp=1744797148
```

**Query Parameters:**

| Param     | Type | Description                                                              |
|-----------|------|------------------------------------------------------------------------- |
| limit     | int  | Number of updated contacts and deleted contactIDs to return per page     |
| page      | int  | Page number                                                              |
| timestamp | int  | A UNIX timestamp representing last sync time                             |


*Limit applies to updated contacts and deleted contactIDs. So if limit 10 is passed maximum what system should return is 10 updated contacts and 10 deleted contactsIDs.*


**Headers:**

```
Authorization: Bearer <JWT_TOKEN>
```

**Sample Response:**
```json
{
  "updated": [
    {
      "id": "100",
      "first_name": "Victor",
      "last_name": "Donnelly",
      "company": "eos",
      "type": "lead",
      "phones": [ ... ],
      "emails": [ ... ]
    }
  ],
  "deleted": ["1555", "1233", ...], //list of contactID-s
  "has_more": true/false //either more updated or deleted
}
```

Required fields are: **first_name** and **last_name**. 

Possible values for phone labels are: `mobile, work, home, fax, other`. If label passed is not some of these values that value will be handled as `other`.

Possible values for email labels are: `private, business, other`. If label passed is not some of these values that value will be handled as `other`. 

Possible values for type: `customer, lead`. If not passed or something else is passed `customer` will be used. 
