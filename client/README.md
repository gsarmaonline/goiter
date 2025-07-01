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
4. After successful login, you'll be redirected to the frontend
5. Open your browser's Developer Tools (F12)
6. Navigate to Application/Storage → Cookies
7. Find the 'session' cookie and copy its value
8. Paste the session cookie value in the terminal

### Finding the Session Cookie

**Chrome/Edge:**
1. Press F12 to open Developer Tools
2. Go to "Application" tab
3. In the left sidebar, expand "Storage" → "Cookies"
4. Click on your domain
5. Find the row with "Name" = "session"
6. Copy the "Value" column

**Firefox:**
1. Press F12 to open Developer Tools  
2. Go to "Storage" tab
3. Expand "Cookies" in the left sidebar
4. Click on your domain
5. Find the "session" cookie and copy its value

**Safari:**
1. Press Cmd+Option+I to open Developer Tools
2. Go to "Storage" tab
3. Select "Cookies" and your domain
4. Find the "session" cookie and copy its value

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
3. After successful login, you'll be redirected to the frontend
4. Open your browser's developer tools (F12)
5. Go to Application/Storage -> Cookies
6. Find the 'session' cookie and copy its value
7. Paste the session cookie value here: your_session_cookie_here

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
- Verify you copied the entire session cookie value
- Check that cookies are enabled in your browser
- Ensure the server's Google OAuth is properly configured

### Connection Issues  
- Verify the server URL with `go run client.go ping`
- Check that the GOITER_BASE_URL environment variable is set correctly
- Ensure there are no firewall issues blocking the connection

### Browser Issues
- If the browser doesn't open automatically, manually visit the provided URL
- Try using an incognito/private browsing window
- Clear your browser cookies if you're having authentication issues

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

- This client uses session-based authentication via cookies
- Sessions typically last 7 days (as configured on the server)
- The client doesn't persist sessions between runs - you'll need to re-authenticate each time
- For production use, consider implementing session persistence to a local file 