package accessgo

import (
	"time"

	"gorm.io/gorm"
)

// User представляет пользователя в системе
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
	UserID   *uint
	GroupID  *uint
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
