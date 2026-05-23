# AICryptoAdvisor

AICryptoAdvisor — backend сервис для анализа криптовалют с использованием AI.

## 🚀 Stack

- Go
- Gin
- PostgreSQL
- pgxpool
- Gemini AI API
- REST API

---

# 📌 Features

## Coins

- Create coin
- Update coin price
- Delete coin
- Get all coins
- AI recommendation for coins

## Users

- Registration
- Login
- Update profile
- Change password
- Delete user

## AI Worker

Background worker automatically updates AI recommendations every 30 minutes.

---

# 📂 Project Structure

```bash
internal/
├── app
├── client
├── domain
├── handler
├── repository
├── service
├── worker
```

---

# ⚙️ Environment Variables

Create `.env` file:

```env
API_KEY=your_gemini_api_key

HOST=localhost
PORT=5432
USER_BASE=postgres
PASSWORD=123456
DB_NAME=aicryptoadvisor
```

---

# 🛠 Run Project

```bash
go mod tidy
go run ./cmd/main.go
```

---

# 📡 API Endpoints

## User

| Method | Endpoint      | Description   |
| ------ | ------------- | ------------- |
| POST   | /register     | Register user |
| POST   | /login        | Login         |
| GET    | /users/:email | Get user      |
| PUT    | /users/:email | Update user   |
| DELETE | /users/:email | Delete user   |

---

## Coins

| Method | Endpoint             | Description   |
| ------ | -------------------- | ------------- |
| POST   | /coins               | Create coin   |
| GET    | /coins               | Get all coins |
| GET    | /coins/:name         | Get coin      |
| PUT    | /coins/:name         | Update price  |
| DELETE | /coins/:name         | Delete coin   |
| POST   | /coins/:name/analyze | AI analyze    |

---

# 🧠 AI Recommendation Example

```json
{
	"recommendation": "BUY"
}
```

---

# 📌 Future Improvements

- JWT authentication
- Middleware
- Swagger documentation
- Unit tests
- CI/CD
- Redis caching
- WebSocket real-time updates

---

# 👨‍💻 Author

[Vana](https://github.com/Neviy)
