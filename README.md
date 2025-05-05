# RESTful Server with JWT Authentication

A simple RESTful server built with Go, using Gin framework, JWT authentication, and SQLite database.

## Features

- User registration and login
- JWT-based authentication
- Protected routes
- SQLite database
- Password hashing with bcrypt
- Environment-based configuration
- Auto-generated TypeScript SDK

## Setup

1. Install Go (version 1.21 or higher)
2. Clone this repository
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Create a `.env` file in the root directory with the following variables:
   ```env
   PORT=8080
   JWT_SECRET=your-secret-key-here
   DB_PATH=app.db
   ```
   You can also set these as environment variables directly.

## Running the Server

```bash
go run .
```

The server will start on the port specified in your configuration (default: 8080).

## API Documentation

The API is documented using Swagger/OpenAPI. You can access the documentation at:
```
http://localhost:8080/swagger/index.html
```

## TypeScript SDK

The project includes an auto-generated TypeScript SDK. To generate the SDK:

1. Install the required tools:
   ```bash
   npm install -g @openapitools/openapi-generator-cli
   ```

2. Run the generation script:
   ```bash
   ./scripts/generate-sdk.sh
   ```

The SDK will be generated in the `sdk/typescript` directory.

### Using the SDK

1. Install the SDK in your TypeScript project:
   ```bash
   npm install @vhybz/api-client
   ```

2. Import and use the client:
   ```typescript
   import { Configuration, DefaultApi } from '@vhybz/api-client';

   const config = new Configuration({
     basePath: 'http://localhost:8080',
   });

   const api = new DefaultApi(config);

   // Register a new user
   const registerResponse = await api.register({
     username: 'testuser',
     password: 'testpass',
     email: 'test@example.com'
   });

   // Login
   const loginResponse = await api.login({
     username: 'testuser',
     password: 'testpass'
   });

   // Get profile (requires authentication)
   const profileResponse = await api.getProfile({
     headers: {
       Authorization: `Bearer ${loginResponse.data.token}`
     }
   });
   ```

## Configuration

The server uses environment variables for configuration. You can set these either in a `.env` file or as system environment variables.

### Available Configuration Options

- `PORT`: The port the server will listen on (default: 8080)
- `JWT_SECRET`: Secret key used for JWT token signing (required)
- `DB_PATH`: Path to the SQLite database file (default: app.db)

## API Endpoints

### Register
- **POST** `/register`
- Request body:
  ```json
  {
    "username": "your_username",
    "password": "your_password",
    "email": "your@email.com"
  }
  ```

### Login
- **POST** `/login`
- Request body:
  ```json
  {
    "username": "your_username",
    "password": "your_password"
  }
  ```
- Returns a JWT token

### Get Profile
- **GET** `/profile`
- Requires Authorization header with JWT token
- Returns user profile information

## Security

- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- Protected routes require valid JWT token
- Configuration values can be set securely through environment variables
