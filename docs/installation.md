# Kurulum Rehberi

Bu belge, Go Web projesinin kurulum ve yapılandırma adımlarını açıklar.

## Gereksinimler

Projeyi çalıştırmak için aşağıdaki gereksinimlere ihtiyacınız vardır:

- Go 1.24.1 veya üzeri
- PostgreSQL 12 veya üzeri
- Git

## Adım 1: Projeyi Klonlama

```bash
git clone https://github.com/UmutTKMN/go-backend.git
cd go-backend
```

## Adım 2: Bağımlılıkları Yükleme

```bash
go mod download
```

## Adım 3: Veritabanı Kurulumu

1. PostgreSQL'i yükleyin ve çalıştırın:

   ```bash
   # Ubuntu için
   sudo apt install postgresql postgresql-contrib
   sudo systemctl start postgresql.service
   
   # macOS için (Homebrew ile)
   brew install postgresql
   brew services start postgresql
   
   # Windows için
   # PostgreSQL'i resmi web sitesinden indirip yükleyin
   ```

2. Veritabanını oluşturun:

   ```bash
   sudo -u postgres psql
   ```

   PostgreSQL komut satırında:

   ```sql
   CREATE DATABASE go-backend;
   CREATE USER go-backenduser WITH ENCRYPTED PASSWORD 'yourpassword';
   GRANT ALL PRIVILEGES ON DATABASE go-backend TO go-backenduser;
   \q
   ```

## Adım 4: Çevresel Değişkenleri Yapılandırma

1. `.env.example` dosyasını `.env` olarak kopyalayın:

   ```bash
   cp .env.example .env
   ```

2. `.env` dosyasını düzenleyerek kendi yapılandırmanızı ayarlayın:

   ```
   # Web Server
   PORT=8080
   
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=go-backenduser
   DB_PASSWORD=yourpassword
   DB_NAME=go-backend
   DB_SSLMODE=disable
   
   # JWT
   JWT_SECRET=your_jwt_secret_key
   JWT_EXPIRATION=24h
   ```

## Adım 5: Uygulamayı Çalıştırma

```bash
go run cmd/main.go
```

Uygulama varsayılan olarak `http://localhost:8080` adresinde çalışmaya başlayacaktır.

## Geliştirme Ortamı İçin Ek Adımlar

### Hot Reload ile Geliştirme

Geliştirme sırasında kod değişikliklerinde otomatik yeniden yükleme için [Air](https://github.com/cosmtrek/air) kullanabilirsiniz:

```bash
# Air yükleme
go install github.com/cosmtrek/air@latest

# Uygulamayı Air ile çalıştırma
air
```

### Veritabanı Migrations

Veritabanı şemasını otomatik olarak oluşturmak için GORM Auto Migration kullanılmaktadır. Uygulama ilk kez çalıştırıldığında gerekli tablolar otomatik olarak oluşturulacaktır.

## Sorun Giderme

### Veritabanı Bağlantı Hatası

Eğer veritabanına bağlanırken hata alıyorsanız:

1. PostgreSQL servisinin çalıştığını kontrol edin
2. `.env` dosyasındaki veritabanı bilgilerinin doğru olduğunu doğrulayın
3. Veritabanı kullanıcısının yeterli izinlere sahip olduğunu kontrol edin

### Port Çakışması

Eğer belirtilen port başka bir uygulama tarafından kullanılıyorsa, `.env` dosyasında `PORT` değerini değiştirin. 