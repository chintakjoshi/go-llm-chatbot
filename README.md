
# Chintak's Portfolio Chatbot - Configuration & Testing Guide

## 📋 Project Overview
This project consists of two main components:

- **Backend**: Go server with dual LLM providers (NVIDIA NIM + OpenRouter fallback)
- **Frontend**: React portfolio website with integrated chatbot

---

## 🔧 Backend Configuration

### Prerequisites
- Go 1.21 or higher
- NVIDIA NIM API account
- OpenRouter account (optional, for fallback)

### Environment Variables
Create a `.env` file in the `chatbot-backend/` directory:

```env
# Server Configuration
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-in-production
ALLOWED_ORIGINS=https://chintakjoshi.github.io,http://localhost:3000
RATE_LIMIT=10
RATE_LIMIT_WINDOW=60

# LLM Providers
NVIDIA_API_KEY=your-nvidia-nim-api-key-here
OPENROUTER_API_KEY=your-openrouter-api-key-here

# LLM Endpoints
NVIDIA_ENDPOINT=https://integrate.api.nvidia.com/v1/chat/completions
OPENROUTER_ENDPOINT=https://openrouter.ai/api/v1/chat/completions

# LLM Models
NVIDIA_MODEL=nvidia/nvidia-nemotron-nano-9b-v2
OPENROUTER_MODEL=z-ai/glm-4.5-air:free

# Logging
LOG_LEVEL=INFO
```

### Getting API Keys

- **NVIDIA NIM API Key**:
  - Visit [NVIDIA API Catalog](https://developer.nvidia.com/) 
  - Sign up for an account
  - Generate API key for the DeepSeek model
  - Copy the key to `NVIDIA_API_KEY`

- **OpenRouter API Key (Optional Fallback)**:
  - Visit [OpenRouter](https://openrouter.ai/)
  - Create an account
  - Generate an API key
  - Copy the key to `OPENROUTER_API_KEY`

### Backend Setup & Running

```bash
# Navigate to backend directory
cd chatbot-backend

# Install dependencies
go mod tidy

# Build the server
go build -o chatbot-server

# Run the server
./chatbot-server

# Or run directly with Go
go run main.go
```

### Backend Health Check

```bash
# Test if backend is running
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "ok",
#   "service": "chintak-chatbot",
#   "time": 1234567890,
#   "providers": "NVIDIA NIM (primary) + OpenRouter (fallback)"
# }
```

### Backend Logs
Logs are stored in `logs/backend.log` with the following structure:

```text
[timestamp] LEVEL file:line - message
```

Available log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`

---

## 🎨 Frontend Configuration

### Environment Variables
Create a `.env` file in your React project root:

```env
REACT_APP_CHATBOT_API_URL=http://localhost:8080/api/v1
REACT_APP_CHATBOT_API_KEY=portfolio-chatbot-key
```

**Note**: For production deployment on GitHub Pages, update the API URL to your deployed backend.

### Frontend Setup & Running

```bash
# Install dependencies (if not already done)
npm install

# Start development server
npm start

# Build for production
npm run build

# Deploy to GitHub Pages
npm run deploy
```

---

## 🧪 Testing the System

1. **Authentication Test**

```bash
# Test simple authentication
curl http://localhost:8080/api/v1/auth/simple

# Expected response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIs...",
#   "expires_in": 86400
# }
```

2. **Chat Endpoint Test**

```bash
# Get authentication token first
TOKEN=$(curl -s http://localhost:8080/api/v1/auth/simple | jq -r '.token')

# Send chat message
curl -X POST http://localhost:8080/api/v1/chat   -H "Authorization: Bearer $TOKEN"   -H "Content-Type: application/json"   -d '{
    "message": "Tell me about your experience",
    "session_id": "test-session-123"
  }'

# Expected response:
# {
#   "response": "I have 4+ years of experience...",
#   "session_id": "test-session-123",
#   "timestamp": 1234567890
# }
```

3. **Provider Testing**

```bash
# Test NVIDIA provider specifically
curl -X POST http://localhost:8080/api/v1/chat   -H "Authorization: Bearer $TOKEN"   -H "Content-Type: application/json"   -d '{"message": "Hello"}'

# Check logs to see which provider was used
tail -f logs/backend.log
```

4. **Frontend Integration Test**
- Start both backend and frontend servers
- Open [http://localhost:3000](http://localhost:3000) in browser
- Click the chat button in the bottom-right corner
- Send test messages:
  - "Who are you?"
  - "What projects have you built?"
  - "Tell me about your skills"

---

## 🔍 Troubleshooting

### Common Issues

- **Backend not starting**:
  - Check if port 8080 is available
  - Verify all required environment variables are set
  - Check `logs/backend.log` for errors

- **Authentication failures**:
  - Verify `JWT_SECRET` is set in backend `.env`
  - Check if API keys are valid and not expired

- **Chat not working**:
  - Check browser console for errors (F12 → Console)
  - Verify backend is running and accessible
  - Check CORS settings in backend

- **Provider failures**:
  - NVIDIA NIM: Check API key and rate limits
  - OpenRouter: Check daily request limit (50/day free tier)
  - Monitor logs for provider switching

### Log Analysis
Check `logs/backend.log` for:

```bash
# Successful request
[INFO] Successfully used NVIDIA NIM provider - Duration: 1.234s

# Provider failure
[WARN] Primary provider (NVIDIA NIM) failed: API rate limit exceeded

# Fallback activation
[INFO] Attempting failover to OpenRouter

# Authentication
[INFO] Successful authentication from 192.168.1.100
```

### Performance Monitoring

```bash
# Monitor slow requests ( > 1 second)
grep "Slow Request" logs/backend.log

# Monitor errors
grep "ERROR" logs/backend.log

# Monitor provider usage
grep "Provider" logs/backend.log
```

---

## 🚀 Deployment Checklist

### Backend Deployment
- Set production environment variables
- Update `ALLOWED_ORIGINS` with production domain
- Use strong `JWT_SECRET`
- Configure proper logging level (`LOG_LEVEL=INFO`)
- Set up process manager (PM2, systemd)
- Configure reverse proxy (nginx, Apache)

### Frontend Deployment
- Update `REACT_APP_CHATBOT_API_URL` to production backend
- Build and test locally before deployment
- Deploy to GitHub Pages: `npm run deploy`

### Post-Deployment Verification
- Backend health check returns `{"status": "ok"}`
- Frontend loads without console errors
- Chatbot responds to basic questions
- Logs are being generated properly
- No 404 or CORS errors in browser console

---

## 📞 Support
If you encounter issues:

- Check the logs: `logs/backend.log`
- Verify all environment variables are set
- Test endpoints with curl commands above
- Ensure both servers are running on correct ports:
  - Backend: [http://localhost:8080](http://localhost:8080)
  - Frontend: [http://localhost:3000](http://localhost:3000)

---

## 🎯 Quick Start Summary
- **Backend**: Set `.env` → `go build` → `./chatbot-server`
- **Frontend**: Set `.env` → `npm start`
- **Test**: Use curl commands or web interface
- **Deploy**: Update URLs → Build → Deploy

The system is now ready for use! 🎉