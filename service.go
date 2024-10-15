package accessgo

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AccessGoService представляет сервис для управления пользователями и группами
type AccessGoService struct {
	db *gorm.DB
}

// NewAccessGoService создает новый экземпляр AccessGoService
func NewAccessGoService(db *gorm.DB) (*AccessGoService, error) {
	if err := db.AutoMigrate(&User{}, &Group{}, &Access{}, &AccessLevel{}); err != nil {
		return nil, err
	}
	res := &AccessGoService{db: db}
	var cnt int64
	if err := db.Model(&Access{}).Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt == 0 {
		if err := res.SetupDefaultPermissions(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// CreateUser создает нового пользователя
func (s *AccessGoService) CreateUser(email, password, name string, userType UserType) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:                email,
		Password:             string(hashedPassword),
		Name:                 name,
		UserType:             string(userType),
		EmailValidationToken: uuid.NewString(),
	}

	result := s.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func (s *AccessGoService) UpdateUser(userID uint, email, password, name string, userType UserType) (*User, error) {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	user.Email = email
	user.Name = name
	user.UserType = string(userType)

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser удаляет пользователя
func (s *AccessGoService) DeleteUser(userID uint) error {
	result := s.db.Delete(&User{}, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}

// CreateGroup создает новую группу
func (s *AccessGoService) CreateGroup(name string) (*Group, error) {
	group := &Group{
		Name: name,
	}

	result := s.db.Create(group)
	if result.Error != nil {
		return nil, result.Error
	}

	return group, nil
}

// ValidateEmail проверяет токен валидации email
func (s *AccessGoService) ValidateEmail(token string) error {
	if token == "" {
		return errors.New("token is required")
	}
	var user User
	err := s.db.Where("email_validation_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	user.EmailValidate = true
	user.EmailValidationToken = ""
	return s.db.Save(&user).Error
}

// UpdateGroup обновляет информацию о группе
func (s *AccessGoService) UpdateGroup(groupID uint, name string) (*Group, error) {
	var group Group
	if err := s.db.First(&group, groupID).Error; err != nil {
		return nil, err
	}

	group.Name = name

	if err := s.db.Save(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// DeleteGroup удаляет группу
func (s *AccessGoService) DeleteGroup(groupID uint) error {
	result := s.db.Delete(&Group{}, groupID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("группа не найдена")
	}
	return nil
}

// CreateAccess создает новое право доступа
func (s *AccessGoService) CreateAccess(name, description string) (*Access, error) {
	access := &Access{
		Name:        name,
		Description: description,
	}

	result := s.db.Create(access)
	if result.Error != nil {
		return nil, result.Error
	}

	return access, nil
}

// UpdateAccess обновляет информацию о праве доступа
func (s *AccessGoService) UpdateAccess(accessID uint, name, description string) (*Access, error) {
	var access Access
	if err := s.db.First(&access, accessID).Error; err != nil {
		return nil, err
	}

	access.Name = name
	access.Description = description

	if err := s.db.Save(&access).Error; err != nil {
		return nil, err
	}

	return &access, nil
}

// DeleteAccess удаляет право доступа
func (s *AccessGoService) DeleteAccess(accessID uint) error {
	result := s.db.Delete(&Access{}, accessID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("право доступа не найдено")
	}
	return nil
}

// GetAccessByName получает право доступа по имени
func (s *AccessGoService) GetAccessByName(name string) (*Access, error) {
	var access Access
	if err := s.db.Where("name = ?", name).First(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("право доступа не найдено")
		}
		return nil, err
	}
	return &access, nil
}

// ListAccesses возвращает список всех прав доступа
func (s *AccessGoService) ListAccesses() ([]Access, error) {
	var accesses []Access
	if err := s.db.Find(&accesses).Error; err != nil {
		return nil, err
	}
	return accesses, nil
}

// AssignUserToGroup добавляет пользователя в группу
func (s *AccessGoService) AssignUserToGroup(userID, groupID uint) error {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("пользователь не найден")
	}

	var group Group
	if err := s.db.First(&group, groupID).Error; err != nil {
		return errors.New("группа не найдена")
	}

	if err := s.db.Model(&user).Association("Groups").Append(&group); err != nil {
		return err
	}

	return nil
}

// ExcludeUserFromGroup удаляет пользователя из группы
func (s *AccessGoService) ExcludeUserFromGroup(userID, groupID uint) error {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("пользователь не найден")
	}

	var group Group
	if err := s.db.First(&group, groupID).Error; err != nil {
		return errors.New("группа не найдена")
	}

	if err := s.db.Model(&user).Association("Groups").Delete(&group); err != nil {
		return err
	}

	return nil
}

// SetUserGroups устанавливает точный список групп для пользователя
func (s *AccessGoService) SetUserGroups(userID uint, groupIDs ...uint) error {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("пользователь не найден")
	}

	var groups []Group
	if len(groupIDs) > 0 {
		if err := s.db.Find(&groups, groupIDs).Error; err != nil {
			return err
		}

		if len(groups) != len(groupIDs) {
			return errors.New("одна или несколько групп не найдены")
		}
	}

	if err := s.db.Model(&user).Association("Groups").Replace(&groups); err != nil {
		return err
	}

	return nil
}

// GetUserGroups возвращает список групп пользователя
func (s *AccessGoService) GetUserGroups(userID uint) ([]Group, error) {
	var user User
	if err := s.db.Preload("Groups").First(&user, userID).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}

	return user.Groups, nil
}

// GetGroupUsers возвращает список пользователей в группе
func (s *AccessGoService) GetGroupUsers(groupID uint) ([]User, error) {
	var group Group
	if err := s.db.Preload("Users").First(&group, groupID).Error; err != nil {
		return nil, errors.New("группа не найдена")
	}

	return group.Users, nil
}

// AddUserAccessLevel добавляет уровень доступа пользователю
func (s *AccessGoService) AddUserAccessLevel(userID uint, accessName string) error {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("пользователь не найден")
	}

	var access Access
	if err := s.db.Where("name = ?", accessName).First(&access).Error; err != nil {
		return errors.New("право доступа не найдено")
	}

	accessLevel := AccessLevel{
		UserID:   &userID,
		AccessID: access.ID,
	}

	if err := s.db.Create(&accessLevel).Error; err != nil {
		return err
	}

	return nil
}

// RemoveUserAccessLevel удаляет уровень доступа у пользователя
func (s *AccessGoService) RemoveUserAccessLevel(userID uint, accessName string) error {
	var access Access
	if err := s.db.Where("name = ?", accessName).First(&access).Error; err != nil {
		return errors.New("право доступа не найдено")
	}

	result := s.db.Where("user_id = ? AND access_id = ?", userID, access.ID).Delete(&AccessLevel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("уровень доступа не найден у пользователя")
	}

	return nil
}

// AddGroupAccessLevel добавляет уровень доступа группе
func (s *AccessGoService) AddGroupAccessLevel(groupID uint, accessName string) error {
	var group Group
	if err := s.db.First(&group, groupID).Error; err != nil {
		return errors.New("группа не найдена")
	}

	var access Access
	if err := s.db.Where("name = ?", accessName).First(&access).Error; err != nil {
		return errors.New("право доступа не найдено")
	}

	accessLevel := AccessLevel{
		GroupID:  &groupID,
		AccessID: access.ID,
	}

	if err := s.db.Create(&accessLevel).Error; err != nil {
		return err
	}

	return nil
}

// RemoveGroupAccessLevel удаляет уровень доступа у группы
func (s *AccessGoService) RemoveGroupAccessLevel(groupID uint, accessName string) error {
	var access Access
	if err := s.db.Where("name = ?", accessName).First(&access).Error; err != nil {
		return errors.New("право доступа не найдено")
	}

	result := s.db.Where("group_id = ? AND access_id = ?", groupID, access.ID).Delete(&AccessLevel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("уровень доступа не найден у группы")
	}

	return nil
}

// CheckUserAccess проверяет, имеет ли пользователь указанный уровень доступа
func (s *AccessGoService) CheckUserAccess(userID uint, accessName string) (bool, error) {
	var user User
	if err := s.db.Preload("Accesses.Access").Preload("Groups.Accesses.Access").First(&user, userID).Error; err != nil {
		return false, errors.New("пользователь не найден")
	}

	// Проверяем прямые доступы пользователя
	for _, access := range user.Accesses {
		if access.Access.Name == accessName {
			return true, nil
		}
	}

	// Проверяем доступы групп пользователя
	for _, group := range user.Groups {
		for _, access := range group.Accesses {
			if access.Access.Name == accessName {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetUserSummaryAccessLevels возвращает все уровни доступа пользователя
func (s *AccessGoService) GetUserSummaryAccessLevels(userID uint) ([]string, error) {
	var user User
	if err := s.db.Preload("Accesses.Access").Preload("Groups.Accesses.Access").First(&user, userID).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}

	accessMap := make(map[string]bool)

	// Добавляем прямые доступы пользователя
	for _, access := range user.Accesses {
		accessMap[access.Access.Name] = true
	}

	// Добавляем доступы групп пользователя
	for _, group := range user.Groups {
		for _, access := range group.Accesses {
			accessMap[access.Access.Name] = true
		}
	}

	accessList := make([]string, 0, len(accessMap))
	for accessName := range accessMap {
		accessList = append(accessList, accessName)
	}

	return accessList, nil
}

// GetGroupAccessLevels возвращает все уровни доступа группы
func (s *AccessGoService) GetGroupAccessLevels(groupID uint) ([]string, error) {
	var group Group
	if err := s.db.Preload("Accesses.Access").First(&group, groupID).Error; err != nil {
		return nil, errors.New("группа не найдена")
	}

	accessList := make([]string, 0, len(group.Accesses))
	for _, access := range group.Accesses {
		accessList = append(accessList, access.Access.Name)
	}

	return accessList, nil
}

// GetUserAccessLevels возвращает все уровни доступа пользователя
func (s *AccessGoService) GetUserAccessLevels(userID uint) ([]string, error) {
	var user User
	if err := s.db.Preload("Accesses.Access").First(&user, userID).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}

	accessList := make([]string, 0, len(user.Accesses))
	for _, access := range user.Accesses {
		accessList = append(accessList, access.Access.Name)
	}
	return accessList, nil
}

// SetupDefaultPermissions создает стандартные права доступа
func (s *AccessGoService) SetupDefaultPermissions() error {
	defaultPermissions := []struct {
		Name        string
		Description string
	}{
		{"user:create", "Создание пользователя"},
		{"user:read", "Чтение информации о пользователе"},
		{"user:update", "Обновление информации о пользователе"},
		{"user:delete", "Удаление пользователя"},
		{"group:create", "Создание группы"},
		{"group:read", "Чтение информации о группе"},
		{"group:update", "Обновление информации о группе"},
		{"group:delete", "Удаление группы"},
		{"access:create", "Создание права доступа"},
		{"access:read", "Чтение информации о праве доступа"},
		{"access:update", "Обновление информации о праве доступа"},
		{"access:delete", "Удаление права доступа"},
		{"user_access:set", "Установка прав доступа пользователю"},
		{"group_access:set", "Установка прав доступа группе"},
	}

	for _, perm := range defaultPermissions {
		existingAccess := &Access{}
		err := s.db.Where("name = ?", perm.Name).First(existingAccess).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newAccess := &Access{
					Name:        perm.Name,
					Description: perm.Description,
				}
				if err := s.db.Create(newAccess).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

// GetUserByEmail возвращает пользователя по email
func (s *AccessGoService) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID возвращает пользователя по ID
func (s *AccessGoService) GetUserByID(userID uint) (*User, error) {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}
	return &user, nil
}

// GetAllUsers возвращает список всех пользователей
func (s *AccessGoService) GetAllUsers() ([]User, error) {
	var users []User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetGroupByID возвращает группу по ID
func (s *AccessGoService) GetGroupByID(groupID uint) (*Group, error) {
	var group Group
	if err := s.db.First(&group, groupID).Error; err != nil {
		return nil, errors.New("группа не найдена")
	}
	return &group, nil
}

// GetAllGroups возвращает список всех групп
func (s *AccessGoService) GetAllGroups() ([]Group, error) {
	var groups []Group
	if err := s.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateDefaultAdminUser создает пользователя-администратора с полными правами
func (s *AccessGoService) CreateDefaultAdminUser(email, password, name string) error {
	admin, err := s.CreateUser(email, password, name, UserTypeAdmin)
	if err != nil {
		return err
	}

	var accesses []Access
	if err := s.db.Find(&accesses).Error; err != nil {
		return err
	}

	for _, access := range accesses {
		if err := s.AddUserAccessLevel(admin.ID, access.Name); err != nil {
			return err
		}
	}
	return nil
}

// AuthenticateUser аутентифицирует пользователя по email и паролю
func (s *AccessGoService) AuthenticateUser(email, password string) (*User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !user.EmailValidate {
		return nil, errors.New("email не подтвержден")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("неверный пароль")
	}
	return user, nil
}
