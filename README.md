# BYOW User Service

A comprehensive user management microservice for the Build Your Own Website (BYOW) platform, built with Go and the Gin framework following Clean Architecture principles.

## 🚀 Features

### Core Features
- **User Authentication**: Registration, login, logout with secure JWT tokens
- **Email/Phone Verification**: Robust OTP-based verification system with expiration
- **Password Management**: Advanced password policies with strength validation
- **Profile Management**: Complete user profile updates with avatar uploads
- **Account Management**: Secure email/phone changes with OTP verification
- **Company Management**: Create, read, and manage company profiles with logo uploads

### Security & Infrastructure
- **Enhanced Security**: AES-GCM encryption, bcrypt hashing (cost 12), secure cookies
- **JWT Token Management**: Token blacklisting and revocation system
- **Input Validation**: Comprehensive validation middleware with structured errors
- **File Upload**: Cloudinary integration with security checks
- **Database Optimization**: MongoDB indexes for optimal performance
- **API Documentation**: Complete Swagger/OpenAPI documentation

### Developer Experience
- **Structured Error Handling**: Consistent error responses across all endpoints
- **General Response Helpers**: Standardized success/error response formats
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Security Best Practices**: No exposed secrets, proper error handling

## 🏗️ Architecture

This service follows Clean Architecture principles with comprehensive error handling and security:

- **Domain**: Core business entities, repository interfaces, and structured error types
- **Usecase**: Business logic with validation, security checks, and error handling
- **Delivery**: HTTP handlers with structured responses and middleware integration
- **Infrastructure**: External services (database, JWT, email, validation, etc.)
- **Repository**: Data access layer with optimized queries and proper error handling
- **Response**: Centralized response formatting with success/error helpers

## 📡 API Endpoints

### Authentication
- `POST /auth/users/register` - Register new user with avatar upload
- `POST /auth/users/login` - User login with structured responses
- `POST /auth/users/change-password-otp` - Change password with OTP validation
- `GET /auth/users/forgot-password/send-otp` - Send OTP for password reset

### Verification
- `GET /verification/users/send-otp` - Send verification OTP
- `POST /verification/users/verify-otp` - Verify OTP with structured responses

### Protected User Routes (requires JWT)
- `GET /api/users/me` - Get current user profile information
- `GET /api/users/onboard` - Mark user as onboarded
- `POST /api/users/update` - Update user profile with validation
- `POST /api/users/logout` - User logout with token blacklisting
- `POST /api/users/change-email` - Change email with OTP verification
- `GET /api/users/change-email/send-otp` - Send OTP for email change
- `POST /api/users/change-phone` - Change phone with OTP verification  
- `GET /api/users/change-phone/send-otp` - Send OTP for phone change
- `POST /api/users/change-password-old` - Change password with old password validation

### Company Management (requires JWT)
- `GET /api/companies/all` - Get all user companies with pagination and search
- `POST /api/companies/create` - Create new company with logo upload
- `GET /api/companies/:id` - Get company details by ID

### Documentation & Health
- `GET /swagger/*any` - Complete Swagger UI documentation
- `GET /health` - Health check endpoint

## 🛠️ Technology Stack

### Core Framework
- **Framework**: Gin (Go web framework) with middleware support
- **Database**: MongoDB with official Go driver and optimized indexes
- **Authentication**: JWT tokens with blacklisting and secure cookie handling

### Security & Validation
- **Encryption**: AES-GCM for sensitive data encryption
- **Password Hashing**: bcrypt with cost factor 12
- **Input Validation**: Comprehensive validation middleware
- **Security Headers**: Secure cookies, CORS configuration

### External Services
- **File Storage**: Cloudinary integration with error handling
- **Email Service**: SMTP with Gomail for OTP delivery
- **Logging**: Uber Zap for structured logging

### Documentation & Tools
- **API Documentation**: Complete Swagger/OpenAPI specification
- **Error Handling**: Structured error responses with proper HTTP status codes
- **Response Formatting**: Centralized response helpers for consistency

## 📋 Prerequisites

- **Go**: Version 1.21 or higher
- **MongoDB**: Version 4.2 or higher (for proper index support)
- **Cloudinary Account**: For file upload functionality
- **SMTP Server**: For email/OTP delivery (Gmail, SendGrid, etc.)

## ⚙️ Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
PORT=8080

# Database Configuration
MONGO_URI=mongodb://localhost:27017
DB_NAME=byow-user-service

# JWT Configuration
JWT_SECRET=your_secure_jwt_secret_key_here
JWT_EXPIRE=3600

# Email Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USER=your_email@gmail.com
EMAIL_PASS=your_app_password

# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret

# Encryption Configuration
DECRYPT_KEY=your_32_character_encryption_key

# CORS Configuration (optional)
CORS_ORIGINS=http://localhost:3000,https://yourdomain.com
```

### Security Notes:
- Use strong, randomly generated keys for `JWT_SECRET` and `DECRYPT_KEY`
- For Gmail, use App Passwords instead of regular passwords
- Never commit `.env` files to version control

## 🚀 Installation & Setup

### 1. Clone and Setup
```bash
# Clone the repository
git clone https://github.com/buildyow/byow-user-service.git
cd byow-user-service

# Install dependencies
go mod download
```

### 2. Environment Configuration
```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configurations
nano .env
```

### 3. Database Setup
```bash
# Start MongoDB (if using Docker)
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or use existing MongoDB instance
# The application will automatically create indexes on startup
```

### 4. Run & Testing the Application
```bash
# Development mode
go run cmd/main.go

# On Command
go test -v -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

# HTML Output
go tool cover -html=coverage.out

# Production build
go build -o byow-user-service cmd/main.go
./byow-user-service
```

### 5. Verify Installation
```bash
# Check health endpoint
curl http://localhost:8080/health

# Access Swagger documentation
open http://localhost:8080/swagger/index.html
```

The server will start on the port specified in your environment variables (default: 8080).

## 📚 API Documentation

### Swagger Documentation
Once the server is running, access the complete API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Response Formats

#### Success Response
```json
{
  "status": "SUCCESS",
  "code": 200,
  "data": {
    "message": "Operation completed successfully",
    "data": { /* response data */ }
  }
}
```

#### Error Response
```json
{
  "status": "ERROR",
  "code": 400,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Password must contain at least one uppercase letter"
  }
}
```

#### Pagination Response
```json
{
  "status": "SUCCESS",
  "code": 200,
  "data": {
    "message": "Companies retrieved successfully",
    "data": [ /* array of items */ ],
    "row_count": 25
  }
}
```

## 📁 Project Structure

```
├── cmd/
│   └── main.go                    # Application entry point
├── constants/
│   └── constants.go               # Application constants and success messages
├── delivery/
│   └── http/
│       ├── user_handler.go        # User HTTP handlers
│       └── company_handler.go     # Company HTTP handlers
├── docs/                          # Swagger documentation files
├── domain/
│   ├── entity/
│   │   ├── user.go               # User entity
│   │   └── company.go            # Company entity
│   ├── repository/
│   │   ├── user_repository.go    # User repository interface
│   │   └── company_repository.go # Company repository interface
│   └── errors/
│       └── errors.go             # Structured error definitions
├── dto/                          # Data transfer objects
├── infrastructure/
│   ├── cors/                     # CORS configuration
│   ├── db/
│   │   └── indexes.go           # Database indexes management
│   ├── jwt/
│   │   ├── middleware.go        # JWT middleware
│   │   └── blacklist.go         # Token blacklisting system
│   ├── logger/                  # Logging configuration
│   ├── mailer/                  # Email service
│   └── validation/              # Input validation middleware
├── lib/
│   └── cloudinary.go            # Cloudinary integration
├── repository/
│   ├── user_mongo.go            # User MongoDB implementation
│   └── company_mongo.go         # Company MongoDB implementation
├── response/
│   └── response.go              # HTTP response utilities and helpers
├── routes/
│   └── routes.go                # Route definitions
├── usecase/
│   ├── user_usecase.go          # User business logic
│   └── company_usecase.go       # Company business logic
├── utils/
│   └── crypto.go                # Encryption/decryption utilities
└── scripts/                     # Database maintenance scripts
```

## 🔐 Security Features

### Password Security
- **Strong Validation**: Minimum 8 characters, uppercase, lowercase, and numbers required
- **Bcrypt Hashing**: Cost factor 12 for enhanced security
- **Password Change**: Secure flows with OTP or old password verification

### Data Protection
- **AES-GCM Encryption**: Secure encryption for sensitive data like OTP
- **Secure Cookies**: HttpOnly, Secure flags for JWT tokens
- **JWT Token Management**: Token blacklisting and revocation system
- **Input Sanitization**: Comprehensive validation middleware

### API Security
- **CORS Configuration**: Configurable allowed origins
- **Rate Limiting**: (Recommended for production)
- **Error Handling**: No sensitive information leaked in error responses
- **Structured Responses**: Consistent error and success formats

## 🧪 Testing

### Manual Testing
```bash
# Test user registration
curl -X POST http://localhost:8080/auth/users/register \
  -F "full_name=John Doe" \
  -F "email=john@example.com" \
  -F "password=SecurePass123" \
  -F "phone_number=628112123123"

# Test user login
curl -X POST http://localhost:8080/auth/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"SecurePass123"}'

# Test company creation (requires authentication)
curl -X POST http://localhost:8080/api/companies/create \
  -H "Cookie: token=your_jwt_token" \
  -F "company_name=My Company" \
  -F "company_email=company@example.com" \
  -F "company_phone=628112999888" \
  -F "company_address=Jakarta, Indonesia"

# Test get companies
curl -X GET "http://localhost:8080/api/companies/all?limit=10&offset=0&keyword=company" \
  -H "Cookie: token=your_jwt_token"
```

### Health Check
```bash
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "OK",
#   "message": "BYOW User Service is healthy", 
#   "version": "1.0.0"
# }
```

## 🚀 Production Deployment

### Environment Considerations
- Set strong, unique `JWT_SECRET` and `DECRYPT_KEY`
- Use production-grade MongoDB setup with authentication
- Configure proper CORS origins for your frontend
- Set up proper SSL/TLS certificates
- Configure email service with proper authentication

### Performance Optimizations
- MongoDB indexes are automatically created for optimal performance
- JWT token blacklisting with TTL for automatic cleanup
- Optimized database queries with proper error handling

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the existing code style
4. Add tests if applicable
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Create a Pull Request

### Code Standards
- Follow Clean Architecture principles
- Use structured error handling with domain/errors
- Implement proper validation for all inputs
- Add comprehensive documentation for new endpoints
- Ensure security best practices are followed

## 📄 License

This project is part of the BYOW (Build Your Own Website) platform.

## 📞 Support

For support and questions:
- Create an issue in this repository
- Check the Swagger documentation for API details
- Review the structured error responses for troubleshooting