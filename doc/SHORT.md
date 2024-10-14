# AccessGo Quick Reference

AccessGo is a Go library for managing user authentication, authorization, and session handling.

## Key Components

1. AccessGoService: Manages users, groups, and access rights
2. SessionService: Handles user sessions

## Main Features

### User Management
- Create, update, delete users
- User authentication
- User types: admin, employee, user, blocked

### Group Management
- Create, update, delete groups
- Assign users to groups

### Access Control
- Create, update, delete access rights
- Assign access rights to users and groups
- Check user access

### Session Management
- Create and manage user sessions
- Support for short-term (24h) and long-term (30 days) sessions

## Quick Start

```go
import "github.com/axgrid/accessgo"

// Initialize AccessGoService
db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
accessService := accessgo.NewAccessGoService(db)

// Create a user
user, _ := accessService.CreateUser("user@example.com", "password", "John Doe", accessgo.UserTypeEmployee)

// Create a group
group, _ := accessService.CreateGroup("Developers")

// Assign user to group
accessService.AssignUserToGroup(user.ID, group.ID)

// Add access right
accessService.AddUserAccessLevel(user.ID, "user:read")

// Check access
hasAccess, _ := accessService.CheckUserAccess(user.ID, "user:read")

// Initialize SessionService
ctx := context.Background()
sessionService := accessgo.NewSessionService(ctx)

// Create a session
sessionID, _ := sessionService.CreateSession(int(user.ID), false)

// Get session
session, _ := sessionService.GetSession(sessionID)

// Don't forget to stop SessionService when shutting down
defer sessionService.Stop()
```

## Key Methods

### AccessGoService
- CreateUser, UpdateUser, DeleteUser, GetUserByEmail, GetUserByID, GetAllUsers
- CreateGroup, UpdateGroup, DeleteGroup, GetGroupByID, GetAllGroups
- CreateAccess, UpdateAccess, DeleteAccess, GetAccessByName, ListAccesses
- AssignUserToGroup, ExcludeUserFromGroup, SetUserGroups, GetUserGroups, GetGroupUsers
- AddUserAccessLevel, RemoveUserAccessLevel, AddGroupAccessLevel, RemoveGroupAccessLevel
- CheckUserAccess, GetUserSummaryAccessLevels, GetUserAccessLevels, GetGroupAccessLevels

### SessionService
- CreateSession, GetSession, DeleteSession, ExtendSession

Remember to handle errors and use appropriate database transactions in production code.