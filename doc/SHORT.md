# AccessGo: Brief Methods Overview

## Service Initialization

### AccessGoService
```go
import (
    "github.com/axgrid/accessgo"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite" // or any other database driver
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
if err != nil {
    // Handle error
}

accessService := accessgo.NewAccessGoService(db)
```

### SessionService
```go
import (
    "context"
    "github.com/axgrid/accessgo"
)

ctx := context.Background()
sessionService := accessgo.NewSessionService(ctx)

// Don't forget to stop the service when your application shuts down
defer sessionService.Stop()
```

## AccessGoService Methods

### User Management
- CreateUser(email, password, name string, userType UserType) (*User, error)
- UpdateUser(userID uint, email, password, name string, userType UserType) (*User, error)
- DeleteUser(userID uint) error
- GetUserByEmail(email string) (*User, error)
- GetUserByID(userID uint) (*User, error)
- GetAllUsers() ([]User, error)
- AuthenticateUser(email, password string) (*User, error)

### Group Management
- CreateGroup(name string) (*Group, error)
- UpdateGroup(groupID uint, name string) (*Group, error)
- DeleteGroup(groupID uint) error
- GetGroupByID(groupID uint) (*Group, error)
- GetAllGroups() ([]Group, error)

### Access Management
- CreateAccess(name, description string) (*Access, error)
- UpdateAccess(accessID uint, name, description string) (*Access, error)
- DeleteAccess(accessID uint) error
- GetAccessByName(name string) (*Access, error)
- ListAccesses() ([]Access, error)

### User-Group Association
- AssignUserToGroup(userID, groupID uint) error
- ExcludeUserFromGroup(userID, groupID uint) error
- SetUserGroups(userID uint, groupIDs ...uint) error
- GetUserGroups(userID uint) ([]Group, error)
- GetGroupUsers(groupID uint) ([]User, error)

### Access Level Management
- AddUserAccessLevel(userID uint, accessName string) error
- RemoveUserAccessLevel(userID uint, accessName string) error
- AddGroupAccessLevel(groupID uint, accessName string) error
- RemoveGroupAccessLevel(groupID uint, accessName string) error
- CheckUserAccess(userID uint, accessName string) (bool, error)
- GetUserSummaryAccessLevels(userID uint) ([]string, error)
- GetUserAccessLevels(userID uint) ([]string, error)
- GetGroupAccessLevels(groupID uint) ([]string, error)

### System Setup
- SetupDefaultPermissions() error
- CreateDefaultAdminUser(email, password, name string) error

## SessionService Methods

- CreateSession(userID int, longTerm bool) (string, error)
- GetSession(sessionID string) (Session, error)
- DeleteSession(sessionID string)
- ExtendSession(sessionID string) error
- Stop()

Note: All methods may return an error. Always check and handle errors in your application code.

## Usage Example

```go
// Initialize services
db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
accessService := accessgo.NewAccessGoService(db)
ctx := context.Background()
sessionService := accessgo.NewSessionService(ctx)

// Create a user
user, _ := accessService.CreateUser("user@example.com", "password", "John Doe", accessgo.UserTypeEmployee)

// Create a group
group, _ := accessService.CreateGroup("Developers")

// Assign user to group
accessService.AssignUserToGroup(user.ID, group.ID)

// Add access right
accessService.AddUserAccessLevel(user.ID, "user:read")

// Create a session
sessionID, _ := sessionService.CreateSession(int(user.ID), false)

// Clean up
defer sessionService.Stop()
```