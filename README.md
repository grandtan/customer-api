# customer-api

1.Install dependencies: go mod tidy
2.Database Initialization: go run cmd/initdb/main.go
3.Run the database initialization script: go run cmd/initdb/main.go
4.Start the application: go run main.go

API Endpoints

1. Create a Customer
   URL: /customers
   Method: POST
   Headers:
   Content-Type: application/json
   Body:
   {
   "name": "Alice Johnson",
   "age": 28
   }

2. Get a Customer by ID
   URL: /customers/{id}
   Method: GET

3. Update a Customer
   URL: /customers/{id}
   Method: PUT
   Headers:
   Content-Type: application/json
   Body:
   {
   "name": "Alice Johnson Updated",
   "age": 30
   }

4. Delete a Customer
   URL: /customers/{id}
   Method: DELETE

Running Tests
go test ./... -cover

Generate a coverage report:
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

Notes
Make sure to remove the test database file after running tests to avoid conflicts:
rm test_customers.db
