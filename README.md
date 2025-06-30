# Clean Architecture User API

A robust REST API built with Go, Gin, and MongoDB following Clean Architecture principles with comprehensive authentication, authorization, and security features.

## Features

- **Clean Architecture**: Properly layered architecture with clear separation of concerns
- **Authentication & Authorization**: JWT-based auth with role-based access control
- **Security**: Password hashing, input validation, CORS, rate limiting
- **Database**: MongoDB with proper indexing and connection pooling
- **Validation**: Comprehensive input validation with custom error messages
- **Logging**: Structured logging with configurable levels
- **Error Handling**: Centralized error handling with proper HTTP status codes
- **Middleware**: Authentication, CORS, request ID, and logging middleware
- **Testing Ready**: Modular design for easy unit and integration testing

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/signin` - User login

### User Management
- `GET /api/v1/profile` - Get current user profile (Protected)
- `GET /api/v1/users/:id` - Get user by ID (Protected)
- `PUT /api/v1/users/:id` - Update user profile (Protected - Self or Admin)
- `DELETE /api/v1/users/:id` - Delete user (Protected - Self or Admin)

### Admin Only
- `GET /api/v1/admin/users` - Get all users with pagination (Admin only)

### Health Check
- `GET /health` - Health check endpoint

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Core business objects and data structures
- **Repositories**: Repository interfaces (data access contracts)
- **Services**: Service interfaces (business logic contracts)

### 2. Use Cases Layer (`internal/usecases/`)
- **User Use Case**: Business logic implementation
- Orchestrates repositories and external services
- Contains the core business rules

### 3. Infrastructure Layer (`internal/infrastructure/`)
- **Database**: MongoDB connection and configuration
- **Repositories**: Repository implementations
- **Security**: JWT, password hashing, middleware

### 4. Interface Layer (`internal/interfaces/`)
- **Handlers**: HTTP handlers (controllers)
- **Routes**: Route definitions and middleware setup
- **DTOs**: Data transfer objects for API communication

### 5. Shared Packages (`pkg/`)
- **Errors**: Custom error types and handling
- **Logger**: Structured logging utilities
- **Response**: Standardized API response format
- **Validator**: Input validation utilities

## Security Features

### Authentication
- JWT tokens with configurable expiration
- Secure password hashing with bcrypt
- Token-based authentication middleware

### Authorization
- Role-based access control (User, Admin)
- Resource-level authorization (users can only modify their own data)
- Admin-only endpoints protection

### Security Best Practices
- Input validation and sanitization
- CORS configuration
- Rate limiting capabilities
- Secure headers
- Password complexity requirements
- Database connection security

## Installation & Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd clean-architecture-go

```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Start MongoDB**
```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or use your local MongoDB installation
```

5. **Run the application**
```bash
go run cmd/api/main.go
```

## Configuration

The application uses environment variables for configuration:

- `ENVIRONMENT`: Application environment (development/production)
- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: MongoDB connection string
- `DATABASE_NAME`: MongoDB database name
- `JWT_SECRET`: Secret key for JWT tokens
- `JWT_EXPIRY_HOURS`: JWT token expiration time
- `LOG_LEVEL`: Logging level (debug/info/warn/error)
- `BCRYPT_COST`: Cost factor for password hashing

## API Usage Examples

### User Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### User Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

### Get User Profile
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Update User Profile
```bash
curl -X PUT http://localhost:8080/api/v1/users/<user-id> \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith"
  }'
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

### Validation Error Response
```json
{
  "success": false,
  "error": "Validation failed",
  "data": [
    "email is required",
    "password must be at least 8 characters long"
  ]
}
```

## Database Schema

### User Collection
```javascript
{
  _id: ObjectId,
  email: String (unique, indexed),
  username: String (unique, indexed),
  password: String (hashed),
  first_name: String,
  last_name: String,
  is_active: Boolean,
  role: String (enum: "user", "admin"),
  created_at: Date,
  updated_at: Date
}
```

## Development

### Project Structure Explanation

1. **cmd/**: Application entry points
2. **internal/**: Private application code
   - **config/**: Configuration management
   - **domain/**: Business logic and entities
   - **infrastructure/**: External dependencies
   - **interfaces/**: HTTP handlers and routes
   - **usecases/**: Application business logic
3. **pkg/**: Shared utility packages

### Adding New Features

1. Define entities in `internal/domain/entities/`
2. Create repository interfaces in `internal/domain/repositories/`
3. Implement repositories in `internal/infrastructure/repositories/`
4. Create use cases in `internal/usecases/`
5. Add handlers in `internal/interfaces/handlers/`
6. Update routes in `internal/interfaces/routes/`

## Testing

The architecture supports easy testing with dependency injection:

```go
// Example test setup
func TestUserUseCase(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockJWT := &MockJWTManager{}
    mockPassword := &MockPasswordManager{}
    
    useCase := NewUserUseCase(mockRepo, mockJWT, mockPassword)
    
    // Test your use case
}
```

## Production Considerations

1. **Environment Variables**: Set secure values for production
2. **Database**: Use MongoDB Atlas or properly configured replica set
3. **Logging**: Configure appropriate log levels
4. **Monitoring**: Add health checks and metrics
5. **HTTPS**: Use TLS certificates
6. **Rate Limiting**: Implement proper rate limiting
7. **Backup**: Set up database backups
8. **Scaling**: Consider horizontal scaling strategies
