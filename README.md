# gofiber_outclass-api
## Deskripsi
Sebuah API untuk aplikasi [OutClass Web](https://github.com/azliR/vue_outclass) menggunakan Vue, dan aplikasi [OutClass Mobile](https://github.com/azliR/flutter_outclass). API ini dibuat dengan menggunakan [Golang](https://golang.org/) dan [Fiber](https://gofiber.io/). Di bagian database, API ini menggunakan [MongoDB](https://www.mongodb.com/) dan [Redis](https://redis.io/).

## Cara Penggunaan
Karena saya menggunakan Linux, maka cara penggunaan ini saya tulis untuk Linux. Jika ingin menggunakan Windows, silahkan sesuaikan penggunaannya dengan sistem operasi yang digunakan.
### Persyaratan
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Golang](https://golang.org/)

### Langkah-langkah
#### Clone Repository
1. Clone repository ini
```bash
git clone https://github.com/azliR/gofiber_outclass-api.git
```
2. Masuk ke folder repository
```bash
cd gofiber_outclass-api
```
#### Menjalankan MongoDB Replica Set Menggunakan Docker
3. Untuk menjalankan MongoDB Replica Set dengan authentication, kita perlu membuat keyfile terlebih dahulu. Jalankan perintah berikut untuk membuat keyfile
```bash
# Buat folder docker-config dan mongo
mkdir -p ./docker-config/mongo
# Membuat keyfile dengan panjang 756 karakter
# Ubah 756 dengan angka yang lebih besar jika ingin membuat keyfile yang lebih aman
openssl rand -base64 756 > ./docker-config/mongo/mongodb-keyfile
# Ubah permission keyfile agar hanya bisa dibaca dan ditulis oleh pemilik file
# Pelajari lebih lanjut tentang chmod di https://en.wikipedia.org/wiki/Chmod
sudo chmod 600 ./docker-config/mongo/mongodb-keyfile
```
4. Jalankan docker-compose. Perhatikan di file `docker-compose.yml`, kita meggunakan username `root` dan password `root` untuk MongoDB. Jika ingin mengganti, silahkan ubah di file `docker-compose.yml`
```bash
docker-compose up
```
5. Masuk ke Docker container untuk setup Replica Set MongoDB
```bash
docker exec -it outclass-api-fiber-mongo-1 mongosh -u <username> -p <password>
```
6. Disini kita hanya akan menggunakan satu node saja karena hanya akan memanfaatkan transaction. Jalankan perintah berikut untuk membuat Replica Set
```bash
rs.initiate()
```
7. Karena hanya untuk digunakan development, kita cukupkan disini. Untuk lebih lanjut, silahkan baca dokumentasi [MongoDB Replica Set](https://docs.mongodb.com/manual/tutorial/deploy-replica-set/).

#### Menjalankan API
8. Install semua dependency
```bash
go mod download
```
9. Buat file `.env` dan isi dengan data berikut
```bash
# Server settings:
SERVER_SCHEME="http"
SERVER_HOST="0.0.0.0"
SERVER_PORT=20109
SERVER_READ_TIMEOUT=60

# JWT settings:
JWT_SECRET_KEY="jwtsecret" # Ubah dengan key yang lebih aman
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT=15
JWT_REFRESH_KEY="jwtrefresh" # Ubah dengan key yang lebih aman
JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT=720

# Database settings:
MONGO_URI="mongodb://<username>:<password>@127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+1.6.0"

# Redis settings:
REDIS_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB_NUMBER=0
```
10. Jalankan server
```bash
go run main.go
```
11. API sudah bisa digunakan.