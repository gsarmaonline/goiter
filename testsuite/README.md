# Goiter CLI Client

A command-line client for accessing the Goiter server API with Google OAuth authentication.

## Features

- 🔐 Google OAuth authentication  
- 👤 User profile management
- 📋 Project listing and management
- 🏢 Account information retrieval
- 🏓 Server connectivity testing

## Installation

1. Make sure you have Go installed
2. Navigate to the client directory:
   ```bash
   cd client
   ```

## Usage

### Basic Commands

```bash
# Test server connectivity
go run client.go ping

# Login with Google OAuth  
go run client.go login

# Get current user information
go run client.go user

# List user's projects
go run client.go projects  

# Get account information
go run client.go account
```

### Environment Variables

Set the following environment variable if your server is running on a different URL:

```bash
export GOITER_BASE_URL="http://localhost:8080"  # Default
# Or for production:
export GOITER_BASE_URL="https://your-goiter-server.com"
```

## Authentication Process

The authentication process requires manual steps due to the nature of OAuth with CLI applications:

1. Run `go run client.go login`
2. Open your browser and visit the provided Google OAuth URL
3. Complete the Google sign-in process
4. After successful login, you'll be redirected to the frontend with a `token` query parameter in the URL.
5. Copy the value of the `token` parameter.
6. Paste the token value in the terminal

### Finding the JWT Token

**Chrome/Edge:**
1. After logging in, look at the URL in the address bar.
2. It should look something like `http://localhost:3000/?token=ey...`
3. Copy the entire string after `token=`

**Firefox:**
1. After logging in, look at the URL in the address bar.
2. It should look something like `http://localhost:3000/?token=ey...`
3. Copy the entire string after `token=`

**Safari:**
1. After logging in, look at the URL in the address bar.
2. It should look something like `http://localhost:3000/?token=ey...`
3. Copy the entire string after `token=`

## Examples

### Test Server Connection
```bash
$ go run client.go ping
Server is running!
```

### Login and Get User Info
```bash
$ go run client.go user
🔐 Goiter Client Login
======================
Please follow these steps to authenticate:

1. Open your browser and visit: http://localhost:8080/auth/google
2. Complete the Google OAuth flow
3. After successful login, you'll be redirected to the frontend with a `token` query parameter in the URL.
4. Copy the value of the `token` parameter.
5. Paste the token value here: your_jwt_token_here

🔄 Testing authentication...
✅ Login successful! Welcome, John Doe (john@example.com)
User: John Doe (john@example.com)
```

### List Projects
```bash
$ go run client.go projects
🔐 Goiter Client Login
======================
[... authentication process ...]
✅ Login successful! Welcome, John Doe (john@example.com)
Projects (2):
- My First Project: A sample project for testing
- Work Dashboard: Internal company dashboard
```

## Troubleshooting

### Authentication Issues
- Make sure the Goiter server is running and accessible
- Verify you copied the entire JWT token value
- Ensure the server's Google OAuth is properly configured

### Connection Issues  
- Verify the server URL with `go run client.go ping`
- Check that the GOITER_BASE_URL environment variable is set correctly
- Ensure there are no firewall issues blocking the connection

### Browser Issues
- If the browser doesn't open automatically, manually visit the provided URL
- Try using an incognito/private browsing window

## API Coverage

The client currently supports these Goiter API endpoints:

- `GET /ping` - Server health check
- `GET /me` - Current user information  
- `GET /projects` - List user's projects
- `GET /account` - Get account information
- `POST /logout` - Logout (clears session)

## Contributing

To add new API endpoints:

1. Add the corresponding struct types if needed
2. Implement a new method in the `GoiterClient` struct
3. Add a new command case in the `main()` function
4. Update this README with the new functionality

## Notes

- This client uses JWT-based authentication.
- The client doesn't persist tokens between runs - you'll need to re-authenticate each time
- For production use, consider implementing token persistence to a local file 