# Goiter Authorization Test Suite

This test suite provides comprehensive testing for authorization permissions in the Goiter application.

## Test Structure

### 1. Basic Functional Tests
- **User Suite**: Tests user authentication and profile management
- **Profile Suite**: Tests user profile operations
- **Account Suite**: Tests account management
- **Project Suite**: Tests basic project CRUD operations

### 2. Authorization Tests
- **Project Permission Tests**: Tests permission levels for project operations
- **Unauthorized Access Tests**: Tests access without proper authentication
- **Resource Permission Tests**: Tests permissions on different resource types
- **Cross-Project Permission Tests**: Tests that users can't access resources across projects
- **Permission Inheritance Tests**: Tests that higher permission levels include lower-level permissions

## Running Tests

### Run All Tests
```bash
make test
```

### Run with Environment Variables
```bash
# Enable verbose output
GOITER_TEST_VERBOSE=true make test

# Disable authorization tests (only run basic functional tests)
GOITER_TEST_AUTH=false make test

# Use custom server URL
GOITER_BASE_URL=http://localhost:8080 make test
```

### Run Specific Test Categories

#### Basic Functional Tests Only
```bash
GOITER_TEST_AUTH=false go run testsuite/run/run.go
```

#### Authorization Tests Only
Create a simple runner script:
```go
// auth_test_runner.go
package main

import (
    "github.com/gsarmaonline/goiter/testsuite"
)

func main() {
    go testsuite.StartServer()
    time.Sleep(2 * time.Second)
    
    client := testsuite.NewAuthTestClient("http://localhost:8090")
    
    // Run only authorization tests
    if err := client.RunProjectPermissionTests(); err != nil {
        log.Fatalf("Project permission tests failed: %v", err)
    }
    
    if err := client.TestResourcePermissions(); err != nil {
        log.Fatalf("Resource permission tests failed: %v", err)
    }
    
    log.Println("Authorization tests completed!")
}
```

## Test Configuration

Environment variables:
- `GOITER_BASE_URL`: Server URL (default: `http://localhost:8090`)
- `GOITER_TEST_AUTH`: Enable authorization tests (default: `true`)
- `GOITER_TEST_VERBOSE`: Enable verbose output (default: `false`)

## Permission Levels Tested

1. **Owner (20)**: Full access to all resources
2. **Admin (19)**: Can manage resources and some project settings
3. **Editor (18)**: Can edit resources
4. **Viewer (17)**: Can only view resources

## Test Coverage

### Project Operations
- ✅ Read project details
- ✅ Update project details
- ✅ Delete project
- ✅ List projects

### Project Member Management
- ✅ Add project members
- ✅ Remove project members
- ✅ View project members

### Resource Access Control
- ✅ Generic resource permissions (`*`)
- ✅ Specific resource permissions (`project`, `project_member`)
- ✅ Action-based permissions (`read`, `create`, `update`, `delete`)

### Security Tests
- ✅ Unauthorized access prevention
- ✅ Cross-project access prevention
- ✅ Permission inheritance validation

## Test Users

The test suite automatically creates test users with the following pattern:
- `owner@example.com` - Project owner
- `user_project_{level}_{action}@example.com` - Users with specific permission levels
- `user_resource_{type}_{action}_{level}@example.com` - Users for resource-specific tests

## Expected Test Results

### Successful Tests
- All permission levels can read resources they have access to
- Higher permission levels can perform actions of lower levels
- Users cannot access resources in projects they're not members of
- Unauthorized requests return 401 status
- Insufficient permissions return 403 status

### Test Failures
Tests will fail if:
- Permission levels don't match expected behavior
- Users can access resources they shouldn't have access to
- Authorization middleware is not properly configured
- Database permissions are not correctly set up

## Debugging Test Failures

1. **Enable verbose output**: Set `GOITER_TEST_VERBOSE=true`
2. **Check server logs**: Look for authentication/authorization errors
3. **Verify database state**: Check that permissions are correctly stored
4. **Run individual tests**: Isolate specific failing test scenarios

## Adding New Tests

### 1. Add Permission Test Scenarios
Add new scenarios to `permission_tests.go`:
```go
var NewResourcePermissionTests = []PermissionTestScenario{
    {
        Name:           "Test name",
        UserLevel:      models.PermissionEditor,
        Action:         "create",
        Method:         "POST",
        Endpoint:       "/api/resource",
        Body:           map[string]string{"field": "value"},
        ExpectedStatus: 200,
        Description:    "Test description",
    },
}
```

### 2. Add Resource-Specific Tests
Add new resource tests to `resource_permission_tests.go`:
```go
var NewResourceTests = []ResourcePermissionTest{
    {
        ResourceType: "new_resource",
        Action:       models.CreateAction,
        MinLevel:     models.PermissionEditor,
        Description:  "Creating new resources",
        ShouldSucceed: map[models.PermissionLevel]bool{
            models.PermissionOwner:  true,
            models.PermissionAdmin:  true,
            models.PermissionEditor: true,
            models.PermissionViewer: false,
        },
    },
}
```

### 3. Update Test Suite
Add your new test function to `testsuite.go`:
```go
if err := authClient.RunNewResourceTests(); err != nil {
    log.Fatalf("❌ New resource tests failed: %v", err)
}
```

## Architecture

The test suite uses:
- **HTTP client testing**: Real HTTP requests to test endpoints
- **JWT authentication**: Proper token-based authentication
- **Test isolation**: Each test creates its own users and projects
- **Cleanup**: Automatic cleanup of test data
- **Comprehensive coverage**: Tests all permission levels and actions

This ensures that your authorization system works correctly in real-world scenarios.