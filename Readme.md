# URL Shortener API

URL shortener with **Go, Gin, Postgres, Redis**.

**Base URL:** `http://localhost:8080`

---

## Endpoints

### 1. **Shorten URL**

**POST** `/shorten`

**Request Body:**

```json
{ "url": "https://example.com/long-url" }
```

**Response:**

```json
{
  "code": "abc123",
  "shortUrl": "http://localhost:8080/abc123",
  "original": "https://example.com/long-url",
  "clickCount": 0,
  "expirationDate": "2026-01-02T12:00:00Z",
  "createdAt": "2026-01-02T10:00:00Z"
}
```

---

### 2. **Resolve Short URL**

**GET** `/:code`

Redirects to the original URL.

- Status `302 Found` -> original URL
- Errors: `404` if code not found

---

### 3. **Get URL Stats**

**GET** `/stats/:code`

Returns info about a short URL including click count.

**Response:**

```json
{
  "code": "EdAZMFz",
  "shortUrl": "http://localhost:8080/EdAZMFz",
  "original": "http://facebook.com/hello/my/kitty/asdsa",
  "clickCount": 1,
  "expirationDate": "2027-01-02T19:31:50.83844Z",
  "createdAt": "2026-01-02T19:31:50.83844Z"
}
```

---

**Config** (env or `config.yaml`):

```env
PORT=8080
BASEURL=http://localhost:8080

DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASS=yourpassword
DB_NAME=shortener

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

- Redis caches short code -> original URL
- Retry on unique short code collisions (5 times max)

---
