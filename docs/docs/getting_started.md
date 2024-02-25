---
title: Memulai
description: Memulai dengan RestFull API Aeroxee Blog
---

# Memulai

Untuk memulai menggunakan API ini, saya harap anda menggunakan Node.js atau console javascript. Semua response yang diberikan berupa `application/json`.

## URL PATH

Berikut URL path ini saya tampilkan agar anda mudah bekerja dalam menghandle API ini. Server ini berjalan pada port `8000`.

| *Method* | *Path*                                          | *Description*                             |
|----------|-------------------------------------------------|-------------------------------------------|
| GET      | /media/*filepath                                | URL media                                 |
| POST     | /v1/register                                    | URL untuk pendaftaran akun                |
| POST     | /v1/get-token                                   | URL Untuk mendapatkan token baru          |
| POST     | /v1/check-email                                 | URL untuk mencek email                    |
| POST     | /v1/check-username                              | URL untuk mencek username                 |
| POST     | /v1/upload                                      | URL untuk upload file                     |
| GET      | /v1/user/auth                                   | URL untuk mencek autentikasi info         |
| GET      | /v1/articles                                    | URL untuk melihat data artikel            |
| GET      | /v1/articles/:username/:slug                    | URL untuk melihat detail artikel          |
| GET      | /v1/articles/:username/:slug/comment            | URL untuk melihat komentar dari artikel   |
| GET      | /v1/articles/:username/:slug/comment/:commentId | URL untuk melihat komentar berdasarkan id |
| POST     | /v1/articles                                    | URL untuk membuat artikel baru            |
| PUT      | /v1/articles/:username/:slug                    | URL untuk memperbaharui artikel           |
| DELETE   | /v1/articles/:username/:slug                    | URL untuk menghapus artikel               |
| POST     | /v1/articles/:username/:slug/comment            | URL untuk membuat komentar                |
| PUT      | /v1/articles/:username/:slug/comment/:commentId | URL untuk mengedit komentar               |
| DELETE   | /v1/articles/:username/:slug/comment/:commentId | URL untuk menghapus komentar              |