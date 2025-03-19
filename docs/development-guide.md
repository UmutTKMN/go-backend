# Geliştirme Rehberi

Bu belge, Go Web projesi için geliştirme süreçlerini ve standartlarını açıklar.

## Geliştirme Ortamı Kurulumu

### Gereksinimler

- Go 1.24.1 veya üzeri
- Git
- PostgreSQL
- IDE (önerilen: VSCode, GoLand)
- Air (Hot Reload için)

### VSCode için Önerilen Eklentiler

- Go (ms-vscode.go)
- Go Test Explorer
- PostgreSQL
- ENV
- Docker
- GitLens

## Kod Standartları

### Go Kod Stili

- [Effective Go](https://golang.org/doc/effective_go) kurallarını izleyin
- Paket adları tekil olmalıdır (`user`, `auth` vb.)
- Çok kelimeli paket adlarından kaçının
- Değişken isimleri anlamlı ve kısa olmalıdır
- Açıklayıcı yorumlar ekleyin

### Kod Düzeni

```go
package user

import (
    "fmt"
    "time"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// User represents a user in the system
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"size:255;not null;unique"`
    Email     string    `gorm:"size:255;not null;unique"`
    Password  string    `gorm:"size:255;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// GetByID fetches a user by ID
func GetByID(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.First(&user, id)
    return &user, result.Error
}
```

### Commit Mesajları

Commit mesajlarınızı [Conventional Commits](https://www.conventionalcommits.org/) standartlarına göre yazın:

```
feat: yeni kullanıcı profili özelliği eklendi
fix: kullanıcı kimlik doğrulama hatası düzeltildi
docs: API dokümantasyonu güncellendi
style: kod formatı düzeltildi
refactor: kullanıcı servisi yeniden yapılandırıldı
test: kullanıcı testi eklendi
chore: bağımlılıklar güncellendi
```

## Git İş Akışı

1. Ana dalda çalışmayın, her özellik için yeni bir dal oluşturun
2. Dal adlandırması şu formatta olmalıdır: `<type>/<description>`
   - Örnek: `feature/user-profile`, `bugfix/auth-error`
3. Yerel olarak çalışmalarınızı sık sık commit edin
4. Değişikliklerinizi uzak depoya push etmeden önce ana dalı pull edin ve merge yapın
5. Pull request oluşturmadan önce kodunuzu gözden geçirin

## Test Yazma

### Birim Testleri

Her paket için birim testleri yazın:

```go
package user

import (
    "testing"
)

func TestGetByID(t *testing.T) {
    // Test kodu
}
```

### Entegrasyon Testleri

API endpointleri için entegrasyon testleri yazın:

```go
package api_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
)

func TestUserEndpoint(t *testing.T) {
    // Test kodu
}
```

## API Geliştirme

### Yeni Endpoint Ekleme

1. `api/handlers` altında yeni bir işleyici fonksiyonu oluşturun
2. `api/routes` altındaki ilgili router dosyasına endpoint'i ekleyin
3. Gerekirse `internal/services` altında yeni bir servis fonksiyonu oluşturun
4. API testleri yazın

### Middleware Kullanımı

Kimlik doğrulama ve loglama için middleware'leri kullanın:

```go
// api/routes/user.go
func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
    userGroup := router.Group("/api/users")
    userGroup.Use(middleware.AuthRequired())
    
    userGroup.GET("/:id", handlers.GetUser)
    userGroup.PUT("/:id", handlers.UpdateUser)
    userGroup.DELETE("/:id", handlers.DeleteUser)
}
```

## Veritabanı İşlemleri

### Migrasyon

Model yapısını değiştirdiğinizde, otomatik migrasyon kullanın:

```go
// internal/database/db.go
func InitDatabase() (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    err = db.AutoMigrate(&models.User{}, &models.Profile{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

### Repository Pattern

Veritabanı işlemleri için repository pattern kullanın:

```go
// internal/repositories/user_repository.go
type UserRepository struct {
    DB *gorm.DB
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    result := r.DB.First(&user, id)
    return &user, result.Error
}
```

## Hata İşleme

### Standart Hata Yanıtları

```go
// pkg/errors/errors.go
func HandleError(c *gin.Context, err error) {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": err.Error(),
    })
}
```

### Doğrulama Hataları

```go
// api/validators/user_validator.go
func ValidateUserInput(c *gin.Context) (*models.UserInput, error) {
    var input models.UserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        return nil, err
    }
    
    // Doğrulama kuralları
    
    return &input, nil
}
```

## Dağıtım (Deployment)

### Docker

`deployments/docker/Dockerfile` dosyasını kullanarak container oluşturun:

```bash
docker build -t go-backend -f deployments/docker/Dockerfile .
docker run -p 8080:8080 go-backend
```

### Kubernetes

`deployments/kubernetes` klasöründeki yapılandırmaları kullanın:

```bash
kubectl apply -f deployments/kubernetes/deployment.yaml
kubectl apply -f deployments/kubernetes/service.yaml
```

## Performans Optimizasyonu

- N+1 sorgu probleminden kaçının, GORM Preload kullanın
- İndeksleri akıllıca kullanın
- Büyük veri kümeleri için sayfalama yapın
- Redis gibi önbellekleme mekanizmaları kullanın

## Güvenlik

- JWT tokenları için kısa ömürler kullanın
- HTTPS kullanın
- SQL enjeksiyonlarından kaçının
- Kullanıcı şifrelerini bcrypt ile hashleyin
- Rate limiting uygulayın 