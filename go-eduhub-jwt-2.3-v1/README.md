# Web Application

## Live Coding - Go EduHub JWT 3

### Implementation technique

Siswa akan melaksanakan sesi live code di 15 menit terakhir dari sesi mentoring dan di awasi secara langsung oleh Mentor. Dengan penjelasan sebagai berikut:

- **Durasi**: 15 menit pengerjaan
- **Submit**: Maximum 10 menit setelah sesi mentoring menggunakan `grader-cli submit`
- **Obligation**: Wajib melakukan _share screen_ di breakout room yang akan dibuatkan oleh Mentor pada saat mengerjakan Live Coding.

### Description

**Go Eduhub JWT** adalah sebuah aplikasi yang dirancang untuk membantu pengelolaan dan manajemen data siswa dan kursus menggunakan bahasa pemrograman Go. Aplikasi ini memungkinkan pengguna untuk melakukan berbagai operasi seperti menambah dan menampilkan data siswa juga menambah kursus yang terkait dengan siswa tersebut.

Dalam live-code ini, kita akan mengimplementasikan API menggunakan _Golang web framework Gin_ untuk mengelola data _student_ dan _course_. API harus mengizinkan client untuk:

- Registrasi pengguna baru
- Login menggunakan user yang telah didaftarkan
- Menghapus kursus dari daftar kursus siswa

Disini sudah ditentukan endpoint untuk setiap operasi untuk mengimplementasikan logika yang diperlukan dari setiap operasi menggunakan repository student dan course di file `main.go` dengan endpoint group sebagai berikut:

```go
users := gin.Group("/user")
{
  users.POST("/login", apiHandler.UserAPIHandler.Login)
  users.POST("/register", apiHandler.UserAPIHandler.Register)
}

course := gin.Group("/course")
{
  course.Use(middleware.Auth())
  course.DELETE("/delete/:course_id", apiHandler.CourseAPIHandler.DeleteCourse)
}
```

### Constraints

Pada live code ini, kamu harus melengkapi fungsi dari repository dan handler api `student` dan `course` sebagai berikut:

📁 **repository**

Ini adalah fungsi yang berinteraksi dengan database Postgres menggunakan GORM:

- `repository/user.go`
  - `GetUserByEmail`: Function ini menggunakan library GORM untuk mengambil data pengguna berdasarkan alamat email yang diberikan sebagai argumen. Pertama-tama, function akan mengeksekusi sebuah query `SELECT` untuk mencari pengguna dengan email yang cocok di dalam tabel `users`. Query tersebut akan menggunakan klausa `WHERE` dengan kondisi email yang diberikan.
    - Jika pengguna dengan email yang cocok ditemukan, data pengguna akan diassign ke variabel `user` yang merupakan objek dari model `User`. Function akan mengembalikan `user` dan `nil` sebagai error.
    - Namun jika tidak ditemukan pengguna dengan email yang cocok, function akan mengembalikan `user` kosong dan `nil` sebagai error.
    - Jika terjadi error lain selama proses tersebut, function akan mengembalikan error yang terjadi.

- `repository/course.go`
  - `Delete`: Function ini akan menghapus data kursus yang memiliki `id` yang sesuai dengan nilai yang diberikan sebagai argumen. Pertama-tama, function akan mengeksekusi sebuah query untuk menghapus data kursus pada tabel `courses` dengan `id` yang sesuai.
    - Jika proses tersebut berhasil, function akan mengembalikan `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, function akan mengembalikan `error` yang terjadi.

📁 **middleware**

Di file `middleware/auth.go` terdapat fungsi `Auth()` yang digunakan untuk melakukan autentikasi pengguna dengan menggunakan JWT (JSON Web Token). Middleware ini berfungsi untuk mengecek apakah user yang mengakses suatu endpoint atau route tertentu sudah terotentikasi atau belum. Fungsi ini terdiri dari beberapa langkah:

- Mengambil cookie dengan nama session_token dari request dengan key `session_token`. Cookie ini berisi JWT token yang digunakan untuk autentikasi.
- Parsing JWT token pada cookie tersebut untuk mendapatkan claims yang berisi informasi mengenai user ID. JWT token pada cookie tersebut akan di-parse menggunakan JWT library pada Go, yaitu jwt-go. Setelah di-parse, claims pada token tersebut akan dimasukkan ke dalam struct `Claims`.

  ```go
  type Claims struct {
    UserID int `json:"user_id"`
    jwt.StandardClaims
  }
  ```

  Claims pada JWT token ini dapat berisi informasi user yang terotentikasi seperti user ID, username, dan lain-lain. Disini, hanya user ID yang dimasukkan ke dalam context.
- Menentukan respons HTTP berdasarkan hasil parsing JWT token dan keberadaan cookie.
  - Jika parsing token gagal, maka akan mengembalikan respon HTTP dengan status code 401 atau 400 tergantung dari jenis error yang terjadi. Jika token tidak valid, maka akan mengembalikan respon HTTP dengan status code 401.
  - Jika cookie session_token tidak ada, maka akan mengembalikan respon HTTP dengan status code 401 jika request memiliki header Content-Type dengan nilai "application/json", atau melakukan redirect ke halaman login jika tidak.
- Menyimpan nilai UserID dari claims ke dalam context dengan key "id". Nilai UserID ini nantinya akan dapat digunakan di handler atau endpoint selanjutnya.
- Setelah semua langkah selesai, middleware akan memanggil Next untuk melanjutkan request ke handler atau endpoint selanjutnya.

📁 **api**

- `api/user.go`
  - method `Login`: adalah sebuah handler yang menerima parameter `*gin.Context`. Method ini akan melakukan login user dengan memanggil userService.Login dengan parameter context dan `*model.User` yang sudah didapatkan dari body request.
    - Method ini wajib mengirim data json dengan contoh format sebagai berikut:

      ```json
      {
        "email": <string>,
        "password": <string>
      }
      ```

    - Jika data `email` atau `password` kosong maka method ini akan mengembalikan response dengan status code `400` dan pesan error sebagai berikut:

      ```json
      {
        "error": "email or password is empty"
      }
      ```

    - Jika terjadi error saat menggunakan `userService.Login`, maka method ini akan mengembalikan response dengan status code `500` dan pesan error sebagai berikut:

      ```json
      {
        "error": "error internal server"
      }
      ```

    - Jika user berhasil login, maka method ini akan membuat token JWT dengan `UserID` sebagai payload dan `expirationTime` sebagai waktu kadaluwarsa. Setelah itu, token JWT akan di-sign dengan menggunakan `model.JwtKey`.
      - Setelah token JWT berhasil di-sign, method akan membuat `cookie` baru dengan nama `session_token` dan value `tokenString` yang sudah didapatkan dari JWT sebelumnya. Jika cookie dengan nama `session_token` sudah ada, maka value cookie tersebut akan diganti dengan `tokenString` yang baru.
      - Method akan mengembalikan response dengan status code `200` dan data user yang sudah login. Jika sukses, maka response akan berisi status `200` dan data JSON berikut:

        ```json
        {
          "user_id": <int>,
          "message": "login success"
        }
        ```

- `api/course.go`
  - `DeleteCourse`: fungsi ini digunakan untuk menghapus data course yang sudah ada di dalam sistem. Course yang akan dihapus diidentifikasi melalui `courseID` yang diambil dari parameter URL.
    - Jika terjadi error dalam proses validasi ID atau proses penghapusan data dari database, maka API akan mengembalikan response JSON dengan status HTTP `400` Bad Request atau `500` Internal Server Error masing-masing beserta pesan error yang dihasilkan.
    - Jika operasi berhasil dilakukan, maka API akan mengembalikan response JSON dengan status HTTP `404` Not Found (terkesan sedikit keliru, karena seharusnya HTTP status code yang digunakan adalah `200` OK) dan pesan sukses.

### Perhatian

Sebelum kalian menjalankan `grader-cli test`, pastikan kalian sudah mengubah database credentials pada file **`main.go`** (line 30) dan **`main_test.go`** (line 49) sesuai dengan database kalian. Kalian cukup mengubah nilai dari  `"username"`, `"password"` dan `"database_name"`saja.

Contoh:

```go
dbCredentials = Credential{
    Host:         "localhost",
    Username:     "postgres", // <- ubah ini
    Password:     "postgres", // <- ubah ini
    DatabaseName: "kampusmerdeka", // <- ubah ini
    Port:         5432,
}
```

### Test Case Examples

#### Test Case 1

**Input**:

- Login dengan `email` dan `password` kosong:

  ```http
  POST /user/login HTTP/1.1
  Host: localhost:8080
  Content-Type: application/json

  {
      "email": "",
      "password": ""
  }
  ```

**Output:**

- Login:

  ```http
  HTTP/1.1 400 Bad Request
  Content-Type: application/json

  {
      "error": "email or password is empty"
  }
  ```

#### Test Case 2

**Input**:

- Delete ID `course` yang valid:

  ```http
  DELETE /course/delete/1 HTTP/1.1
  Host: localhost:8080

**Output:**

- Data `course` berhasil dihapus:

  ```http
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
      "message": "course delete success"
  }
  ```

#### Test Case 3

**Input**:

- Delete ID `course` yang **tidak valid**:

  ```http
  DELETE /course/delete/abc HTTP/1.1
  Host: localhost:8080

**Output:**

- Gagal menghapus `course` karena ID `course` tidak valid:

  ```http
  HTTP/1.1 400 Bad Request
  Content-Type: application/json

  {
      "error": "Invalid course ID"
  }
  ```
  