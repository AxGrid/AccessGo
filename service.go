package accessgo

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserGoService представляет сервис для управления пользователями и группами
type UserGoService struct {
	db *gorm.DB
}

// NewUserGoService создает новый экземпляр UserGoService
func NewUserGoService(db *gorm.DB) *UserGoService {
	return &UserGoService{db: db}
}

// CreateUser создает нового пользователя
func (s *UserGoService) CreateUser(email, password, name string, userType UserType) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		UserType: string(userType),
	}

	result := s.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func (s *UserGoService) UpdateUser(userID uint, email, password, name string, userType UserType) (*User, error) {
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
func (s *UserGoService) DeleteUser(userID uint) error {
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
func (s *UserGoService) CreateGroup(name string) (*Group, error) {
	group := &Group{
		Name: name,
	}

	result := s.db.Create(group)
	if result.Error != nil {
		return nil, result.Error
	}

	return group, nil
}

// UpdateGroup обновляет информацию о группе
func (s *UserGoService) UpdateGroup(groupID uint, name string) (*Group, error) {
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
func (s *UserGoService) DeleteGroup(groupID uint) error {
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
func (s *UserGoService) CreateAccess(name, description string) (*Access, error) {
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
func (s *UserGoService) UpdateAccess(accessID uint, name, description string) (*Access, error) {
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
func (s *UserGoService) DeleteAccess(accessID uint) error {
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
func (s *UserGoService) GetAccessByName(name string) (*Access, error) {
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
func (s *UserGoService) ListAccesses() ([]Access, error) {
	var accesses []Access
	if err := s.db.Find(&accesses).Error; err != nil {
		return nil, err
	}
	return accesses, nil
}

// AssignUserToGroup добавляет пользователя в группу
func (s *UserGoService) AssignUserToGroup(userID, groupID uint) error {
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
func (s *UserGoService) ExcludeUserFromGroup(userID, groupID uint) error {
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
func (s *UserGoService) SetUserGroups(userID uint, groupIDs ...uint) error {
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
func (s *UserGoService) GetUserGroups(userID uint) ([]Group, error) {
	var user User
	if err := s.db.Preload("Groups").First(&user, userID).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}

	return user.Groups, nil
}

// GetGroupUsers возвращает список пользователей в группе
func (s *UserGoService) GetGroupUsers(groupID uint) ([]User, error) {
	var group Group
	if err := s.db.Preload("Users").First(&group, groupID).Error; err != nil {
		return nil, errors.New("группа не найдена")
	}

	return group.Users, nil
}

// AddUserAccessLevel добавляет уровень доступа пользователю
func (s *UserGoService) AddUserAccessLevel(userID uint, accessName string) error {
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
func (s *UserGoService) RemoveUserAccessLevel(userID uint, accessName string) error {
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
func (s *UserGoService) AddGroupAccessLevel(groupID uint, accessName string) error {
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
func (s *UserGoService) RemoveGroupAccessLevel(groupID uint, accessName string) error {
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
func (s *UserGoService) CheckUserAccess(userID uint, accessName string) (bool, error) {
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

// GetUserAccessLevels возвращает все уровни доступа пользователя
func (s *UserGoService) GetUserAccessLevels(userID uint) ([]string, error) {
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

// SetupDefaultPermissions создает стандартные права доступа
func (s *UserGoService) SetupDefaultPermissions() error {
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
func (s *UserGoService) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}

// CreateDefaultAdminUser создает пользователя-администратора с полными правами
func (s *UserGoService) CreateDefaultAdminUser(email, password, name string) error {
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
func (s *UserGoService) AuthenticateUser(email, password string) (*User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("неверный пароль")
	}

	return user, nil
}
