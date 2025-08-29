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
	mcpServer.AddTool(m.toolReact(), m.handleReact)
	mcpServer.AddTool(m.toolDelete(), m.handleDelete)
	mcpServer.AddTool(m.toolGetMessages(), m.handleGetMessages)
	mcpServer.AddTool(m.toolMarkAsRead(), m.handleMarkAsRead)
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

func (m *MessageHandler) toolGetMessages() mcp.Tool {
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

func (m *MessageHandler) handleGetMessages(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	limit := 10
	if l, ok := request.GetArguments()["limit"].(float64); ok {
		limit = int(l)
	}

	// TODO: This needs to use ChatService.GetChatMessages instead of MessageService
	// For now, return a descriptive message about the limitation
	return mcp.NewToolResultText(fmt.Sprintf("Message retrieval for %s requires chat service integration (limit: %d). Use whatsapp_get_chat_list to see available chats first.", phone, limit)), nil
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