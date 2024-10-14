# UserGoService

UserGoService - это Go-библиотека для управления пользователями, группами и правами доступа. Она предоставляет функциональность для аутентификации, авторизации и управления правами доступа в приложениях.

## Основные возможности

- Управление пользователями (создание, обновление, удаление, аутентификация)
- Управление группами (создание, обновление, удаление)
- Управление правами доступа (создание, назначение пользователям и группам)
- Проверка прав доступа пользователей
- Назначение пользователей в группы

## Установка

Для использования UserGoService в вашем проекте, выполните следующую команду:

```
go get github.com/axgrid/accessgo
```

## Использование

Вот пример базового использования UserGoService:

```go
import (
    "github.com/axgrid/accessgo"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

// Инициализация базы данных
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
if err != nil {
    panic("failed to connect database")
}

// Создание экземпляра сервиса
service := accessgo.NewUserGoService(db)

// Создание пользователя
user, err := service.CreateUser("user@example.com", "password", "John Doe", accessgo.UserTypeUser)
if err != nil {
    panic(err)
}

// Аутентификация пользователя
authenticatedUser, err := service.AuthenticateUser("user@example.com", "password")
if err != nil {
    panic(err)
}

// Создание группы
group, err := service.CreateGroup("Administrators")
if err != nil {
    panic(err)
}

// Назначение пользователя в группу
err = service.AssignUserToGroup(user.ID, group.ID)
if err != nil {
    panic(err)
}

// Создание права доступа
access, err := service.CreateAccess("user:read", "Право на чтение информации о пользователях")
if err != nil {
    panic(err)
}

// Назначение права доступа группе
err = service.AddGroupAccessLevel(group.ID, "user:read")
if err != nil {
    panic(err)
}

// Проверка права доступа у пользователя
hasAccess, err := service.CheckUserAccess(user.ID, "user:read")
if err != nil {
    panic(err)
}
fmt.Printf("User has access: %v\n", hasAccess)
```

## Основные методы

- `CreateUser(email, password, name string, userType UserType) (*User, error)`
- `AuthenticateUser(email, password string) (*User, error)`
- `UpdateUser(userID uint, email, password, name string, userType UserType) (*User, error)`
- `DeleteUser(userID uint) error`
- `CreateGroup(name string) (*Group, error)`
- `UpdateGroup(groupID uint, name string) (*Group, error)`
- `DeleteGroup(groupID uint) error`
- `AssignUserToGroup(userID, groupID uint) error`
- `ExcludeUserFromGroup(userID, groupID uint) error`
- `CreateAccess(name, description string) (*Access, error)`
- `AddUserAccessLevel(userID uint, accessName string) error`
- `AddGroupAccessLevel(groupID uint, accessName string) error`
- `CheckUserAccess(userID uint, accessName string) (bool, error)`
- `GetUserAccessLevels(userID uint) ([]string, error)`
- `SetupDefaultPermissions() error`

## Тестирование

Для запуска тестов выполните следующую команду:

```
go test github.com/axgrid/accessgo
```

## Вклад в проект

Мы приветствуем вклад в развитие проекта! Пожалуйста, создавайте issue или pull request на GitHub.

## Лицензия

Этот проект лицензирован под [MIT License](https://opensource.org/licenses/MIT).