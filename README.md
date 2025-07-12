# Goiter - Full-Stack Application Boilerplate

A comprehensive boilerplate for building modern SaaS applications with authentication, authorization, project management, and billing capabilities.

## 🚀 Features

### 🔐 Authentication & Authorization

- **Google OAuth 2.0** integration for secure user authentication
- **Session-based authentication** with secure cookie management
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
- **PostgreSQL** database with GORM ORM
- **Stripe** for payment processing
- **Google OAuth 2.0** for authentication
- **Air** for hot reloading during development

### Frontend (Coming Soon)

- **React** with modern hooks
- **TypeScript** for type safety
- **Tailwind CSS** for styling
- **Stripe Elements** for payment forms

### Infrastructure

- **Render** for cloud deployment
- **Docker** for containerization
- **Make** for build automation

## 📁 Project Structure

```
goiter/
├── backend/                    # Go backend application
│   ├── core/                  # Core application logic
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── middleware/        # Authentication & authorization middleware
│   │   ├── models/           # Database models and business logic
│   │   ├── services/         # External service integrations
│   │   └── server.go         # Server configuration
│   ├── data/                 # Seed data and migrations
│   ├── main.go               # Application entry point
│   └── tmp/                  # Temporary files (Air hot reload)
├── client/                   # Go client SDK
│   ├── client.go             # Client implementation
│   └── README.md             # Client documentation
├── Makefile                  # Development workflow commands
└── render.yaml              # Deployment configuration
```

## 🚦 Quick Start

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 13 or higher
- Stripe account (for billing features)
- Google OAuth credentials

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

Create a `.env` file in the `backend/` directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=goiter
DB_SSLMODE=disable

# Server Configuration
PORT=8080
FRONTEND_URL=http://localhost:3000

# Google OAuth Configuration
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_CALLBACK_URL=http://localhost:8080/auth/google/callback

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
go run client/client.go ping

# Test authentication flow
go run client/client.go login
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

- `GET /auth/google` - Initiate Google OAuth flow
- `GET /auth/google/callback` - Handle OAuth callback
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
   - `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`
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

Easily add new plans by modifying `backend/data/seed.json`:

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

- **Session-based authentication** with secure cookie settings
- **CORS configuration** for cross-origin requests
- **Input validation** and sanitization
- **SQL injection protection** via GORM
- **Rate limiting ready** (middleware available)
- **HTTPS enforcement** in production

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./core/handlers -v
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
- [Google OAuth](https://developers.google.com/identity/protocols/oauth2) for authentication
- [Render](https://render.com/) for deployment platform

---

**Happy coding! 🚀**

> Goiter provides everything you need to build a modern SaaS application. Focus on your unique business logic while we handle the boilerplate.
