package model

import (
	"database/sql"
	"time"
)

// Role kullanıcı rollerini tanımlar
type Role struct {
	ID              uint      `json:"id" gorm:"primaryKey;column:role_id"`
	RoleName        string    `json:"role_name" gorm:"unique;not null"`
	Description     string    `json:"description" gorm:"type:text"`
	PermissionLevel int       `json:"permission_level" gorm:"not null;default:0"`
	IsSystemRole    bool      `json:"is_system_role" gorm:"not null;default:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uint      `json:"created_by"`

	// İlişkiler
	Creator *User   `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Staff   []Staff `json:"staff,omitempty" gorm:"foreignKey:RoleID"`
	Users   []User  `json:"users,omitempty" gorm:"many2many:user_roles"`
}

// Staff personel bilgilerini tanımlar
type Staff struct {
	ID               uint           `json:"id" gorm:"primaryKey;column:staff_id"`
	UserID           uint           `json:"user_id" gorm:"not null;unique"`
	RoleID           uint           `json:"role_id" gorm:"not null"`
	Department       string         `json:"department"`
	Position         string         `json:"position"`
	HireDate         time.Time      `json:"hire_date"`
	ManagerID        *uint          `json:"manager_id"`
	EmployeeStatus   string         `json:"employee_status" gorm:"default:'active'"`
	AccessLevel      int            `json:"access_level" gorm:"default:1"`
	ShiftPattern     string         `json:"shift_pattern"`
	Permissions      sql.NullString `json:"permissions" gorm:"type:json"`
	StartDate        time.Time      `json:"start_date"`
	EndDate          *time.Time     `json:"end_date"`
	EmergencyContact string         `json:"emergency_contact"`
	Notes            string         `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`

	// İlişkiler
	User    *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role    *Role  `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Manager *Staff `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
}
