package accessgo

import (
	"time"

	"gorm.io/gorm"
)

// User представляет пользователя в системе
type User struct {
	gorm.Model
	Email                string        `gorm:"unique;not null"`
	EmailValidate        bool          `gorm:"not null;default:false"`
	EmailValidationToken string        `gorm:"size:64;index:idx_user_email_validation_token"`
	Password             string        `gorm:"size:64; not null"`
	Name                 string        `gorm:"size:255; not null"`
	UserType             string        `gorm:"size:15;not null"`
	CreatedAt            time.Time     `gorm:"not null"`
	UpdatedAt            time.Time     `gorm:"not null;index:idx_user_updated_at"`
	Accesses             []AccessLevel `gorm:"foreignKey:UserID"`
	Groups               []Group       `gorm:"many2many:user_groups;"`
}

// Group представляет группу пользователей
type Group struct {
	gorm.Model
	Name     string        `gorm:"unique;not null"`
	Accesses []AccessLevel `gorm:"foreignKey:GroupID"`
	Users    []User        `gorm:"many2many:user_groups;"`
}

// Access представляет право доступа
type Access struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
}

// AccessLevel представляет уровень доступа для пользователя или группы
type AccessLevel struct {
	gorm.Model
	AccessID uint
	UserID   *uint  `gorm:"uniqueIndex:idx_user_group_access"`
	GroupID  *uint  `gorm:"uniqueIndex:idx_user_group_access"`
	Access   Access `gorm:"foreignKey:AccessID"`
}

// UserType представляет типы пользователей
type UserType string

const (
	UserTypeAdmin    UserType = "admin"
	UserTypeEmployee UserType = "employee"
	UserTypeUser     UserType = "user"
	UserTypeBlocked  UserType = "blocked"
)

type Session struct {
	ID         string
	UserID     int
	CreatedAt  time.Time
	ExpiresAt  time.Time
	IsLongTerm bool
}
