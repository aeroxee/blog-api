# Blog API dengan Golang

Selamat datang di repository Blog API!

## Panduan Penggunaan

### Prasyarat

1. Pastikan Golang versi `1.21.5` telah terinstal di komputer Anda.
2. Clone repositori ini ke lokal komputer Anda.

```bash
git clone https://github.com/nama-akun-anda/blog-client-nextjs.git
```
3. Proyek ini dibuat untuk menjalankan API dari Web Client Blog. Silahkan clone [repo saya berikut ini](https://github.com/Aeroxee/blog-app)

### Instalasi

1. Buat file `.env`

```env
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=email anda
SMTP_PASSWORD=password app anda
HOSTNAME=ip address anda
PORT=8000
```

Sesuaikan dengan IP address anda.

2. Install dependensi

```bash
go mod tidy
```

3. Jalankan server.

```bash
go run cmd/blog-api/main.go
```

Server akan berjalan pada  `http:localhost:8000`

### Struktur Proyek

- `auth/`: Autentikasi dan enkripsi password.
- `cmd/`: File eksekusi.
- `docs/`: Dokumentasi dari Mkdocs.
- `handlers/`: Handler HTTP.
- `models/`: Model database.
- `controllers.go`: Controller group.
- `middlewares.go`: Middleware.
- `router.go`: Router dan Pengaturan server address.


## Kontribusi

Jika Anda ingin berkontribusi pada proyek ini, Anda dapat melakukan langkah-langkah berikut:

1. Fork repositori ini.
2. Lakukan perubahan pada forked repository Anda.
3. Submit pull request dengan deskripsi yang jelas tentang perubahan yang Anda lakukan.

## Lisensi

Proyek ini dilisensikan di bawah lisensi MIT - lihat file LICENSE untuk detailnya.

Terima kasih telah menggunakan Blog API! Semoga proyek ini bermanfaat untuk pengembangan aplikasi blog Anda. Jangan ragu untuk memberikan masukan atau melaporkan isu jika Anda menemui kendala. Selamat berkoding!