package mcp

import (
	"context"
	"fmt"

	domainUser "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/user"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type UserHandler struct {
	userService domainUser.IUserUsecase
}

func InitMcpUser(userService domainUser.IUserUsecase) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) AddUserTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(u.toolGetInfo(), u.handleGetInfo)
	mcpServer.AddTool(u.toolGetAvatar(), u.handleGetAvatar)
	mcpServer.AddTool(u.toolMyListGroups(), u.handleMyListGroups)
	mcpServer.AddTool(u.toolCheckPhone(), u.handleCheckPhone)
	mcpServer.AddTool(u.toolGetMyPrivacy(), u.handleGetMyPrivacy)
}

func (u *UserHandler) toolGetInfo() mcp.Tool {
	return mcp.NewTool("whatsapp_get_user_info",
		mcp.WithDescription("Get information about a WhatsApp user."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number with country code"),
		),
	)
}

func (u *UserHandler) handleGetInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)

	response, err := u.userService.Info(ctx, domainUser.InfoRequest{
		Phone: phone,
	})
	
	if err != nil {
		return nil, err
	}

	result := "User Info:\n"
	for _, data := range response.Data {
		result += fmt.Sprintf("Status: %s\nPicture ID: %s\nVerified Name: %s\n", 
			data.Status, data.PictureID, data.VerifiedName)
		if len(data.Devices) > 0 {
			result += fmt.Sprintf("Devices: %d\n", len(data.Devices))
		}
	}
	return mcp.NewToolResultText(result), nil
}

func (u *UserHandler) toolGetAvatar() mcp.Tool {
	return mcp.NewTool("whatsapp_get_avatar",
		mcp.WithDescription("Get user's WhatsApp profile picture."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number with country code"),
		),
		mcp.WithBoolean("preview",
			mcp.Description("Get preview size (default: false for full size)"),
		),
	)
}

func (u *UserHandler) handleGetAvatar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)
	preview := false
	if p, ok := request.GetArguments()["preview"].(bool); ok {
		preview = p
	}

	response, err := u.userService.Avatar(ctx, domainUser.AvatarRequest{
		Phone:      phone,
		IsPreview:  preview,
		IsCommunity: false,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Avatar URL: %s\nID: %s", response.URL, response.ID)), nil
}

func (u *UserHandler) toolMyListGroups() mcp.Tool {
	return mcp.NewTool("whatsapp_get_my_groups",
		mcp.WithDescription("Get list of groups the logged-in account has joined."),
	)
}

func (u *UserHandler) handleMyListGroups(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	response, err := u.userService.MyListGroups(ctx)
	
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("My Groups (%d):\n", len(response.Data))
	for _, group := range response.Data {
		result += fmt.Sprintf("- %s (ID: %s)\n", group.GroupName.Name, group.JID.String())
	}
	return mcp.NewToolResultText(result), nil
}

func (u *UserHandler) toolCheckPhone() mcp.Tool {
	return mcp.NewTool("whatsapp_check_phone",
		mcp.WithDescription("Check if a phone number is registered on WhatsApp."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number to check"),
		),
	)
}

func (u *UserHandler) handleCheckPhone(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)

	response, err := u.userService.IsOnWhatsApp(ctx, domainUser.CheckRequest{
		Phone: phone,
	})
	
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Phone %s is on WhatsApp: %v", phone, response.IsOnWhatsApp)
	return mcp.NewToolResultText(result), nil
}

func (u *UserHandler) toolGetMyPrivacy() mcp.Tool {
	return mcp.NewTool("whatsapp_get_my_privacy",
		mcp.WithDescription("Get privacy settings of the logged-in WhatsApp account."),
	)
}

func (u *UserHandler) handleGetMyPrivacy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	response, err := u.userService.MyPrivacySetting(ctx)
	
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Privacy Settings:\nLast Seen: %s\nProfile: %s\nStatus: %s\nRead Receipts: %s\nGroups: %s", 
		response.LastSeen, response.Profile, response.Status, response.ReadReceipts, response.GroupAdd)
	return mcp.NewToolResultText(result), nil
}