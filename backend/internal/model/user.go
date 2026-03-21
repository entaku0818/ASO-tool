package model

import "time"

type User struct {
	ID                   string     `json:"id"`
	Email                string     `json:"email"`
	PasswordHash         string     `json:"-"`
	Name                 string     `json:"name"`
	IsAdmin              bool       `json:"is_admin"`
	Plan                 string     `json:"plan"` // "free" or "pro"
	StripeCustomerID     string     `json:"-"`
	StripeSubscriptionID string     `json:"-"`
	PlanExpiresAt        *time.Time `json:"plan_expires_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
}

func (u *User) IsPro() bool {
	if u.Plan != "pro" {
		return false
	}
	if u.PlanExpiresAt != nil && u.PlanExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if len(r.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}
