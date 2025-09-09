# Daily Coffee API / API Project

A backend project built with **Go** using modern tools, libraries and **Clean Architecture** for building robust, scalable, and secure APIs.

---

## 🚀 Features
- User authentication with **JWT** tokens (access + refresh)
- PostgreSQL database integration
- File/image upload with **Cloudinary**
- Structured logging with **Logrus**
- Input validation with **Validator**
- Configuration management with **Viper**
- Middleware support with **Midtrants**
- Clean architecture ready
- Ready for REST API or microservice integration

---

## 🛠 Tools & Libraries

| Tool / Library | Description | Link |
|----------------|-------------|------|
| **Fiber** | Fast, Express-inspired web framework for Go | [https://gofiber.io](https://gofiber.io) |
| **PostgreSQL** | Open-source relational database | [https://www.postgresql.org](https://www.postgresql.org) |
| **GORM** | ORM library for Go | [https://gorm.io](https://gorm.io) |
| **Logrus** | Structured logger for Go | [https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) |
| **Viper** | Go configuration management | [https://github.com/spf13/viper](https://github.com/spf13/viper) |
| **Midtrants** | Middleware chaining for Fiber | [https://github.com/gofiber/fiber#middleware](https://github.com/gofiber/fiber#middleware) |
| **JWT** | JSON Web Tokens for authentication | [https://github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt) |
| **Validator** | Struct and field validation for Go | [https://github.com/go-playground/validator](https://github.com/go-playground/validator) |
| **Cloudinary Go SDK** | Cloud-based media storage and CDN for images/videos | [https://github.com/cloudinary/cloudinary-go](https://github.com/cloudinary/cloudinary-go) |

---

## 📦 Project Structure

.
cmd/                  # Entry point for the application
internal/
├─ config/            # Configuration (Viper, Gorm, Logger, etc)
├─ delivery/
│  ├─ http/           # HTTP handlers / controllers
│  │  ├─ auth_controller.go
│  │  ├─ category_controller.go
│  │  ├─ customer_controller.go
│  │  ├─ order_controller.go
│  │  ├─ product_controller.go
│  │  ├─ user_controller.go
│  │  └─ route/       # API route definitions
├─ middleware/        # Custom middlewares (e.g., auth)
├─ entity/            # Core entities / domain models
├─ migration/         # Database migrations
├─ model/             # Database models
│  └─ converter/      # Model converters
├─ repository/        # Database repository layer
├─ service/           # Business logic / service layer
├─ usecase/           # Application use cases
└─ utils/             # Utilities (e.g., validator, cloudinary helpers)

---

## ⚡ Installation

1. Clone the repository
```bash
git clone https://github.com/ojihalaw/ojihalaw-daily-coffee-api.git
cd project-name
```

2. Install dependencies
```bash
go mod tidy
```

3. Setup environment variables
```bash
# Copy the example environment file
cp .env_example .env

# Edit the .env file 
```

4. Run the project
```bash
go run cmd/web/main.go
```

## 🔑 Authentication

- Login with **email/password**
- Returns an `AuthResponse` containing:
  - **User info**
  - **Access token** (short-lived)
  - **Refresh token** (long-lived)
- **Access token** is used for API requests
- **Refresh token** is used to generate a new access token

---

## 📄 Logging

- All logs use **Logrus**
- Supports log levels: `debug`, `info`, `warn`, `error`, `fatal`

---

## 🔧 Configuration

- **Viper** handles configuration from:
  - `.env` file
  - Environment variables
  - Default values

---

## 📂 File Upload / Media

- **Cloudinary** handles file/image uploads
- Integrates with API endpoints for user-uploaded images
- Automatically stores files in cloud and returns CDN URLs

---

## 📚 References

- [Fiber Documentation](https://gofiber.io)  
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)  
- [GORM Documentation](https://gorm.io/docs/)  
- [Logrus Documentation](https://github.com/sirupsen/logrus)  
- [Viper Documentation](https://github.com/spf13/viper)  
- [JWT Go Documentation](https://github.com/golang-jwt/jwt)  
- [Validator Go Documentation](https://github.com/go-playground/validator)  
- [Cloudinary Go SDK](https://github.com/cloudinary/cloudinary-go)  

---

## 💡 Notes

- Ensure **PostgreSQL** is running and `.env` is configured correctly
- Use migrations for database schema updates
- Follow clean architecture for maintainability
- Use **Validator** for all input data before processing
- Use **Cloudinary SDK** for media uploads to avoid storing files locally


