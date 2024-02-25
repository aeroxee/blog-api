# Pendaftaran Akun

Pada URL path pendaftaran akun `/v1/register` ini menggunakan method `POST` dan body `application/json`.

## Payload

```json title="Payload"
{
    "first_name": "Aeroxee", // (1)
    "last_name": "Blog", // (2)
    "username": "aeroxee", // (3)
    "email": "youremail@gmail.com", // (4)
    "password": "yourpassword" // (5)
}
```

1. required(dibutuhkan).
2. required(dibutuhkan).
3. required(dibutuhkan).
4. required(dibutuhkan).
5. required(dibutuhkan).

## Kirim data

Disini saya akan mengirim data melalui Javascript saja.

```js title="Kirim data"
async function register() {
    const response = await fetch(`http://localhost:8000/v1/register`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json" // (1)
        },
        body: JSON.stringify({ // (2)
            first_name: "Your first name",
            last_name: "Your last name",
            username: "Your username",
            email: "Your email",
            password: "Your password"
        })
    })

    const data = await response.json()

    console.log(data)
}
```

1. Pastikan kirim header `Content-Type` dengan nilai `application/json`.
2. Dan kirim body data dengan format json.

Berikut output atau response jika pendaftaran berhasil.

```json
{
  "message": "Berhasil mendaftarkan akun dengan alaman email: blabla@mail.com",
  "status": "success",
  "user": {
    "id": 2,
    "first_name": "Aeroxee",
    "last_name": "Blog",
    "username": "ae",
    "email": "blabla@mail.com",
    // ---
  }
}
```