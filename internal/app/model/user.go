package model

import (
	"database/sql"
	"time"
)

// User temel kullanıcı modelini tanımlar
type User struct {
	// Temel Bilgiler
	ID             uint   `json:"id" gorm:"primaryKey"`
	Username       string `json:"username" gorm:"unique;not null"`
	Email          string `json:"email" gorm:"unique;not null"`
	Password       string `json:"-" gorm:"not null"` // JSON'da gösterilmeyecek
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DisplayName    string `json:"display_name"`
	PhoneNumber    string `json:"phone_number"`
	SecondaryEmail string `json:"secondary_email"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio" gorm:"type:text"`

	// Adres Bilgileri
	AddressLine1  string  `json:"address_line1"`
	AddressLine2  string  `json:"address_line2"`
	City          string  `json:"city"`
	StateProvince string  `json:"state_province"`
	Country       string  `json:"country"`
	PostalCode    string  `json:"postal_code"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`

	// Zaman Bilgileri
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	LastLogin          *time.Time `json:"last_login"`
	LastPasswordChange *time.Time `json:"last_password_change"`
	LastActivity       *time.Time `json:"last_activity"`
	RegistrationDate   time.Time  `json:"registration_date"`
	BirthDate          *time.Time `json:"birth_date"`
	Age                int        `json:"age" gorm:"-"` // Database'de saklanmayacak, hesaplanacak

	// Durum ve Doğrulama
	IsActive           bool       `json:"is_active" gorm:"default:true"`
	IsVerified         bool       `json:"is_verified" gorm:"default:false"`
	VerificationToken  string     `json:"-"`
	VerificationSentAt *time.Time `json:"verification_sent_at"`
	EmailVerifiedAt    *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt    *time.Time `json:"phone_verified_at"`
	TwoFactorEnabled   bool       `json:"two_factor_enabled" gorm:"default:false"`
	LoginAttempts      int        `json:"login_attempts" gorm:"default:0"`
	AccountLockedUntil *time.Time `json:"account_locked_until"`

	// Tercihler
	PreferredLanguage       string         `json:"preferred_language" gorm:"default:'tr'"`
	Timezone                string         `json:"timezone" gorm:"default:'Europe/Istanbul'"`
	NotificationPreferences sql.NullString `json:"notification_preferences" gorm:"type:json"`
	PrivacySettings         sql.NullString `json:"privacy_settings" gorm:"type:json"`
	ThemePreference         string         `json:"theme_preference" gorm:"default:'light'"`
	AccessibilitySettings   sql.NullString `json:"accessibility_settings" gorm:"type:json"`

	// Sistem Bilgileri
	AccountType         string         `json:"account_type" gorm:"default:'standard'"`
	SubscriptionID      string         `json:"subscription_id"`
	SubscriptionStatus  string         `json:"subscription_status"`
	SubscriptionEndDate *time.Time     `json:"subscription_end_date"`
	AccountLimits       sql.NullString `json:"account_limits" gorm:"type:json"`
	ApiKey              string         `json:"-"`
	ApiUsageCount       int            `json:"api_usage_count" gorm:"default:0"`
	IpRegistration      string         `json:"ip_registration"`
	ReferralCode        string         `json:"referral_code"`
	ReferredBy          string         `json:"referred_by"`

	// Sosyal Medya Entegrasyonu
	FacebookID   string         `json:"facebook_id"`
	GoogleID     string         `json:"google_id"`
	TwitterID    string         `json:"twitter_id"`
	LinkedinID   string         `json:"linkedin_id"`
	GithubID     string         `json:"github_id"`
	AppleID      string         `json:"apple_id"`
	SocialLogins sql.NullString `json:"social_logins" gorm:"type:json"`

	// İlişkiler
	UserRequests       []UserRequest       `json:"user_requests,omitempty" gorm:"foreignKey:UserID"`
	UserDocuments      []UserDocument      `json:"user_documents,omitempty" gorm:"foreignKey:UserID"`
	UserDevices        []UserDevice        `json:"user_devices,omitempty" gorm:"foreignKey:UserID"`
	UserPreferences    []UserPreference    `json:"user_preferences,omitempty" gorm:"foreignKey:UserID"`
	UserCommunications []UserCommunication `json:"user_communications,omitempty" gorm:"foreignKey:UserID"`

	// Rol ve Personel İlişkileri
	Roles []Role `json:"roles,omitempty" gorm:"many2many:user_roles"`
	Staff *Staff `json:"staff,omitempty" gorm:"foreignKey:UserID"`

	// Oluşturduğu Roller
	CreatedRoles []Role `json:"created_roles,omitempty" gorm:"foreignKey:CreatedBy"`
}

// UserRequest kullanıcı talep tablosu
type UserRequest struct {
	ID                  uint       `json:"id" gorm:"primaryKey"`
	UserID              uint       `json:"user_id"`
	RequestType         string     `json:"request_type"`
	RequestDetails      string     `json:"request_details" gorm:"type:text"`
	Priority            string     `json:"priority"`
	Status              string     `json:"status"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	ProcessingStartedAt *time.Time `json:"processing_started_at"`
	CompletedAt         *time.Time `json:"completed_at"`
	HandledByStaffID    uint       `json:"handled_by_staff_id"`
	DecisionNotes       string     `json:"decision_notes" gorm:"type:text"`
	FollowUpRequired    bool       `json:"follow_up_required"`
	FollowUpDate        *time.Time `json:"follow_up_date"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// İlişkiler
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserDocument kullanıcı belge tablosu
type UserDocument struct {
	ID                 uint       `json:"id" gorm:"primaryKey"`
	UserID             uint       `json:"user_id"`
	DocumentType       string     `json:"document_type"`
	DocumentPath       string     `json:"document_path"`
	UploadDate         time.Time  `json:"upload_date"`
	VerificationStatus string     `json:"verification_status"`
	VerifiedBy         uint       `json:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at"`
	ExpiryDate         *time.Time `json:"expiry_date"`
	Notes              string     `json:"notes" gorm:"type:text"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// İlişkiler
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserDevice kullanıcı cihaz tablosu
type UserDevice struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id"`
	DeviceType    string    `json:"device_type"`
	DeviceName    string    `json:"device_name"`
	DeviceModel   string    `json:"device_model"`
	OSVersion     string    `json:"os_version"`
	AppVersion    string    `json:"app_version"`
	DeviceToken   string    `json:"device_token"`
	LastIP        string    `json:"last_ip"`
	LastLoginDate time.Time `json:"last_login_date"`
	IsTrusted     bool      `json:"is_trusted"`
	GeoLocation   string    `json:"geo_location"`
	UserAgent     string    `json:"user_agent" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// İlişkiler
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserPreference kullanıcı tercih tablosu
type UserPreference struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id"`
	PreferenceKey   string    `json:"preference_key"`
	PreferenceValue string    `json:"preference_value"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`

	// İlişkiler
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserCommunication kullanıcı iletişim tablosu
type UserCommunication struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	UserID            uint       `json:"user_id"`
	CommunicationType string     `json:"communication_type"`
	SentAt            time.Time  `json:"sent_at"`
	DeliveryStatus    string     `json:"delivery_status"`
	ReadAt            *time.Time `json:"read_at"`
	Subject           string     `json:"subject"`
	Content           string     `json:"content" gorm:"type:text"`
	AttachmentPaths   string     `json:"attachment_paths"`
	SenderID          uint       `json:"sender_id"`
	Importance        string     `json:"importance"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// İlişkiler
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
