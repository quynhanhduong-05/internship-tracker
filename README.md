# Internship Application Tracker API

A RESTful backend service built with Go to help users track internship applications, manage application statuses, and securely authenticate using JWT.

## 🚀 Features

- User authentication (Register / Login / Logout)
- JWT-based authentication with token revocation (logout)
- CRUD internship applications
- Authorization: users can only access their own data
- Pagination & filtering for application listing
- MySQL database with GORM ORM
- Clean modular project structure

---

## 🛠 Tech Stack

- Language: Go (net/http)
- Database: MySQL
- ORM: GORM
- Authentication: JWT (golang-jwt/jwt v5)
- Architecture: Handler / Service / Model separation

---

## 📂 Project Structure


---

## 🔐 Authentication Flow

- User logs in and receives a JWT access token
- Token is sent via `Authorization: Bearer <token>`
- Middleware validates token and injects `user_id` into request context
- Logout revokes token by storing it in a blacklist until expiration

---

## 📌 API Endpoints

### Auth
- POST `/intern/auth/register`
- POST `/intern/auth/login`
- POST `/intern/auth/logout`

### Applications
- POST `/applications`
- GET `/applications`
- PUT `/applications/{id}`
- DELETE `/applications/{id}`

---

## 📄 Pagination & Filtering

```http
GET /applications?page=2&limit=5&status=applied


---
