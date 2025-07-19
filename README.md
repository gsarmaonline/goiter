# Goiter - Full-Stack Application Boilerplate

A comprehensive boilerplate for building modern SaaS applications with authentication, authorization, project management, and billing capabilities.

## 🚀 Features

### 🔐 Authentication & Authorization

- **JWT-based authentication** with secure token management
- **Role-based access control** with granular permissions
- **User profile management** with automatic profile and account creation

### 💳 Billing & Subscriptions

- **Stripe integration** for payment processing
- **Multiple subscription plans** (Free, Pro, Enterprise)
- **Feature-based plan limitations**
- **Automatic subscription management**
- **Webhook handling** for subscription events

### 📊 Project & Account Management

- **Multi-tenant architecture** with account isolation
- **Project management** with user-based access control
- **Account-level billing** and subscription management
- **User profiles** with customizable settings

### 🛠️ Developer Experience

- **Hot reloading** with Air for rapid development
- **Comprehensive Go client SDK** for API integration
- **Database migrations** and seeding
- **Docker support** for containerized development
- **Makefile** for streamlined development workflow

## 🏗️ Tech Stack

### Backend

- **Go 1.23+** with Gin web framework
- **PostgreSQL/SQLite** database with GORM ORM
- **Stripe** for payment processing
- **JWT** for authentication
- **Air** for hot reloading during development

### Infrastructure

- **Render** for cloud deployment
- **Docker** for containerization
- **Make** for build automation

## 📁 Project Structure

```
goiter/
├── core/                      # Core application logic
│   ├── handlers/              # HTTP request handlers
│   │   ├── account_handler.go # Account management endpoints
│   │   ├── auth_handler.go    # Authentication endpoints
│   │   ├── billing_handler.go # Billing and subscription endpoints
│   │   ├── handler.go         # Base handler utilities
│   │   ├── plan_handler.go    # Plan management endpoints
│   │   ├── profile_handler.go # User profile endpoints
│   │   └── project_handler.go # Project management endpoints
│   ├── middleware/            # Authentication & authorization middleware
│   │   ├── authentication_middleware.go # JWT authentication
│   │   ├── authorisation_middleware.go  # Role-based authorization
│   │   └── middleware.go      # Base middleware utilities
│   ├── models/               # Database models and business logic
│   │   ├── account.go        # Account model
│   │   ├── authorisation.go  # Authorization model
│   │   ├── base_model.go     # Base model structure
│   │   ├── db.go            # Database connection
│   │   ├── plan.go          # Subscription plan model
│   │   ├── profile.go       # User profile model
│   │   ├── project.go       # Project model
│   │   ├── seed.go          # Database seeding
│   │   └── user.go          # User model
│   ├── services/             # External service integrations
│   │   └── stripe_service.go # Stripe payment integration
│   └── server.go             # Server configuration
├── config/                   # Configuration management
│   └── config.go            # Application configuration
├── data/                     # Seed data and migrations
│   └── seed.json            # Initial data seeding
├── testsuite/               # Test suite
│   ├── run/                 # Test runner
│   │   └── run.go          # Test execution
│   ├── account.go          # Account tests
│   ├── profile.go          # Profile tests
│   ├── project.go          # Project tests
│   ├── server.go           # Test server setup
│   ├── testsuite.go        # Test suite utilities
│   ├── user.go             # User tests
│   └── README.md           # Test documentation
├── main.go                  # Application entry point
├── Makefile                 # Development workflow commands
├── render.yaml              # Deployment configuration
├── go.mod                   # Go module definition
├── go.sum                   # Go dependency checksums
└── gorm.db                  # SQLite database file (development)
```

## 🚦 Quick Start

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 13 or higher (or SQLite for development)
- Stripe account (for billing features)

### 1. Clone the Repository

```bash
git clone https://github.com/gsarmaonline/goiter.git
cd goiter
```

### 2. Database Setup

```bash
# Create PostgreSQL database
createdb goiter

# Or using psql
psql -U postgres -c "CREATE DATABASE goiter;"
```

### 3. Environment Configuration

Create a `.env` file in the root directory:

```env
# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=goiter

# Or use SQLite for development (comment out PostgreSQL config above)
# DB_NAME=gorm.db

# Server Configuration
PORT=8080
MODE=dev
GIN_MODE=debug

# JWT Configuration
JWT_SECRET=your_jwt_secret

# Stripe Configuration
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
# Start with hot reloading (recommended for development)
make start-backend

# Or start without hot reloading
make start-backend-no-air
```

### 6. Test the Installation

```bash
# Test server connectivity
curl http://localhost:8080/ping

# Run the test suite
make test
```

## 🔧 Development Workflow

### Available Make Commands

```bash
# Start backend with hot reloading
make start-backend

# Start backend without hot reloading
make start-backend-no-air

# Stop backend server
make stop-backend

# Clean up processes and ports
make clean-air

# Database operations
make db          # Connect to database
make clean       # Reset database

# Testing
make test        # Run test suite
```

### Database Management

```bash
# Connect to database
make db

# Reset database (drops and recreates)
make clean

# The application automatically:
# - Runs migrations on startup
# - Seeds initial data (plans, role permissions)
# - Creates user profiles and accounts
```

## 📚 API Documentation

### Authentication Endpoints

- `POST /login` - User login with credentials
- `POST /register` - User registration
- `GET /me` - Get current user information
- `POST /logout` - Logout current user

### User Management

- `GET /me` - Get current user profile
- `GET /profile` - Get detailed user profile
- `PUT /profile` - Update user profile

### Project Management

- `GET /projects` - List user's projects
- `POST /projects` - Create new project
- `GET /projects/:id` - Get project details
- `PUT /projects/:id` - Update project
- `DELETE /projects/:id` - Delete project

### Account & Billing

- `GET /account` - Get account information
- `PUT /account` - Update account settings
- `GET /plans` - List available subscription plans
- `POST /billing/subscribe` - Create subscription
- `POST /billing/portal` - Access billing portal

### Utility Endpoints

- `GET /ping` - Health check
- `GET /plans` - List available plans

## 🚀 Deployment

### Render Deployment

The project includes a `render.yaml` file for easy deployment to Render:

1. **Push to GitHub**: Ensure your code is in a GitHub repository
2. **Create Render Account**: Sign up at [render.com](https://render.com)
3. **Create New Web Service**: Connect your GitHub repository
4. **Configure Environment Variables**:
   - `JWT_SECRET`
   - `STRIPE_PUBLISHABLE_KEY`, `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`
   - Database credentials (auto-configured by Render)

### Manual Deployment

```bash
# Build the application
go build -o main

# Run in production
./main
```

## 🎯 Subscription Plans

The boilerplate includes a flexible plan system:

### Free Plan

- 1 project limit
- Basic features
- No billing required

### Pro Plan

- 10 project limit
- Advanced features
- $10/month

### Custom Plans

Easily add new plans by modifying `data/seed.json`:

```json
{
  "plans": [
    {
      "name": "Enterprise",
      "price": 50,
      "description": "Enterprise plan with unlimited projects",
      "features": [
        {
          "name": "Projects",
          "limit": -1
        }
      ]
    }
  ]
}
```

## 🔒 Security Features

- **JWT-based authentication** with secure token settings
- **CORS configuration** for cross-origin requests
- **Input validation** and sanitization
- **SQL injection protection** via GORM
- **Rate limiting ready** (middleware available)
- **HTTPS enforcement** in production

## 🧪 Testing

```bash
# Run the test suite
make test

# Run tests manually
go run testsuite/run/run.go

# Run individual test files
go run testsuite/user.go
go run testsuite/project.go
go run testsuite/account.go
```

## 📈 Monitoring & Logging

The application includes structured logging and is ready for monitoring integration:

- **Structured JSON logging** for production
- **Request/response logging** middleware
- **Error tracking** with detailed stack traces
- **Performance metrics** ready for integration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new features
- Update documentation for API changes
- Use meaningful commit messages
- Ensure all tests pass before submitting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: Report bugs and request features via [GitHub Issues](https://github.com/gsarmaonline/goiter/issues)
- **Discussions**: Join the community in [GitHub Discussions](https://github.com/gsarmaonline/goiter/discussions)
- **Documentation**: Visit our [Wiki](https://github.com/gsarmaonline/goiter/wiki) for detailed guides

## 🙏 Acknowledgments

- [Gin Framework](https://gin-gonic.com/) for the excellent web framework
- [GORM](https://gorm.io/) for the powerful ORM
- [Stripe](https://stripe.com/) for payment processing
- [JWT](https://jwt.io/) for authentication
- [Render](https://render.com/) for deployment platform

---

**Happy coding! 🚀**

> Goiter provides everything you need to build a modern SaaS application. Focus on your unique business logic while we handle the boilerplate.
