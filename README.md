# AccessGo

AccessGo - это пакет на Go для управления пользователями, группами и правами доступа в приложениях. Он предоставляет гибкую и расширяемую систему управления доступом, которая может быть легко интегрирована в различные проекты.

## Возможности

- Управление пользователями (создание, обновление, удаление)
- Управление группами пользователей
- Определение и назначение прав доступа
- Проверка прав доступа для пользователей и групп
- Аутентификация пользователей
- Управление сессиями
- Подтверждение email пользователей

## Установка

```bash
go get github.com/yourusername/accessgo
```

## Использование

### Инициализация сервиса

```go
import (
    "github.com/yourusername/accessgo"
    "gorm.io/gorm"
)

db, err := gorm.Open(...)
if err != nil {
    // обработка ошибки
}

service := accessgo.NewAccessGoService(db)
```

### Создание пользователя

```go
user, err := service.CreateUser("user@example.com", "password", "John Doe", accessgo.UserTypeEmployee)
if err != nil {
    // обработка ошибки
}
```

### Аутентификация пользователя

```go
user, err := service.AuthenticateUser("user@example.com", "password")
if err != nil {
    // обработка ошибки
}
```

### Управление правами доступа

```go
err := service.AddUserAccessLevel(userID, "user:read")
if err != nil {
    // обработка ошибки
}

hasAccess, err := service.CheckUserAccess(userID, "user:read")
if err != nil {
    // обработка ошибки
}
```

### Управление сессиями

```go
sessionService := accessgo.NewSessionService(context.Background())

sessionID, err := sessionService.CreateSession(userID, false)
if err != nil {
    // обработка ошибки
}

session, err := sessionService.GetSession(sessionID)
if err != nil {
    // обработка ошибки
}
```

## Основные методы

### Управление пользователями

- `CreateUser(email, password, name string, userType UserType) (*User, error)`: Создает нового пользователя.
- `UpdateUser(userID uint, email, password, name string, userType UserType) (*User, error)`: Обновляет информацию о пользователе.
- `DeleteUser(userID uint) error`: Удаляет пользователя.
- `GetUserByEmail(email string) (*User, error)`: Получает пользователя по email.
- `GetUserByID(userID uint) (*User, error)`: Получает пользователя по ID.
- `GetAllUsers() ([]User, error)`: Получает список всех пользователей.
- `ValidateEmail(token string) error`: Подтверждает email пользователя.

### Управление группами

- `CreateGroup(name string) (*Group, error)`: Создает новую группу.
- `UpdateGroup(groupID uint, name string) (*Group, error)`: Обновляет информацию о группе.
- `DeleteGroup(groupID uint) error`: Удаляет группу.
- `GetGroupByID(groupID uint) (*Group, error)`: Получает группу по ID.
- `GetAllGroups() ([]Group, error)`: Получает список всех групп.
- `AssignUserToGroup(userID, groupID uint) error`: Добавляет пользователя в группу.
- `ExcludeUserFromGroup(userID, groupID uint) error`: Удаляет пользователя из группы.
- `SetUserGroups(userID uint, groupIDs ...uint) error`: Устанавливает точный список групп для пользователя.
- `GetUserGroups(userID uint) ([]Group, error)`: Получает список групп пользователя.
- `GetGroupUsers(groupID uint) ([]User, error)`: Получает список пользователей в группе.

### Управление правами доступа

- `CreateAccess(name, description string) (*Access, error)`: Создает новое право доступа.
- `UpdateAccess(accessID uint, name, description string) (*Access, error)`: Обновляет информацию о праве доступа.
- `DeleteAccess(accessID uint) error`: Удаляет право доступа.
- `GetAccessByName(name string) (*Access, error)`: Получает право доступа по имени.
- `ListAccesses() ([]Access, error)`: Возвращает список всех прав доступа.
- `AddUserAccessLevel(userID uint, accessName string) error`: Добавляет уровень доступа пользователю.
- `RemoveUserAccessLevel(userID uint, accessName string) error`: Удаляет уровень доступа у пользователя.
- `AddGroupAccessLevel(groupID uint, accessName string) error`: Добавляет уровень доступа группе.
- `RemoveGroupAccessLevel(groupID uint, accessName string) error`: Удаляет уровень доступа у группы.
- `CheckUserAccess(userID uint, accessName string) (bool, error)`: Проверяет, имеет ли пользователь указанный уровень доступа.
- `GetUserSummaryAccessLevels(userID uint) ([]string, error)`: Возвращает все уровни доступа пользователя (включая групповые).
- `GetGroupAccessLevels(groupID uint) ([]string, error)`: Возвращает все уровни доступа группы.
- `GetUserAccessLevels(userID uint) ([]string, error)`: Возвращает все прямые уровни доступа пользователя.

### Аутентификация и сессии

- `AuthenticateUser(email, password string) (*User, error)`: Аутентифицирует пользователя по email и паролю.
- `CreateSession(userID int, longTerm bool) (string, error)`: Создает новую сессию для пользователя.
- `GetSession(sessionID string) (Session, error)`: Получает информацию о сессии.
- `DeleteSession(sessionID string)`: Удаляет сессию.
- `ExtendSession(sessionID string) error`: Продлевает сессию.

### Дополнительные методы

- `SetupDefaultPermissions() error`: Создает стандартные права доступа.
- `CreateDefaultAdminUser(email, password, name string) error`: Создает пользователя-администратора с полными правами.

## Структура проекта

- `structs.go`: Определения основных структур данных
- `service.go`: Основная логика сервиса управления доступом
- `session.go`: Управление сессиями пользователей

## Зависимости

- [GORM](https://gorm.io/): ORM библиотека для Go
- [bcrypt](golang.org/x/crypto/bcrypt): Для хеширования паролей
- [uuid](github.com/google/uuid): Для генерации уникальных идентификаторов

## Лицензия

[MIT License](https://opensource.org/licenses/MIT)

## Вклад в проект

Мы приветствуем вклад в развитие проекта! Пожалуйста, создавайте issues для обсуждения предлагаемых изменений перед отправкой pull request.

