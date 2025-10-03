# Smithery.ai Deployment Guide

## Overview
This WhatsApp MCP server is fully configured for deployment on Smithery.ai with HTTP Streamable transport.

## Features
- ✅ HTTP Streamable MCP transport (Smithery compatible)
- ✅ Stateless mode enabled for scalability
- ✅ CORS headers configured
- ✅ Session configuration support
- ✅ Docker containerization
- ✅ Health check endpoint
- ✅ 5 WhatsApp messaging tools

## Available Tools

### 1. whatsapp_send_text
Send text messages to WhatsApp contacts or groups.
- **Parameters:**
  - `phone` (required): Phone number or group ID
  - `message` (required): Text message to send
  - `is_forwarded` (optional): Whether message is forwarded
  - `reply_message_id` (optional): Message ID to reply to

### 2. whatsapp_send_contact
Send contact cards via WhatsApp.
- **Parameters:**
  - `phone` (required): Phone number or group ID
  - `contact_name` (required): Name of the contact
  - `contact_phone` (required): Phone number of the contact
  - `is_forwarded` (optional): Whether message is forwarded

### 3. whatsapp_send_link
Send links with captions via WhatsApp.
- **Parameters:**
  - `phone` (required): Phone number or group ID
  - `link` (required): URL to send
  - `caption` (required): Caption for the link
  - `is_forwarded` (optional): Whether message is forwarded

### 4. whatsapp_send_location
Send location coordinates via WhatsApp.
- **Parameters:**
  - `phone` (required): Phone number or group ID
  - `latitude` (required): Latitude coordinate (as string)
  - `longitude` (required): Longitude coordinate (as string)
  - `is_forwarded` (optional): Whether message is forwarded

### 5. whatsapp_send_image
Send images via WhatsApp.
- **Parameters:**
  - `phone` (required): Phone number or group ID
  - `image_url` (required): URL of the image
  - `caption` (optional): Caption for the image
  - `view_once` (optional): View once mode
  - `compress` (optional): Compress image (default: true)
  - `is_forwarded` (optional): Whether message is forwarded

## Deployment Configuration

### smithery.yaml
```yaml
runtime: "container"
build:
  dockerfile: "Dockerfile"
  dockerBuildPath: "."
startCommand:
  type: "http"
```

### Key Features for Smithery
1. **Port 8081**: Configured as required by Smithery
2. **Stateless Mode**: Enabled for better scalability
3. **CORS Support**: Full CORS headers for cross-origin requests
4. **Health Check**: Available at `/health`
5. **Tools Debug**: Available at `/tools`
6. **MCP Endpoint**: Available at `/mcp`

## Testing on Smithery

Once deployed, the server exposes:
- **MCP Endpoint**: `https://your-server.smithery.ai/mcp`
- **Health Check**: `https://your-server.smithery.ai/health`
- **Tools List**: `https://your-server.smithery.ai/tools`

## WhatsApp Authentication

On first run, you'll need to scan the QR code to authenticate with WhatsApp:
1. The server will generate a QR code
2. Open WhatsApp on your phone
3. Go to Settings > Linked Devices
4. Scan the QR code
5. The session will be saved for future use

## Persistent Storage

The server uses SQLite databases stored in `/app/storages/`:
- `whatsapp.db` - WhatsApp session data (authentication, contacts, groups)
- `chatstorage.db` - Chat messages and history

**IMPORTANT**: The Dockerfile declares `/app/storages` as a VOLUME for persistence.
- On Smithery deployment, this volume persists across container restarts
- Your WhatsApp session remains logged in even after server restarts
- Contact names, chat history, and media are preserved

## Environment Variables

The server uses these default values suitable for Smithery:
- `PORT=8081` (Set by Smithery)
- `MCP_HOST=0.0.0.0`
- Auto-reconnect enabled
- Session persistence enabled via SQLite volumes

## Troubleshooting

### Tools Not Showing
The server registers 5 WhatsApp tools on startup. If tools aren't visible:
1. Check `/health` endpoint for server status
2. Check `/tools` endpoint for tool list
3. Verify WhatsApp authentication is complete

### Connection Issues
The server includes auto-reconnection logic for WhatsApp connection stability.

## Architecture
- **Go 1.24**: High-performance runtime
- **whatsmeow**: Official WhatsApp Web protocol
- **MCP-Go**: Model Context Protocol implementation
- **Docker**: Container-based deployment
- **Alpine Linux**: Minimal container size