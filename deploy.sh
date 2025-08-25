#!/bin/bash

# WhatsApp Go MCP Deployment Script for Smithery.ai

echo "🚀 Deploying WhatsApp Go MCP to Smithery.ai..."

# Build the Go binary
echo "📦 Building Go binary..."
cd src
go build -o ../whatsapp .
cd ..

# Create deployment package
echo "📋 Creating deployment package..."
tar -czf whatsapp-go-mcp.tar.gz \
  whatsapp \
  smithery.yaml \
  src/statics \
  src/views \
  readme.md \
  docs/

# Deploy to Smithery
echo "☁️ Deploying to Smithery.ai..."
smithery deploy whatsapp-go-mcp.tar.gz

echo "✅ Deployment complete!"
echo ""
echo "📱 To use in Claude Desktop, add to your config:"
echo ""
echo '{
  "mcpServers": {
    "whatsapp": {
      "command": "npx",
      "args": [
        "-y",
        "@smithery/cli@latest",
        "run",
        "@samihalawa/whatsapp-go-mcp"
      ]
    }
  }
}'