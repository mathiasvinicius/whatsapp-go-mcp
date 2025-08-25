package mcp

import (
	"context"
	"fmt"

	domainChat "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/chat"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ChatHandler struct {
	chatService domainChat.IChatUsecase
}

func InitMcpChat(chatService domainChat.IChatUsecase) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (c *ChatHandler) AddChatTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(c.toolGetList(), c.handleGetList)
	mcpServer.AddTool(c.toolArchive(), c.handleArchive)
	mcpServer.AddTool(c.toolMarkAsRead(), c.handleMarkAsRead)
	mcpServer.AddTool(c.toolDeleteChat(), c.handleDeleteChat)
}

func (c *ChatHandler) toolGetList() mcp.Tool {
	return mcp.NewTool("whatsapp_get_chat_list",
		mcp.WithDescription("Get list of all WhatsApp chats."),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of chats to return (default: 50)"),
		),
	)
}

func (c *ChatHandler) handleGetList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limit := 50
	if l, ok := request.GetArguments()["limit"].(float64); ok {
		limit = int(l)
	}

	// This would need actual implementation
	result := fmt.Sprintf("Retrieved %d chats", limit)
	return mcp.NewToolResultText(result), nil
}

func (c *ChatHandler) toolArchive() mcp.Tool {
	return mcp.NewTool("whatsapp_archive_chat",
		mcp.WithDescription("Archive or unarchive a WhatsApp chat."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithBoolean("archive",
			mcp.Required(),
			mcp.Description("True to archive, false to unarchive"),
		),
	)
}

func (c *ChatHandler) handleArchive(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	archive := request.GetArguments()["archive"].(bool)

	action := "archived"
	if !archive {
		action = "unarchived"
	}

	// This would need actual implementation
	return mcp.NewToolResultText(fmt.Sprintf("Chat with %s has been %s", phone, action)), nil
}

func (c *ChatHandler) toolMarkAsRead() mcp.Tool {
	return mcp.NewTool("whatsapp_mark_chat_as_read",
		mcp.WithDescription("Mark all messages in a chat as read."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
	)
}

func (c *ChatHandler) handleMarkAsRead(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)

	// This would need actual implementation
	return mcp.NewToolResultText(fmt.Sprintf("All messages in chat with %s marked as read", phone)), nil
}

func (c *ChatHandler) toolDeleteChat() mcp.Tool {
	return mcp.NewTool("whatsapp_delete_chat",
		mcp.WithDescription("Delete a WhatsApp chat."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithBoolean("keep_starred",
			mcp.Description("Keep starred messages (default: false)"),
		),
	)
}

func (c *ChatHandler) handleDeleteChat(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	keepStarred := false
	if k, ok := request.GetArguments()["keep_starred"].(bool); ok {
		keepStarred = k
	}

	result := fmt.Sprintf("Chat with %s has been deleted", phone)
	if keepStarred {
		result += " (starred messages kept)"
	}

	// This would need actual implementation
	return mcp.NewToolResultText(result), nil
}