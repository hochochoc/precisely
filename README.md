# Precisely CRUD APIs

Building a RESTful HTTP application in Go to create, read, update, delete documents containing content

## Prerequisites
- Mysql 
- Golang v1.16

## Local development
1. Setup environment
- Create database (on local or docker MySQL instance)
- Sync dependencies
```
go mod download
```
- Copy/modify environment variables
```
cp .env.example .env
```
Notes: Replace the values of vars in `.env`
- Run migrate 
```
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -source file://./db/migration -database "mysql://username:password@tcp(host:port)/database" up
```

2. Run
```
go run main.go
```

3. Test
```
go test ./... -v
```

## APIs

### API Response
| Element | Type   | Description                |
|---------|--------|----------------------------|
| data | array or object | list of documents or a document |
| error | string | error message or empty |
| code | int | http status code of response |
| status | bool | true when there is no error |
### Listing
- Returns a list of documents wrapped in `data`
```shell
curl -X GET \
  'http://localhost:8000/documents'
```

- Status Code: 
    - `200`: successfully got all the documents
    - `500`: internal server error, ex: database error, etc...

### Get by id
- Returns a specific document wrapped in `data`
```shell
curl -X GET \
  'http://localhost:8000/documents/{id:[0-9]+}'
```
| Element      | Description | Type   | Required | Notes                                                                         |
|--------------|-------------|--------|----------|-------------------------------------------------------------------------------|
| id    | body param      | integer | required | id of document

- Status Code
    - `200`: successfully got the document by its `id`
    - `404`: the queried document is not found in database
    - `500`: internal server error, ex: database error, etc...
### Create a document
- Create a new document and return it in `data` of response
```shell
curl -X POST \
  http://localhost:8000/documents \
  -H 'content-type: application/json' \
  -d '{
	"title": "a",
    "content": {
        "header": "header",
        "data": "data"
    },
    "signee": "signee"
}'
```
| Element      | Description | Type   | Required | Notes                                                                         |
|--------------|-------------|--------|----------|-------------------------------------------------------------------------------|
| title    | body      | string | required | title of document |
| content | body | json | optional | content of document `{"header": "", "data": ""}` |
| signee | body | string | required | signee |

- Status Code
    - `201`: successfully created the document
    - `400`: bad request, invalid json input; eg: wrong data types, etc...
    - `500`: internal server error; eg: database error, etc...
    - `422`: invalid entity, empty `title` or `signee`
### Update a document
- Update a document and return updated one in `data`
```shell
curl -X PUT \
  http://localhost:8000/documents/{id:[0-9]+} \
  -H 'content-type: application/json' \
  -d '{
	"title": "a",
    "content": {
        "header": "header",
        "data": "data"
    },
    "signee": "signee"
}'
```
| Element      | Description | Type   | Required | Notes                                                                         |
|--------------|-------------|--------|----------|-------------------------------------------------------------------------------|
| id    | body param      | int | required | the id of document
| title    | body      | string | required | title of document |
| content | body | json | optional | content of document `{"header": "", "data": ""}` |
| signee | body | string | required | signee |

- Status Code
    - `200`: successfully updated the document by its `id`
    - `404`: the updated document is not found in database
    - `500`: internal server error, ex: database error, etc...

### Delete a document
```shell
curl -X DELETE \
  http://localhost:8000/documents/{id:[0-9]+}
```
| Element      | Description | Type   | Required | Notes                                                                         |
|--------------|-------------|--------|----------|-------------------------------------------------------------------------------|
| id    | body param      | string | required | the id of document |

- Status Code
    - `200`: successfully deleted the document by its `id`
    - `404`: the updated document is not found in database
    - `500`: internal server error, ex: database error, etc...
    - `422`: invalid entity, empty `title` or `signee`

