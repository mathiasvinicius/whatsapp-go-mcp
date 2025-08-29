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
	// User information
	mcpServer.AddTool(u.toolGetInfo(), u.handleGetInfo)
	mcpServer.AddTool(u.toolCheckPhone(), u.handleCheckPhone)
	mcpServer.AddTool(u.toolBusinessProfile(), u.handleBusinessProfile)
	
	// Profile management
	mcpServer.AddTool(u.toolGetAvatar(), u.handleGetAvatar)
	mcpServer.AddTool(u.toolChangeAvatar(), u.handleChangeAvatar)
	mcpServer.AddTool(u.toolChangePushName(), u.handleChangePushName)
	
	// Listings
	mcpServer.AddTool(u.toolMyListGroups(), u.handleMyListGroups)
	mcpServer.AddTool(u.toolMyListNewsletter(), u.handleMyListNewsletter)
	mcpServer.AddTool(u.toolMyListContacts(), u.handleMyListContacts)
	
	// Privacy
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

// ===== BUSINESS PROFILE TOOLS =====

func (u *UserHandler) toolBusinessProfile() mcp.Tool {
	return mcp.NewTool("whatsapp_get_business_profile",
		mcp.WithDescription("Get business profile information for a WhatsApp Business account."),
		mcp.WithString("phone",
			mcp.Required(),
			mcp.Description("Phone number of the business account"),
		),
	)
}

func (u *UserHandler) handleBusinessProfile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phone := request.GetArguments()["phone"].(string)

	response, err := u.userService.BusinessProfile(ctx, domainUser.BusinessProfileRequest{
		Phone: phone,
	})
	
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Business Profile:\nJID: %s\nEmail: %s\nAddress: %s\nTimezone: %s", 
		response.JID, response.Email, response.Address, response.BusinessHoursTimeZone)
	
	if len(response.Categories) > 0 {
		result += "\nCategories:\n"
		for _, cat := range response.Categories {
			result += fmt.Sprintf("- %s (ID: %s)\n", cat.Name, cat.ID)
		}
	}
	
	if len(response.BusinessHours) > 0 {
		result += "\nBusiness Hours:\n"
		for _, hours := range response.BusinessHours {
			result += fmt.Sprintf("- %s: %s %s-%s\n", hours.DayOfWeek, hours.Mode, hours.OpenTime, hours.CloseTime)
		}
	}
	
	return mcp.NewToolResultText(result), nil
}

// ===== PROFILE MANAGEMENT TOOLS =====

func (u *UserHandler) toolChangeAvatar() mcp.Tool {
	return mcp.NewTool("whatsapp_change_avatar",
		mcp.WithDescription("Change the profile avatar/picture of the logged-in WhatsApp account. Note: File upload in MCP context has limitations."),
		mcp.WithString("avatar_path",
			mcp.Required(),
			mcp.Description("Local file path to the new avatar image (JPG/PNG)"),
		),
	)
}

func (u *UserHandler) handleChangeAvatar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	avatarPath := request.GetArguments()["avatar_path"].(string)

	// Note: This implementation has limitations in MCP context due to multipart.FileHeader requirement
	// The actual file upload would need to be handled differently in production
	return mcp.NewToolResultText(fmt.Sprintf("Avatar change requested for file: %s\nNote: File upload functionality requires direct REST API access for multipart form handling", avatarPath)), nil
}

func (u *UserHandler) toolChangePushName() mcp.Tool {
	return mcp.NewTool("whatsapp_change_push_name",
		mcp.WithDescription("Change the display name (push name) of the logged-in WhatsApp account."),
		mcp.WithString("push_name",
			mcp.Required(),
			mcp.Description("New display name for the account"),
		),
	)
}

func (u *UserHandler) handleChangePushName(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pushName := request.GetArguments()["push_name"].(string)

	err := u.userService.ChangePushName(ctx, domainUser.ChangePushNameRequest{
		PushName: pushName,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Display name changed to: %s", pushName)), nil
}

// ===== LISTING TOOLS =====

func (u *UserHandler) toolMyListNewsletter() mcp.Tool {
	return mcp.NewTool("whatsapp_get_my_newsletters",
		mcp.WithDescription("Get list of newsletters the logged-in account has subscribed to."),
	)
}

func (u *UserHandler) handleMyListNewsletter(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	response, err := u.userService.MyListNewsletter(ctx)
	
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return mcp.NewToolResultText("No newsletters subscribed"), nil
	}

	result := fmt.Sprintf("My Newsletters (%d):\n", len(response.Data))
	for _, newsletter := range response.Data {
		result += fmt.Sprintf("- ID: %s\n", newsletter.ID.String())
		result += fmt.Sprintf("  State: %d\n", newsletter.State)
		result += fmt.Sprintf("  Created: %s\n", newsletter.ThreadMeta.CreationTime.Format("2006-01-02 15:04:05"))
	}
	return mcp.NewToolResultText(result), nil
}

func (u *UserHandler) toolMyListContacts() mcp.Tool {
	return mcp.NewTool("whatsapp_get_my_contacts",
		mcp.WithDescription("Get list of contacts in the logged-in account's address book."),
	)
}

func (u *UserHandler) handleMyListContacts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	response, err := u.userService.MyListContacts(ctx)
	
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return mcp.NewToolResultText("No contacts found"), nil
	}

	result := fmt.Sprintf("My Contacts (%d):\n", len(response.Data))
	for _, contact := range response.Data {
		result += fmt.Sprintf("- %s (JID: %s)\n", contact.Name, contact.JID.String())
	}
	return mcp.NewToolResultText(result), nil
}