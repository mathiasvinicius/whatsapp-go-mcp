package mcp

import (
	"context"
	"fmt"

	domainGroup "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/group"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GroupHandler struct {
	groupService domainGroup.IGroupUsecase
}

func InitMcpGroup(groupService domainGroup.IGroupUsecase) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

func (g *GroupHandler) AddGroupTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(g.toolCreateGroup(), g.handleCreateGroup)
	mcpServer.AddTool(g.toolLeaveGroup(), g.handleLeaveGroup)
	mcpServer.AddTool(g.toolGetGroupInfo(), g.handleGetGroupInfo)
	mcpServer.AddTool(g.toolJoinWithLink(), g.handleJoinWithLink)
	mcpServer.AddTool(g.toolGetInviteLink(), g.handleGetInviteLink)
	mcpServer.AddTool(g.toolSetGroupName(), g.handleSetGroupName)
	mcpServer.AddTool(g.toolSetGroupLocked(), g.handleSetGroupLocked)
	mcpServer.AddTool(g.toolSetGroupAnnounce(), g.handleSetGroupAnnounce)
}

func (g *GroupHandler) toolCreateGroup() mcp.Tool {
	return mcp.NewTool("whatsapp_create_group",
		mcp.WithDescription("Create a new WhatsApp group."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Group name"),
		),
		mcp.WithArray("participants",
			mcp.Required(),
			mcp.Description("Array of phone numbers to add as participants"),
		),
	)
}

func (g *GroupHandler) handleCreateGroup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetArguments()["name"].(string)
	participantsRaw := request.GetArguments()["participants"].([]interface{})
	
	participants := make([]string, len(participantsRaw))
	for i, p := range participantsRaw {
		participants[i] = p.(string)
	}

	groupID, err := g.groupService.CreateGroup(ctx, domainGroup.CreateGroupRequest{
		Title:        name,
		Participants: participants,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Group created successfully!\nID: %s\nName: %s", groupID, name)), nil
}

func (g *GroupHandler) toolLeaveGroup() mcp.Tool {
	return mcp.NewTool("whatsapp_leave_group",
		mcp.WithDescription("Leave a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
	)
}

func (g *GroupHandler) handleLeaveGroup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)

	err := g.groupService.LeaveGroup(ctx, domainGroup.LeaveGroupRequest{
		GroupID: groupID,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully left group %s", groupID)), nil
}

func (g *GroupHandler) toolGetGroupInfo() mcp.Tool {
	return mcp.NewTool("whatsapp_get_group_info",
		mcp.WithDescription("Get information about a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
	)
}

func (g *GroupHandler) handleGetGroupInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)

	response, err := g.groupService.GroupInfo(ctx, domainGroup.GroupInfoRequest{
		GroupID: groupID,
	})
	
	if err != nil {
		return nil, err
	}

	// The response.Data contains the actual group info
	result := fmt.Sprintf("Group Info:\n%+v", response.Data)
	return mcp.NewToolResultText(result), nil
}

func (g *GroupHandler) toolJoinWithLink() mcp.Tool {
	return mcp.NewTool("whatsapp_join_group_link",
		mcp.WithDescription("Join a WhatsApp group using an invite link."),
		mcp.WithString("link",
			mcp.Required(),
			mcp.Description("WhatsApp group invite link"),
		),
	)
}

func (g *GroupHandler) handleJoinWithLink(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	link := request.GetArguments()["link"].(string)

	groupID, err := g.groupService.JoinGroupWithLink(ctx, domainGroup.JoinGroupWithLinkRequest{
		Link: link,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully joined group: %s", groupID)), nil
}

func (g *GroupHandler) toolGetInviteLink() mcp.Tool {
	return mcp.NewTool("whatsapp_get_invite_link",
		mcp.WithDescription("Get invite link for a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
	)
}

func (g *GroupHandler) handleGetInviteLink(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)

	response, err := g.groupService.GetGroupInviteLink(ctx, domainGroup.GetGroupInviteLinkRequest{
		GroupID: groupID,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Invite link: %s", response.InviteLink)), nil
}

func (g *GroupHandler) toolSetGroupName() mcp.Tool {
	return mcp.NewTool("whatsapp_set_group_name",
		mcp.WithDescription("Change the name of a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("New name for the group"),
		),
	)
}

func (g *GroupHandler) handleSetGroupName(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	name := request.GetArguments()["name"].(string)

	err := g.groupService.SetGroupName(ctx, domainGroup.SetGroupNameRequest{
		GroupID: groupID,
		Name:    name,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Group name changed to: %s", name)), nil
}

func (g *GroupHandler) toolSetGroupLocked() mcp.Tool {
	return mcp.NewTool("whatsapp_set_group_locked",
		mcp.WithDescription("Lock or unlock group settings (only admins can edit)."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithBoolean("locked",
			mcp.Required(),
			mcp.Description("True to lock, false to unlock"),
		),
	)
}

func (g *GroupHandler) handleSetGroupLocked(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	locked := request.GetArguments()["locked"].(bool)

	err := g.groupService.SetGroupLocked(ctx, domainGroup.SetGroupLockedRequest{
		GroupID: groupID,
		Locked:  locked,
	})
	
	if err != nil {
		return nil, err
	}

	status := "unlocked"
	if locked {
		status = "locked"
	}

	return mcp.NewToolResultText(fmt.Sprintf("Group settings %s", status)), nil
}

func (g *GroupHandler) toolSetGroupAnnounce() mcp.Tool {
	return mcp.NewTool("whatsapp_set_group_announce",
		mcp.WithDescription("Set group to announcement mode (only admins can send messages)."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithBoolean("announce",
			mcp.Required(),
			mcp.Description("True for announcement mode, false for all participants"),
		),
	)
}

func (g *GroupHandler) handleSetGroupAnnounce(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	announce := request.GetArguments()["announce"].(bool)

	err := g.groupService.SetGroupAnnounce(ctx, domainGroup.SetGroupAnnounceRequest{
		GroupID:  groupID,
		Announce: announce,
	})
	
	if err != nil {
		return nil, err
	}

	mode := "all participants can send messages"
	if announce {
		mode = "only admins can send messages"
	}

	return mcp.NewToolResultText(fmt.Sprintf("Group mode changed: %s", mode)), nil
}