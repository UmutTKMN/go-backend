# Go Web Projesi

Bu proje, modern web uygulamaları geliştirmek için Go programlama dili kullanılarak oluşturulmuş bir web uygulamasıdır. Aynı zamanda Go dilini öğrenme sürecimin bir parçası olarak geliştirilmektedir.

## 🚀 Özellikler

- Gin web framework'ü ile RESTful API desteği
- PostgreSQL veritabanı entegrasyonu (GORM ORM)
- JWT tabanlı kimlik doğrulama
- Güvenli şifreleme işlemleri
- Modüler ve ölçeklenebilir mimari

## 🛠 Teknolojiler

- Go 1.24.1
- Gin Web Framework
- GORM (PostgreSQL)
- JWT Authentication
- PostgreSQL

## 📁 Proje Yapısı

```
.
├── api/        # API endpoint'leri
├── cmd/        # Ana uygulama giriş noktası
├── configs/    # Yapılandırma dosyaları
├── deployments/# Deployment yapılandırmaları
├── docs/       # Dokümantasyon
├── internal/   # İç paketler
├── pkg/        # Paylaşılan paketler
├── scripts/    # Yardımcı scriptler
├── test/       # Test dosyaları
└── web/        # Web arayüzü dosyaları
```

## 🚀 Başlangıç

### Gereksinimler

- Go 1.24.1 veya üzeri
- PostgreSQL

### Kurulum

1. Projeyi klonlayın:

```bash
git clone https://github.com/UmutTKMN/gobackend.git
cd gobackend
```

2. Bağımlılıkları yükleyin:

```bash
go mod download
```

3. `.env` dosyasını yapılandırın:

```bash
cp .env.example .env
# .env dosyasını düzenleyin
```

4. Uygulamayı çalıştırın:

```bash
go run cmd/main.go
```

## 🔒 Güvenlik

- JWT tabanlı kimlik doğrulama
- Güvenli şifreleme işlemleri
- Çevresel değişkenler ile hassas bilgilerin yönetimi

## 📝 Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.

## 🤝 Katkıda Bulunma

1. Bu depoyu fork edin
2. Yeni bir özellik dalı oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'feat: Add some amazing feature'`)
4. Dalınıza push yapın (`git push origin feature/amazing-feature`)
5. Bir Pull Request oluşturun

## 📞 İletişim

Umut TKMN - [@UmutTKMN](https://github.com/UmutTKMN)

Proje Linki: [https://github.com/UmutTKMN/gobackend](https://github.com/UmutTKMN/gobackend)

## 📚 Öğrenme Kaynakları

Bu projede Go dilini öğrenirken faydalandığım kaynaklar:

- [Go Resmi Dokümantasyonu](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Tour of Go](https://tour.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Gin Framework Dokümantasyonu](https://gin-gonic.com/docs/)
- Claude AI - Kod geliştirme ve öğrenme sürecinde yapay zeka desteği
