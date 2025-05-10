# 🚀 GoMailer - Email API with Scheduler & WorkerPool

*A lightweight Go API to send emails via REST or an HTMX-powered interface, featuring scheduled sending and concurrent processing with a goroutine-based worker pool.*

## ✨ Features
- **📤 Email sending** via REST API or HTML/HTMX interface
- **⏱ Scheduled sending** with a built-in scheduler
- **👷 WorkerPool** using goroutines for concurrent task processing
- **🎯 Rate limiting** to control throughput
- **📦 Built-in HTML templates**

## 🛠 Installation

### Prerequisites
- Go 1.24+
- SMTP server credentials (e.g., Gmail with app [password](https://itsupport.umd.edu/itsupport?id=kb_article_view&sysparm_article=KB0015112))
- Environment variables (create a `app.env` file):

```
FROM="your.gmail"
PASSWORD="your.app.password"
HOST="your.stmp.server"
PORT=587
```

### Run the project
```
git clone https://github.com/Moth13/Mailer.git
cd mailer
go mod download
go build -o moth-mailer .
./moth-mailer
```

or

```
make server
```

### Test the project
```
make test
```

## 🎮 Usage

### API Endpoints
**Send an immediate email (POST)**
```
curl -X POST http://localhost:8080/api/mail/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": "dest@example.com",
    "subject": "Hello GoMailer!",
    "body": "This is a test 🌟"
    "scheduled_at": "2025-05-10T15:15"
  }'
```

### HTMX Interface
Visit `http://localhost:8080/` to access the interactive form

## 📚 Dependencies
- [Gin](https://gin-gonic.com) - HTTP routing
- [Viper](https://github.com/spf13/viper) - Environment variable management
- [HTMX](https://htmx.org) - Dynamic frontend interface

## 🤝 Contributing
1. Fork the project
2. Create your feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Push to your branch
5. Open a Pull Request

## 📄 License
MIT License - See [LICENSE](LICENSE)

---

**Developed with ❤️ by [Moth13] | 2025**
[Full documentation](docs/) | [Advanced examples](examples/)
