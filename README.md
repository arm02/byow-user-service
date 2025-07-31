# BYOW User Service

A comprehensive user management microservice for the Build Your Own Website (BYOW) platform, built with Go and the Gin framework.

## Features

- **User Authentication**: Registration, login, logout with JWT tokens
- **Email/Phone Verification**: OTP-based verification system
- **Password Management**: Change password with OTP or old password
- **Profile Management**: Update user information including avatar uploads
- **Account Management**: Change email/phone with OTP verification
- **File Upload**: Cloudinary integration for avatar uploads
- **API Documentation**: Swagger/OpenAPI documentation
- **Security**: JWT middleware, password hashing, CORS support

## Architecture

This service follows Clean Architecture principles with the following layers:

- **Domain**: Core business entities and repository interfaces
- **Usecase**: Business logic and application services
- **Delivery**: HTTP handlers and API endpoints
- **Infrastructure**: External services (database, JWT, email, etc.)
- **Repository**: Data access layer

## API Endpoints

### Authentication
- `POST /auth/users/register` - Register new user with avatar
- `POST /auth/users/login` - User login
- `POST /auth/users/change-password-otp` - Change password with OTP
- `GET /auth/users/forgot-password/send-otp` - Send OTP for password reset

### Verification
- `GET /verification/users/send-otp` - Send verification OTP
- `POST /verification/users/verify-otp` - Verify OTP

### Protected User Routes (requires JWT)
- `GET /api/users/me` - Get current user info
- `GET /api/users/onboard` - Mark user as onboarded
- `POST /api/users/update` - Update user profile
- `POST /api/users/logout` - User logout
- `POST /api/users/change-email` - Change email with OTP
- `GET /api/users/change-email/send-otp` - Send OTP for email change
- `POST /api/users/change-phone` - Change phone with OTP
- `GET /api/users/change-phone/send-otp` - Send OTP for phone change
- `POST /api/users/change-password-old` - Change password with old password

### Documentation
- `GET /swagger/*any` - Swagger UI documentation

## Technology Stack

- **Framework**: Gin (Go web framework)
- **Database**: MongoDB with official Go driver
- **Authentication**: JWT tokens
- **File Storage**: Cloudinary
- **Email Service**: SMTP with Gomail
- **Logging**: Uber Zap
- **Documentation**: Swagger/OpenAPI
- **CORS**: Gin CORS middleware

## Prerequisites

- Go 1.24.5 or higher
- MongoDB instance
- Cloudinary account (for file uploads)
- SMTP server for email functionality

## Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
DB_NAME=byow_user_service
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE=3600
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USER=your_email@gmail.com
EMAIL_PASS=your_email_password
```

## Installation & Setup

1. Clone the repository:
```bash
git clone https://github.com/buildyow/byow-user-service.git
cd byow-user-service
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (see above)

4. Run the application:
```bash
go run cmd/main.go
```

The server will start on the port specified in your environment variables (default: 8080).

## API Documentation

Once the server is running, access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Project Structure

```
├── cmd/
│   └── main.go              # Application entry point
├── constants/
│   └── constants.go         # Application constants
├── delivery/
│   └── http/
│       └── user_handler.go  # HTTP handlers
├── docs/                    # Swagger documentation files
├── domain/
│   ├── entity/
│   │   └── user.go         # User entity
│   └── repository/
│       └── user_repository.go # Repository interface
├── dto/                     # Data transfer objects
├── infrastructure/
│   ├── cors/               # CORS configuration
│   ├── db/                 # Database connection
│   ├── jwt/                # JWT middleware and utilities
│   ├── logger/             # Logging configuration
│   └── mailer/             # Email service
├── lib/
│   └── cloudinary.go       # Cloudinary integration
├── repository/
│   └── user_mongo.go       # MongoDB implementation
├── response/
│   └── response.go         # HTTP response utilities
├── routes/
│   └── routes.go           # Route definitions
├── usecase/
│   └── user_usecase.go     # Business logic
└── utils/                  # Utility functions
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Commit your changes
6. Push to the branch
7. Create a Pull Request

## License

This project is part of the BYOW platform.