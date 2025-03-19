# API Dokümantasyonu

Bu belge, Go Web projesinin sunduğu RESTful API'leri açıklar.

## Genel Bilgiler

- Base URL: `http://localhost:8080` (Geliştirme ortamı)
- Tüm API istekleri JSON formatında veri döndürür
- Hata durumunda, API standart bir hata yanıtı döndürür

### Kimlik Doğrulama

API isteklerinin çoğu için JWT tabanlı kimlik doğrulama gereklidir. Token'ı istek başlıklarında aşağıdaki gibi gönderin:

```
Authorization: Bearer <your_jwt_token>
```

## Kullanıcı Endpointleri

### Kullanıcı Kaydı

```
POST /api/auth/register
```

Yeni bir kullanıcı hesabı oluşturur.

**İstek Gövdesi:**

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Başarılı Yanıt (200 OK):**

```json
{
  "id": "integer",
  "username": "string",
  "email": "string",
  "created_at": "timestamp"
}
```

### Kullanıcı Girişi

```
POST /api/auth/login
```

Kullanıcı kimlik bilgilerini doğrular ve bir JWT token döndürür.

**İstek Gövdesi:**

```json
{
  "email": "string",
  "password": "string"
}
```

**Başarılı Yanıt (200 OK):**

```json
{
  "token": "string",
  "user": {
    "id": "integer",
    "username": "string",
    "email": "string"
  }
}
```

### Kullanıcı Bilgilerini Alma

```
GET /api/users/me
```

Oturum açmış kullanıcının bilgilerini döndürür.

**Başarılı Yanıt (200 OK):**

```json
{
  "id": "integer",
  "username": "string",
  "email": "string",
  "created_at": "timestamp"
}
```

## Örnek İstek ve Yanıtlar

### Kullanıcı Kaydı Örneği

**İstek:**

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securePassword123"
  }'
```

**Yanıt:**

```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "created_at": "2023-01-01T12:00:00Z"
}
```

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| 400 | Geçersiz istek (Bad Request) |
| 401 | Kimlik doğrulama hatası (Unauthorized) |
| 403 | Yetkisiz erişim (Forbidden) |
| 404 | Kaynak bulunamadı (Not Found) |
| 500 | Sunucu hatası (Internal Server Error) | 