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
	mcpServer.AddTool(c.toolGetMessages(), c.handleGetMessages)
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

	// Call actual service
	response, err := c.chatService.ListChats(ctx, domainChat.ListChatsRequest{
		Limit:  limit,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get chat list: %w", err)
	}

	// Format for AI consumption with proper JSON structure
	result := fmt.Sprintf("Found %d chats:\n\n", len(response.Data))
	
	for i, chat := range response.Data {
		chatType := "Contact"
		if chat.IsGroup {
			chatType = "Group"
		}
		
		result += fmt.Sprintf("%d. %s (%s)\n", i+1, chat.Name, chatType)
		result += fmt.Sprintf("   JID: %s\n", chat.JID)
		
		if chat.LastMessage != "" {
			sender := "You"
			if chat.LastMessageFrom != "" && chat.LastMessageFrom != "Me" {
				sender = chat.LastMessageFrom
			}
			result += fmt.Sprintf("   Last: %s: %s\n", sender, chat.LastMessage)
		}
		
		if chat.LastMessageTime != "" {
			result += fmt.Sprintf("   Time: %s\n", chat.LastMessageTime)
		}
		
		if chat.UnreadCount > 0 {
			result += fmt.Sprintf("   Unread: %d messages\n", chat.UnreadCount)
		}
		
		if chat.IsPinned {
			result += "   📌 Pinned\n"
		}
		
		if chat.IsArchived {
			result += "   🗄️ Archived\n"
		}
		
		result += "\n"
	}
	
	result += fmt.Sprintf("\nTotal: %d chats (showing %d)\n", response.Pagination.Total, len(response.Data))

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

	resp, err := c.chatService.ArchiveChat(ctx, domainChat.ArchiveChatRequest{
		ChatJID: phone,
		Archive: archive,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to archive chat: %w", err)
	}

	return mcp.NewToolResultText(resp.Message), nil
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

	resp, err := c.chatService.MarkChatAsRead(ctx, domainChat.MarkChatAsReadRequest{
		ChatJID: phone,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to mark chat as read: %w", err)
	}

	return mcp.NewToolResultText(resp.Message), nil
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

	resp, err := c.chatService.DeleteChat(ctx, domainChat.DeleteChatRequest{
		ChatJID:     phone,
		KeepStarred: keepStarred,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to delete chat: %w", err)
	}

	return mcp.NewToolResultText(resp.Message), nil
}

func (c *ChatHandler) toolGetMessages() mcp.Tool {
	return mcp.NewTool("whatsapp_get_messages",
		mcp.WithDescription("Get recent messages from a chat."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Number of messages to retrieve (default: 10)"),
		),
	)
}

func (c *ChatHandler) handleGetMessages(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	limit := 10
	if l, ok := request.GetArguments()["limit"].(float64); ok {
		limit = int(l)
	}

	response, err := c.chatService.GetChatMessages(ctx, domainChat.GetChatMessagesRequest{
		ChatJID: phone,
		Limit:   limit,
		Offset:  0,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	if len(response.Data) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No messages found for chat %s", phone)), nil
	}

	// Format messages for AI consumption
	result := fmt.Sprintf("Messages from %s (%d messages):\n", phone, len(response.Data))
	for i, msg := range response.Data {
		sender := msg.SenderJID
		if msg.IsFromMe {
			sender = "Me"
		}
		result += fmt.Sprintf("%d. [%s] %s: %s\n", i+1, msg.Timestamp, sender, msg.Content)
		if msg.MediaType != "" && msg.MediaType != "text" {
			result += fmt.Sprintf("   Type: %s", msg.MediaType)
			if msg.Filename != "" {
				result += fmt.Sprintf(" | File: %s", msg.Filename)
			}
			result += "\n"
		}
	}
	
	return mcp.NewToolResultText(result), nil
}