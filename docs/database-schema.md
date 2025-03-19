# Veritabanı Şeması

Bu belge, Go Web projesinin veritabanı modellerini ve ilişkilerini açıklar.

## Modeller

### User Model

Kullanıcı hesaplarını temsil eder.

| Alan | Tür | Açıklama |
|------|-----|----------|
| id | serial | Birincil anahtar |
| username | string | Kullanıcı adı (benzersiz) |
| email | string | E-posta adresi (benzersiz) |
| password | string | Şifrelenmiş parola |
| created_at | timestamp | Oluşturulma zamanı |
| updated_at | timestamp | Son güncelleme zamanı |

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"size:255;not null;unique"`
    Email     string    `gorm:"size:255;not null;unique"`
    Password  string    `gorm:"size:255;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

### Profile Model

Kullanıcıların profil bilgilerini temsil eder.

| Alan | Tür | Açıklama |
|------|-----|----------|
| id | serial | Birincil anahtar |
| user_id | integer | User tablosu ile ilişki |
| full_name | string | Tam ad |
| bio | text | Biyografi |
| avatar_url | string | Profil resmi URL'si |
| created_at | timestamp | Oluşturulma zamanı |
| updated_at | timestamp | Son güncelleme zamanı |

```go
type Profile struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID"`
    FullName  string    `gorm:"size:255"`
    Bio       string    `gorm:"type:text"`
    AvatarURL string    `gorm:"size:255"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

### Post Model

Kullanıcıların paylaşımlarını temsil eder.

| Alan | Tür | Açıklama |
|------|-----|----------|
| id | serial | Birincil anahtar |
| user_id | integer | User tablosu ile ilişki |
| title | string | Başlık |
| content | text | İçerik |
| created_at | timestamp | Oluşturulma zamanı |
| updated_at | timestamp | Son güncelleme zamanı |

```go
type Post struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID"`
    Title     string    `gorm:"size:255;not null"`
    Content   string    `gorm:"type:text;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

## İlişkiler

### User - Profile (Bire Bir)

Her kullanıcının bir profili vardır.

```go
// User modelinden Profile'a erişim
user.Profile

// Profile modelinden User'a erişim
profile.User
```

### User - Post (Bire Çok)

Bir kullanıcı birden fazla paylaşım yapabilir.

```go
// User modelinden Post'lara erişim
user.Posts

// Post modelinden User'a erişim
post.User
```

## Veritabanı Diyagramı

```
+----------------+       +----------------+       +----------------+
|      User      |       |     Profile    |       |      Post      |
+----------------+       +----------------+       +----------------+
| id             |<----->| id             |       | id             |
| username       |       | user_id        |       | user_id        |
| email          |       | full_name      |       | title          |
| password       |       | bio            |       | content        |
| created_at     |       | avatar_url     |       | created_at     |
| updated_at     |       | created_at     |       | updated_at     |
+----------------+       | updated_at     |       +----------------+
                         +----------------+              ^
                                                        |
                                                        |
                                                        +
                                                   +----------------+
                                                   |      User      |
                                                   +----------------+
```

## Migrations

GORM Auto Migration kullanılarak veritabanı şeması otomatik olarak oluşturulmaktadır. Uygulama ilk kez çalıştırıldığında, veritabanı modelleri için gerekli tablolar otomatik olarak oluşturulur.

Migrations kodu aşağıdaki gibidir:

```go
func InitDatabase() (*gorm.DB, error) {
    // Veritabanı bağlantısı
    db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Auto Migration
    err = db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Post{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
} 