# AzizChatbot WhatsApp Bot 📱

A WhatsApp bot integration for AzizChatbot - allowing users to interact with the AI chatbot directly through WhatsApp.  Built with Go and the WhatsApp Web API.

## ✨ Features

- **WhatsApp Web Integration**: Uses whatsmeow library for WhatsApp Web API
- **QR Code Authentication**: Easy login via QR code scanning
- **Rate Limiting**: Redis integration for message rate limiting
- **Docker Support**: Multi-stage build for optimized container size

## 🛠️ Tech Stack

- **Language**: Go 1.23
- **WhatsApp API**: whatsmeow
- **Database**: PostgreSQL (via pgx driver)
- **Session Management**: Redis
- **QR Code**: qrterminal
- **Environment**: godotenv

## 📋 Prerequisites

- Go 1.23 or higher
- PostgreSQL database
- Redis server
- AI service running (see [azizchatbot-ai](https://github.com/AzizChatbot/azizchatbot-ai))

## 🚀 Installation

### Local Setup

1. Clone the repository:
   ```bash
   git clone https://github. com/AzizChatbot/azizchatbot-whatsapp-bot.git
   cd azizchatbot-whatsapp-bot
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file based on `. env.example`:
   ```env
   REDIS_URL="redis://localhost:6379"
   AI_URL="http://localhost:5000"
   DB_URL="postgresql://user:password@localhost/whatsapp_bot"
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5.  Scan the QR code with WhatsApp to authenticate. 

### Docker Setup

1. Build the Docker image:
   ```bash
   docker build -t azizchatbot-whatsapp-bot . 
   ```

2.  Run the container:
   ```bash
   docker run -it --env-file .env azizchatbot-whatsapp-bot
   ```

> **Note**: For QR code scanning, run the container in interactive mode (`-it`) to see the QR code in the terminal.

## 📁 Project Structure

```
azizchatbot-whatsapp-bot/
├── main.go               # Application entry point
├── lib/
│   └── msgHandler/       # Message handling logic
├── go.mod                # Go module definition
├── go. sum                # Dependency checksums
├── Dockerfile            # Multi-stage Docker build
└── . env. example          # Environment variables template
```

## 🔧 How It Works

1. **Initialization**: The bot connects to WhatsApp Web using stored credentials or prompts for QR code login
2. **Session Storage**: Device credentials are stored in PostgreSQL for persistent sessions
3.  **Message Handling**: Incoming messages are processed and sent to the AI service
4. **Response**: AI responses are sent back to the user via WhatsApp
5. **Group Protection**: Bot automatically leaves any group it's added to

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

<p align="center">Made with ❤️ for King Abdulaziz University</p>
