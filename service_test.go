package accessgo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&User{}, &Group{}, &Access{}, &AccessLevel{})
	require.NoError(t, err)

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	user, err := service.CreateUser("test@example.com", "password", "Test User", UserTypeUser)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, string(UserTypeUser), user.UserType)
}

func TestCreateAndAuthenticateUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	// Создаем пользователя
	_, err := service.CreateUser("auth@example.com", "password", "Auth User", UserTypeUser)
	assert.NoError(t, err)

	// Аутентифицируем пользователя
	user, err := service.AuthenticateUser("auth@example.com", "password")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "auth@example.com", user.Email)

	// Проверяем неверный пароль
	_, err = service.AuthenticateUser("auth@example.com", "wrongpassword")
	assert.Error(t, err)
}

func TestCreateAndDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	user, err := service.CreateUser("delete@example.com", "password", "Delete User", UserTypeUser)
	assert.NoError(t, err)

	err = service.DeleteUser(user.ID)
	assert.NoError(t, err)

	// Проверяем, что пользователь удален
	_, err = service.GetUserByEmail("delete@example.com")
	assert.Error(t, err)
}

func TestCreateAndUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	user, err := service.CreateUser("update@example.com", "password", "Update User", UserTypeUser)
	assert.NoError(t, err)

	updatedUser, err := service.UpdateUser(user.ID, "newemail@example.com", "newpassword", "New Name", UserTypeEmployee)
	assert.NoError(t, err)
	assert.Equal(t, "newemail@example.com", updatedUser.Email)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, string(UserTypeEmployee), updatedUser.UserType)
}

func TestCreateGroupAndAssignUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	user, err := service.CreateUser("group@example.com", "password", "Group User", UserTypeUser)
	assert.NoError(t, err)

	group, err := service.CreateGroup("Test Group")
	assert.NoError(t, err)

	err = service.AssignUserToGroup(user.ID, group.ID)
	assert.NoError(t, err)

	userGroups, err := service.GetUserGroups(user.ID)
	assert.NoError(t, err)
	assert.Len(t, userGroups, 1)
	assert.Equal(t, group.ID, userGroups[0].ID)
}

func TestSetupDefaultPermissionsAndCreateAdmin(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	err := service.SetupDefaultPermissions()
	assert.NoError(t, err)

	err = service.CreateDefaultAdminUser("admin@example.com", "adminpass", "Admin User")
	assert.NoError(t, err)

	admin, err := service.GetUserByEmail("admin@example.com")
	assert.NoError(t, err)
	assert.Equal(t, string(UserTypeAdmin), admin.UserType)

	accessLevels, err := service.GetUserAccessLevels(admin.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessLevels)
}

func TestAddAndCheckUserAccess(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	user, err := service.CreateUser("access@example.com", "password", "Access User", UserTypeUser)
	assert.NoError(t, err)

	err = service.SetupDefaultPermissions()
	assert.NoError(t, err)

	err = service.AddUserAccessLevel(user.ID, "user:read")
	assert.NoError(t, err)

	hasAccess, err := service.CheckUserAccess(user.ID, "user:read")
	assert.NoError(t, err)
	assert.True(t, hasAccess)

	hasAccess, err = service.CheckUserAccess(user.ID, "user:write")
	assert.NoError(t, err)
	assert.False(t, hasAccess)
}

func TestAddAndCheckUserGroupAccess(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	user, err := service.CreateUser("access@example.com", "password", "Access User", UserTypeUser)
	assert.NoError(t, err)

	err = service.SetupDefaultPermissions()
	assert.NoError(t, err)

	group, err := service.CreateGroup("Test Group")
	assert.NoError(t, err)

	err = service.AssignUserToGroup(user.ID, group.ID)
	assert.NoError(t, err)

	err = service.AddGroupAccessLevel(user.ID, "user:read")
	assert.NoError(t, err)

	hasAccess, err := service.CheckUserAccess(user.ID, "user:read")
	assert.NoError(t, err)
	assert.True(t, hasAccess)

	hasAccess, err = service.CheckUserAccess(user.ID, "user:write")
	assert.NoError(t, err)
	assert.False(t, hasAccess)
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserGoService(db)

	// Создаем пользователя
	email := "getuser@example.com"
	_, err := service.CreateUser(email, "password", "Get User", UserTypeUser)
	assert.NoError(t, err)

	// Получаем пользователя по email
	user, err := service.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "Get User", user.Name)
	assert.Equal(t, string(UserTypeUser), user.UserType)

	// Пробуем получить несуществующего пользователя
	_, err = service.GetUserByEmail("nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, "пользователь не найден", err.Error())
}
