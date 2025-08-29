package mcp

import (
	"context"
	"fmt"

	domainMessage "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/message"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MessageHandler struct {
	messageService domainMessage.IMessageUsecase
}

func InitMcpMessage(messageService domainMessage.IMessageUsecase) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (m *MessageHandler) AddMessageTools(mcpServer *server.MCPServer) {
	// Message operations (get_messages moved to ChatHandler for proper service access)
	mcpServer.AddTool(m.toolMarkAsRead(), m.handleMarkAsRead)
	mcpServer.AddTool(m.toolReact(), m.handleReact)
	mcpServer.AddTool(m.toolDelete(), m.handleDelete)
	
	// Advanced message management
	mcpServer.AddTool(m.toolUpdate(), m.handleUpdate)
	mcpServer.AddTool(m.toolRevoke(), m.handleRevoke)
	mcpServer.AddTool(m.toolStar(), m.handleStar)
	mcpServer.AddTool(m.toolUnstar(), m.handleUnstar)
	mcpServer.AddTool(m.toolDownloadMedia(), m.handleDownloadMedia)
}

func (m *MessageHandler) toolReact() mcp.Tool {
	return mcp.NewTool("whatsapp_react_message",
		mcp.WithDescription("React to a WhatsApp message with an emoji."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithString("message_id",
			mcp.Required(),
			mcp.Description("ID of the message to react to"),
		),
		mcp.WithString("emoji",
			mcp.Required(),
			mcp.Description("Emoji reaction (e.g., 👍, ❤️, 😂)"),
		),
	)
}

func (m *MessageHandler) handleReact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageID := request.GetArguments()["message_id"].(string)
	emoji := request.GetArguments()["emoji"].(string)

	_, err := m.messageService.ReactMessage(ctx, domainMessage.ReactionRequest{
		Phone:     phone,
		MessageID: messageID,
		Emoji:     emoji,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Reacted with %s to message", emoji)), nil
}

func (m *MessageHandler) toolDelete() mcp.Tool {
	return mcp.NewTool("whatsapp_delete_message",
		mcp.WithDescription("Delete a WhatsApp message."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithString("message_id",
			mcp.Required(),
			mcp.Description("ID of the message to delete"),
		),
	)
}

func (m *MessageHandler) handleDelete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageID := request.GetArguments()["message_id"].(string)

	err := m.messageService.DeleteMessage(ctx, domainMessage.DeleteRequest{
		Phone:     phone,
		MessageID: messageID,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message %s deleted successfully", messageID)), nil
}


func (m *MessageHandler) toolMarkAsRead() mcp.Tool {
	return mcp.NewTool("whatsapp_mark_as_read",
		mcp.WithDescription("Mark messages as read."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithArray("message_ids",
			mcp.Required(),
			mcp.Description("Array of message IDs to mark as read"),
		),
	)
}

func (m *MessageHandler) handleMarkAsRead(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageIDs := request.GetArguments()["message_ids"].([]interface{})

	ids := make([]string, len(messageIDs))
	for i, id := range messageIDs {
		ids[i] = id.(string)
	}

	// Mark each message as read individually
	for _, id := range ids {
		_, err := m.messageService.MarkAsRead(ctx, domainMessage.MarkAsReadRequest{
			Phone:     phone,
			MessageID: id,
		})
		if err != nil {
			return nil, err
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("Marked %d messages as read", len(ids))), nil
}

// ===== ADVANCED MESSAGE MANAGEMENT TOOLS =====

func (m *MessageHandler) toolUpdate() mcp.Tool {
	return mcp.NewTool("whatsapp_update_message",
		mcp.WithDescription("Update/edit a WhatsApp message."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithString("message_id",
			mcp.Required(),
			mcp.Description("ID of the message to update"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("New message content"),
		),
	)
}

func (m *MessageHandler) handleUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageID := request.GetArguments()["message_id"].(string)
	message := request.GetArguments()["message"].(string)

	_, err := m.messageService.UpdateMessage(ctx, domainMessage.UpdateMessageRequest{
		Phone:     phone,
		MessageID: messageID,
		Message:   message,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message %s updated successfully", messageID)), nil
}

func (m *MessageHandler) toolRevoke() mcp.Tool {
	return mcp.NewTool("whatsapp_revoke_message",
		mcp.WithDescription("Revoke/recall a WhatsApp message for everyone."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithString("message_id",
			mcp.Required(),
			mcp.Description("ID of the message to revoke"),
		),
	)
}

func (m *MessageHandler) handleRevoke(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageID := request.GetArguments()["message_id"].(string)

	_, err := m.messageService.RevokeMessage(ctx, domainMessage.RevokeRequest{
		Phone:     phone,
		MessageID: messageID,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message %s revoked successfully", messageID)), nil
}

func (m *MessageHandler) toolStar() mcp.Tool {
	return mcp.NewTool("whatsapp_star_message",
		mcp.WithDescription("Star a WhatsApp message for bookmarking."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithString("message_id",
			mcp.Required(),
			mcp.Description("ID of the message to star"),
		),
	)
}

func (m *MessageHandler) handleStar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageID := request.GetArguments()["message_id"].(string)

	err := m.messageService.StarMessage(ctx, domainMessage.StarRequest{
		Phone:     phone,
		MessageID: messageID,
		IsStarred: true,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message %s starred successfully", messageID)), nil
}

func (m *MessageHandler) toolUnstar() mcp.Tool {
	return mcp.NewTool("whatsapp_unstar_message",
		mcp.WithDescription("Unstar a WhatsApp message (remove bookmark)."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithString("message_id",
			mcp.Required(),
			mcp.Description("ID of the message to unstar"),
		),
	)
}

func (m *MessageHandler) handleUnstar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageID := request.GetArguments()["message_id"].(string)

	err := m.messageService.StarMessage(ctx, domainMessage.StarRequest{
		Phone:     phone,
		MessageID: messageID,
		IsStarred: false,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message %s unstarred successfully", messageID)), nil
}

func (m *MessageHandler) toolDownloadMedia() mcp.Tool {
	return mcp.NewTool("whatsapp_download_media",
		mcp.WithDescription("Download media from a WhatsApp message (image, video, audio, document)."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID"),
		),
		mcp.WithString("message_id",
			mcp.Required(),
			mcp.Description("ID of the message containing media"),
		),
	)
}

func (m *MessageHandler) handleDownloadMedia(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	messageID := request.GetArguments()["message_id"].(string)

	response, err := m.messageService.DownloadMedia(ctx, domainMessage.DownloadMediaRequest{
		Phone:     phone,
		MessageID: messageID,
	})
	
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Media downloaded successfully:\nType: %s\nFilename: %s\nPath: %s\nSize: %d bytes", 
		response.MediaType, response.Filename, response.FilePath, response.FileSize)
	
	return mcp.NewToolResultText(result), nil
}