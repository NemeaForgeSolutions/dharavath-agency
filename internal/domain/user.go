package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	RoleAdmin = "admin"
	RoleAgent = "agent"
	RoleUser  = "user"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email is already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("forbidden: insufficient privileges")
	ErrLoginRestricted    = errors.New("login is currently restricted to administrators and agents only")
)

type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"` // "admin", "agent", "user"
	Phone        string     `json:"phone,omitempty"`
	Avatar       string     `json:"avatar,omitempty"`
	IsActive     bool       `json:"isActive"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"`
	CreatedBy    *string    `json:"createdBy,omitempty"`
	UpdatedBy    *string    `json:"updatedBy,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

func (u *User) IsAdmin() bool {
	return u != nil && strings.ToLower(u.Role) == RoleAdmin
}

func (u *User) IsAgent() bool {
	return u != nil && (strings.ToLower(u.Role) == RoleAgent || u.IsAdmin())
}

func (u *User) Sanitized() *User {
	if u == nil {
		return nil
	}
	clone := *u
	clone.PasswordHash = ""
	return &clone
}
