# AccessGo Developer Documentation

## Installation

To use AccessGo in your project, first add it to your `go.mod` file:

```
require github.com/axgrid/accessgo latest
```

Then, run:

```
go mod tidy
```

## Usage

Import the library in your Go code:

```go
import "github.com/axgrid/accessgo"
```

## Available Services

### AccessGoService

#### Initialization

Create a new instance of `AccessGoService`:

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
if err != nil {
    // Handle error
}

service := accessgo.NewAccessGoService(db)
```

#### Available Methods

##### User Management

1. `CreateUser(email, password, name string, userType UserType) (*User, error)`
   - Creates a new user.

2. `UpdateUser(userID uint, email, password, name string, userType UserType) (*User, error)`
   - Updates an existing user's information.

3. `DeleteUser(userID uint) error`
   - Deletes a user.

4. `GetUserByEmail(email string) (*User, error)`
   - Retrieves a user by their email address.

5. `GetUserByID(userID uint) (*User, error)`
   - Retrieves a user by their ID.

6. `GetAllUsers() ([]User, error)`
   - Retrieves all users.

7. `AuthenticateUser(email, password string) (*User, error)`
   - Authenticates a user with their email and password.

8. `GetUserSummaryAccessLevels(userID uint) ([]string, error)`
   - Retrieves all access levels of a user, including those inherited from groups.

9. `GetUserAccessLevels(userID uint) ([]string, error)`
   - Retrieves direct access levels of a user (not including those from groups).

##### Group Management

1. `CreateGroup(name string) (*Group, error)`
   - Creates a new group.

2. `UpdateGroup(groupID uint, name string) (*Group, error)`
   - Updates an existing group's information.

3. `DeleteGroup(groupID uint) error`
   - Deletes a group.

4. `GetGroupByID(groupID uint) (*Group, error)`
   - Retrieves a group by its ID.

5. `GetAllGroups() ([]Group, error)`
   - Retrieves all groups.

##### Access Management

1. `CreateAccess(name, description string) (*Access, error)`
   - Creates a new access right.

2. `UpdateAccess(accessID uint, name, description string) (*Access, error)`
   - Updates an existing access right.

3. `DeleteAccess(accessID uint) error`
   - Deletes an access right.

4. `GetAccessByName(name string) (*Access, error)`
   - Retrieves an access right by its name.

5. `ListAccesses() ([]Access, error)`
   - Lists all access rights.

6. `GetGroupAccessLevels(groupID uint) ([]string, error)`
   - Retrieves all access levels assigned to a group.

##### User-Group Association

1. `AssignUserToGroup(userID, groupID uint) error`
   - Assigns a user to a group.

2. `ExcludeUserFromGroup(userID, groupID uint) error`
   - Removes a user from a group.

3. `SetUserGroups(userID uint, groupIDs ...uint) error`
   - Sets the exact list of groups for a user.

4. `GetUserGroups(userID uint) ([]Group, error)`
   - Retrieves the list of groups a user belongs to.

5. `GetGroupUsers(groupID uint) ([]User, error)`
   - Retrieves the list of users in a group.

##### Access Level Management

1. `AddUserAccessLevel(userID uint, accessName string) error`
   - Adds an access level to a user.

2. `RemoveUserAccessLevel(userID uint, accessName string) error`
   - Removes an access level from a user.

3. `AddGroupAccessLevel(groupID uint, accessName string) error`
   - Adds an access level to a group.

4. `RemoveGroupAccessLevel(groupID uint, accessName string) error`
   - Removes an access level from a group.

5. `CheckUserAccess(userID uint, accessName string) (bool, error)`
   - Checks if a user has a specific access level.

##### System Setup

1. `SetupDefaultPermissions() error`
   - Sets up default permissions in the system.

2. `CreateDefaultAdminUser(email, password, name string) error`
   - Creates a default admin user with full permissions.

### SessionService

The SessionService provides functionality for managing user sessions.

#### Initialization

Create a new instance of `SessionService`:

```go
ctx := context.Background()
sessionService := accessgo.NewSessionService(ctx)
```

#### Methods

1. `CreateSession(userID int, longTerm bool) (string, error)`
   - Creates a new session for a user.
   - Parameters:
      - `userID`: The ID of the user.
      - `longTerm`: Whether this is a long-term session (30 days) or short-term (24 hours).
   - Returns the session ID and an error if any.

2. `GetSession(sessionID string) (Session, error)`
   - Retrieves a session by its ID.
   - Returns the Session struct and an error if the session is not found or has expired.

3. `DeleteSession(sessionID string)`
   - Deletes a session.

4. `ExtendSession(sessionID string) error`
   - Extends the expiration time of a session.
   - For short-term sessions, extends by 24 hours. For long-term sessions, extends by 30 days.

5. `Stop()`
   - Stops the session cleanup routine.

#### Session Cleanup

The SessionService automatically runs a cleanup routine every hour to remove expired sessions. This routine runs in the background and doesn't require manual intervention.

## Data Types

- `User`: Represents a user in the system.
- `Group`: Represents a group of users.
- `Access`: Represents an access right.
- `AccessLevel`: Represents an access level for a user or group.
- `UserType`: Enum representing user types (admin, employee, user, blocked).
- `Session`: Represents a user session.

### User

```go
type User struct {
	gorm.Model
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Name      string
	UserType  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Accesses  []AccessLevel `gorm:"foreignKey:UserID"`
	Groups    []Group       `gorm:"many2many:user_groups;"`
}
```

### Group

```go
type Group struct {
	gorm.Model
	Name     string        `gorm:"unique;not null"`
	Accesses []AccessLevel `gorm:"foreignKey:GroupID"`
	Users    []User        `gorm:"many2many:user_groups;"`
}
```

### Access

```go
type Access struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
}
```

### AccessLevel

```go
type AccessLevel struct {
	gorm.Model
	AccessID uint
	UserID   *uint
	GroupID  *uint
	Access   Access `gorm:"foreignKey:AccessID"`
}
```

### UserType

```go
type UserType string

const (
	UserTypeAdmin    UserType = "admin"
	UserTypeEmployee UserType = "employee"
	UserTypeUser     UserType = "user"
	UserTypeBlocked  UserType = "blocked"
)
```

### Session

```go
type Session struct {
	ID         string
	UserID     int
	CreatedAt  time.Time
	ExpiresAt  time.Time
	IsLongTerm bool
}
```

## Error Handling

Most methods return an error as the last return value. Always check for errors and handle them appropriately in your application.

## Best Practices

1. Always use password hashing (already implemented in the service) when dealing with user passwords.
2. Regularly update and review access rights and user permissions.
3. Use transactions when performing multiple database operations that need to be atomic.
4. Implement proper logging and monitoring in your application to track access and changes to permissions.
5. Manage session lifecycle carefully, creating sessions when users log in and deleting them when they log out.
6. Consider the security implications of long-term sessions and use them judiciously.

## Examples

### Using AccessGoService

```go
service := accessgo.NewAccessGoService(db)

// Create a new user
user, err := service.CreateUser("john@example.com", "password123", "John Doe", accessgo.UserTypeEmployee)
if err != nil {
    // Handle error
}

// Create a new group
group, err := service.CreateGroup("Developers")
if err != nil {
    // Handle error
}

// Assign user to group
err = service.AssignUserToGroup(user.ID, group.ID)
if err != nil {
    // Handle error
}

// Add access level to user
err = service.AddUserAccessLevel(user.ID, "user:read")
if err != nil {
    // Handle error
}

// Check user access
hasAccess, err := service.CheckUserAccess(user.ID, "user:read")
if err != nil {
    // Handle error
}
if hasAccess {
    fmt.Println("User has 'user:read' access")
}

// Get all groups
groups, err := service.GetAllGroups()
if err != nil {
    // Handle error
}
fmt.Println("All groups:", groups)

// Get group access levels
groupAccessLevels, err := service.GetGroupAccessLevels(group.ID)
if err != nil {
    // Handle error
}
fmt.Println("Group access levels:", groupAccessLevels)

// Get user summary access levels
userSummaryAccessLevels, err := service.GetUserSummaryAccessLevels(user.ID)
if err != nil {
    // Handle error
}
fmt.Println("User summary access levels:", userSummaryAccessLevels)
```

### Using SessionService

```go
ctx := context.Background()
sessionService := accessgo.NewSessionService(ctx)

// Create a new session
sessionID, err := sessionService.CreateSession(userID, false)
if err != nil {
    // Handle error
}

// Get session information
session, err := sessionService.GetSession(sessionID)
if err != nil {
    // Handle error (session not found or expired)
}

// Extend session
err = sessionService.ExtendSession(sessionID)
if err != nil {
    // Handle error
}

// Delete session
sessionService.DeleteSession(sessionID)

// Don't forget to stop the session service when your application is shutting down
sessionService.Stop()
```

Remember to properly manage the lifecycle of the SessionService in your application. Create it when your application starts and stop it when your application is shutting down.

This documentation provides an overview of the AccessGo library's capabilities. For more detailed information about each method and its parameters, refer to the source code and comments in the `service.go` and `session.go` files.