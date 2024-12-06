# Sales App

This is a Go-based application for managing transactions, redemptions, and brands in a sales system. It uses Gin for the HTTP framework, GORM for database interactions, and MySQL as the database.

## Prerequisites

Before running the application, ensure you have the following installed:

- Go 1.18+ [Download Go](https://golang.org/dl/)
- MySQL or a MySQL-compatible database
- Git [Download Git](https://git-scm.com/downloads)

## Setup Instructions

### 1. Clone the Repository

Clone the repository to your local machine:

```
git clone https://github.com/AndreRidho/sales-app.git
cd sales-app
```

### 2. Set Up Environment Variables
Create a .env file in the root of the project and add the necessary environment variables. Example:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=sales_app
APP_PORT=8080
```
Note: Make sure to replace your_password with your actual database password.

### 3. Install Dependencies
Install the Go dependencies for the project:

```
go mod tidy
```
This will install the necessary packages as defined in the go.mod file.

### 4. Run the Application
To run the application locally, use the following command:

```
go run main.go
```

### 5. Running Tests
To run unit tests for the application, execute the following:

```
go test ./test/controllers -v
```
