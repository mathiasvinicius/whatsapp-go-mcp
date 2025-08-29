package mcp

import (
	"context"
	"errors"
	"fmt"

	domainSend "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/send"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SendHandler struct {
	sendService domainSend.ISendUsecase
}

func InitMcpSend(sendService domainSend.ISendUsecase) *SendHandler {
	return &SendHandler{
		sendService: sendService,
	}
}

func (s *SendHandler) AddSendTools(mcpServer *server.MCPServer) {
	// Text and basic media
	mcpServer.AddTool(s.toolSendText(), s.handleSendText)
	mcpServer.AddTool(s.toolSendImage(), s.handleSendImage)
	
	// Advanced multimedia
	mcpServer.AddTool(s.toolSendAudio(), s.handleSendAudio)
	mcpServer.AddTool(s.toolSendVideo(), s.handleSendVideo)
	mcpServer.AddTool(s.toolSendFile(), s.handleSendFile)
	
	// Interactions
	mcpServer.AddTool(s.toolSendContact(), s.handleSendContact)
	mcpServer.AddTool(s.toolSendLink(), s.handleSendLink)
	mcpServer.AddTool(s.toolSendLocation(), s.handleSendLocation)
	mcpServer.AddTool(s.toolSendPoll(), s.handleSendPoll)
	
	// Presence
	mcpServer.AddTool(s.toolSendPresence(), s.handleSendPresence)
}

func (s *SendHandler) toolSendText() mcp.Tool {
	sendTextTool := mcp.NewTool("whatsapp_send_text",
		mcp.WithDescription("Send a text message to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send message to"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The text message to send"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
		mcp.WithString("reply_message_id",
			mcp.Description("Message ID to reply to (optional)"),
		),
	)

	return sendTextTool
}

func (s *SendHandler) handleSendText(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	message, ok := request.GetArguments()["message"].(string)
	if !ok {
		return nil, errors.New("message must be a string")
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	replyMessageId, ok := request.GetArguments()["reply_message_id"].(string)
	if !ok {
		replyMessageId = ""
	}

	res, err := s.sendService.SendText(ctx, domainSend.MessageRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		Message:        message,
		ReplyMessageID: &replyMessageId,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message sent successfully with ID %s", res.MessageID)), nil
}

func (s *SendHandler) toolSendContact() mcp.Tool {
	sendContactTool := mcp.NewTool("whatsapp_send_contact",
		mcp.WithDescription("Send a contact card to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send contact to"),
		),
		mcp.WithString("contact_name",
			mcp.Required(),
			mcp.Description("Name of the contact to send"),
		),
		mcp.WithString("contact_phone",
			mcp.Required(),
			mcp.Description("Phone number of the contact to send"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)

	return sendContactTool
}

func (s *SendHandler) handleSendContact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	contactName, ok := request.GetArguments()["contact_name"].(string)
	if !ok {
		return nil, errors.New("contact_name must be a string")
	}

	contactPhone, ok := request.GetArguments()["contact_phone"].(string)
	if !ok {
		return nil, errors.New("contact_phone must be a string")
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	res, err := s.sendService.SendContact(ctx, domainSend.ContactRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		ContactName:  contactName,
		ContactPhone: contactPhone,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Contact sent successfully with ID %s", res.MessageID)), nil
}

func (s *SendHandler) toolSendLink() mcp.Tool {
	sendLinkTool := mcp.NewTool("whatsapp_send_link",
		mcp.WithDescription("Send a link with caption to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send link to"),
		),
		mcp.WithString("link",
			mcp.Required(),
			mcp.Description("URL link to send"),
		),
		mcp.WithString("caption",
			mcp.Required(),
			mcp.Description("Caption or description for the link"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)

	return sendLinkTool
}

func (s *SendHandler) handleSendLink(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	link, ok := request.GetArguments()["link"].(string)
	if !ok {
		return nil, errors.New("link must be a string")
	}

	caption, ok := request.GetArguments()["caption"].(string)
	if !ok {
		caption = ""
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	res, err := s.sendService.SendLink(ctx, domainSend.LinkRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		Link:    link,
		Caption: caption,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Link sent successfully with ID %s", res.MessageID)), nil
}

func (s *SendHandler) toolSendLocation() mcp.Tool {
	sendLocationTool := mcp.NewTool("whatsapp_send_location",
		mcp.WithDescription("Send a location coordinates to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send location to"),
		),
		mcp.WithString("latitude",
			mcp.Required(),
			mcp.Description("Latitude coordinate (as string)"),
		),
		mcp.WithString("longitude",
			mcp.Required(),
			mcp.Description("Longitude coordinate (as string)"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)

	return sendLocationTool
}

func (s *SendHandler) handleSendLocation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	latitude, ok := request.GetArguments()["latitude"].(string)
	if !ok {
		return nil, errors.New("latitude must be a string")
	}

	longitude, ok := request.GetArguments()["longitude"].(string)
	if !ok {
		return nil, errors.New("longitude must be a string")
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	res, err := s.sendService.SendLocation(ctx, domainSend.LocationRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		Latitude:  latitude,
		Longitude: longitude,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Location sent successfully with ID %s", res.MessageID)), nil
}

func (s *SendHandler) toolSendImage() mcp.Tool {
	sendImageTool := mcp.NewTool("whatsapp_send_image",
		mcp.WithDescription("Send an image to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send image to"),
		),
		mcp.WithString("image_url",
			mcp.Description("URL of the image to send"),
		),
		mcp.WithString("caption",
			mcp.Description("Caption or description for the image"),
		),
		mcp.WithBoolean("view_once",
			mcp.Description("Whether this image should be viewed only once (default: false)"),
		),
		mcp.WithBoolean("compress",
			mcp.Description("Whether to compress the image (default: true)"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)

	return sendImageTool
}

func (s *SendHandler) handleSendImage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	imageURL, imageURLOk := request.GetArguments()["image_url"].(string)
	if !imageURLOk {
		return nil, errors.New("image_url must be a string")
	}

	caption, ok := request.GetArguments()["caption"].(string)
	if !ok {
		caption = ""
	}

	viewOnce, ok := request.GetArguments()["view_once"].(bool)
	if !ok {
		viewOnce = false
	}

	compress, ok := request.GetArguments()["compress"].(bool)
	if !ok {
		compress = true
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	// Create image request
	imageRequest := domainSend.ImageRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		Caption:  caption,
		ViewOnce: viewOnce,
		Compress: compress,
	}

	if imageURLOk && imageURL != "" {
		imageRequest.ImageURL = &imageURL
	}
	res, err := s.sendService.SendImage(ctx, imageRequest)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Image sent successfully with ID %s", res.MessageID)), nil
}

// ===== MULTIMEDIA TOOLS =====

func (s *SendHandler) toolSendAudio() mcp.Tool {
	return mcp.NewTool("whatsapp_send_audio",
		mcp.WithDescription("Send an audio file to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send audio to"),
		),
		mcp.WithString("audio_url",
			mcp.Required(),
			mcp.Description("URL of the audio file to send"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)
}

func (s *SendHandler) handleSendAudio(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	audioURL, ok := request.GetArguments()["audio_url"].(string)
	if !ok {
		return nil, errors.New("audio_url must be a string")
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	res, err := s.sendService.SendAudio(ctx, domainSend.AudioRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		AudioURL: &audioURL,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Audio sent successfully with ID %s", res.MessageID)), nil
}

func (s *SendHandler) toolSendVideo() mcp.Tool {
	return mcp.NewTool("whatsapp_send_video",
		mcp.WithDescription("Send a video file to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send video to"),
		),
		mcp.WithString("video_url",
			mcp.Required(),
			mcp.Description("URL of the video file to send"),
		),
		mcp.WithString("caption",
			mcp.Description("Caption or description for the video"),
		),
		mcp.WithBoolean("view_once",
			mcp.Description("Whether this video should be viewed only once (default: false)"),
		),
		mcp.WithBoolean("compress",
			mcp.Description("Whether to compress the video (default: true)"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)
}

func (s *SendHandler) handleSendVideo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	videoURL, ok := request.GetArguments()["video_url"].(string)
	if !ok {
		return nil, errors.New("video_url must be a string")
	}

	caption, ok := request.GetArguments()["caption"].(string)
	if !ok {
		caption = ""
	}

	viewOnce, ok := request.GetArguments()["view_once"].(bool)
	if !ok {
		viewOnce = false
	}

	compress, ok := request.GetArguments()["compress"].(bool)
	if !ok {
		compress = true
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	res, err := s.sendService.SendVideo(ctx, domainSend.VideoRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		Caption:  caption,
		ViewOnce: viewOnce,
		Compress: compress,
		VideoURL: &videoURL,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Video sent successfully with ID %s", res.MessageID)), nil
}

func (s *SendHandler) toolSendFile() mcp.Tool {
	return mcp.NewTool("whatsapp_send_file",
		mcp.WithDescription("Send a file/document to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send file to"),
		),
		mcp.WithString("file_url",
			mcp.Required(),
			mcp.Description("URL of the file to send"),
		),
		mcp.WithString("caption",
			mcp.Description("Caption or description for the file"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)
}

func (s *SendHandler) handleSendFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Note: File sending via URL is not directly supported by the FileRequest struct
	// which expects multipart.FileHeader. For MCP integration with AI agents,
	// we'll return a helpful message about this limitation.
	return mcp.NewToolResultText("File sending via URL requires file upload capability not available in MCP context. Use image, audio, or video tools for media files with URLs."), nil
}

func (s *SendHandler) toolSendPoll() mcp.Tool {
	return mcp.NewTool("whatsapp_send_poll",
		mcp.WithDescription("Send a poll to a WhatsApp contact or group."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number or group ID to send poll to"),
		),
		mcp.WithString("question",
			mcp.Required(),
			mcp.Description("Poll question"),
		),
		mcp.WithArray("options",
			mcp.Required(),
			mcp.Description("Array of poll options (strings)"),
		),
		mcp.WithNumber("max_answer",
			mcp.Description("Maximum number of answers allowed (default: 1)"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this message is being forwarded (default: false)"),
		),
	)
}

func (s *SendHandler) handleSendPoll(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone, ok := request.GetArguments()["phone"].(string)
	if !ok {
		return nil, errors.New("phone must be a string")
	}

	question, ok := request.GetArguments()["question"].(string)
	if !ok {
		return nil, errors.New("question must be a string")
	}

	optionsRaw, ok := request.GetArguments()["options"].([]interface{})
	if !ok {
		return nil, errors.New("options must be an array")
	}

	options := make([]string, len(optionsRaw))
	for i, opt := range optionsRaw {
		if optStr, ok := opt.(string); ok {
			options[i] = optStr
		} else {
			return nil, fmt.Errorf("option at index %d must be a string", i)
		}
	}

	maxAnswer, ok := request.GetArguments()["max_answer"].(float64)
	if !ok {
		maxAnswer = 1
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	res, err := s.sendService.SendPoll(ctx, domainSend.PollRequest{
		BaseRequest: domainSend.BaseRequest{
			Phone:       phone,
			IsForwarded: isForwarded,
		},
		Question:  question,
		Options:   options,
		MaxAnswer: int(maxAnswer),
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Poll sent successfully with ID %s", res.MessageID)), nil
}

func (s *SendHandler) toolSendPresence() mcp.Tool {
	return mcp.NewTool("whatsapp_send_presence",
		mcp.WithDescription("Send typing indicator or online presence to WhatsApp."),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Presence type: 'typing', 'recording', 'online', 'offline'"),
		),
		mcp.WithBoolean("is_forwarded",
			mcp.Description("Whether this is forwarded (default: false)"),
		),
	)
}

func (s *SendHandler) handleSendPresence(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	presenceType, ok := request.GetArguments()["type"].(string)
	if !ok {
		return nil, errors.New("type must be a string")
	}

	isForwarded, ok := request.GetArguments()["is_forwarded"].(bool)
	if !ok {
		isForwarded = false
	}

	res, err := s.sendService.SendPresence(ctx, domainSend.PresenceRequest{
		Type:        presenceType,
		IsForwarded: isForwarded,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Presence '%s' sent successfully with ID %s", presenceType, res.MessageID)), nil
}
