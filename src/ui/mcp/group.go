package mcp

import (
	"context"
	"fmt"

	domainGroup "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/group"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.mau.fi/whatsmeow"
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
	// Basic group operations
	mcpServer.AddTool(g.toolCreateGroup(), g.handleCreateGroup)
	mcpServer.AddTool(g.toolLeaveGroup(), g.handleLeaveGroup)
	mcpServer.AddTool(g.toolGetGroupInfo(), g.handleGetGroupInfo)
	mcpServer.AddTool(g.toolJoinWithLink(), g.handleJoinWithLink)
	mcpServer.AddTool(g.toolGetInviteLink(), g.handleGetInviteLink)
	
	// Group settings
	mcpServer.AddTool(g.toolSetGroupName(), g.handleSetGroupName)
	mcpServer.AddTool(g.toolSetGroupLocked(), g.handleSetGroupLocked)
	mcpServer.AddTool(g.toolSetGroupAnnounce(), g.handleSetGroupAnnounce)
	mcpServer.AddTool(g.toolSetGroupTopic(), g.handleSetGroupTopic)
	
	// Participant management
	mcpServer.AddTool(g.toolAddParticipants(), g.handleAddParticipants)
	mcpServer.AddTool(g.toolRemoveParticipants(), g.handleRemoveParticipants)
	mcpServer.AddTool(g.toolPromoteAdmin(), g.handlePromoteAdmin)
	mcpServer.AddTool(g.toolDemoteAdmin(), g.handleDemoteAdmin)
	
	// Advanced group features
	mcpServer.AddTool(g.toolGetGroupInfoFromLink(), g.handleGetGroupInfoFromLink)
	mcpServer.AddTool(g.toolGetGroupRequestParticipants(), g.handleGetGroupRequestParticipants)
	mcpServer.AddTool(g.toolManageGroupRequestParticipants(), g.handleManageGroupRequestParticipants)
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

// ===== GROUP SETTINGS TOOLS =====

func (g *GroupHandler) toolSetGroupTopic() mcp.Tool {
	return mcp.NewTool("whatsapp_set_group_topic",
		mcp.WithDescription("Set the topic/description of a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithString("topic",
			mcp.Required(),
			mcp.Description("New topic/description for the group"),
		),
	)
}

func (g *GroupHandler) handleSetGroupTopic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	topic := request.GetArguments()["topic"].(string)

	err := g.groupService.SetGroupTopic(ctx, domainGroup.SetGroupTopicRequest{
		GroupID: groupID,
		Topic:   topic,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Group topic updated successfully")), nil
}

// ===== PARTICIPANT MANAGEMENT TOOLS =====

func (g *GroupHandler) toolAddParticipants() mcp.Tool {
	return mcp.NewTool("whatsapp_add_group_participants",
		mcp.WithDescription("Add participants to a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithArray("participants",
			mcp.Required(),
			mcp.Description("Array of phone numbers to add to the group"),
		),
	)
}

func (g *GroupHandler) handleAddParticipants(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	participantsRaw := request.GetArguments()["participants"].([]interface{})
	
	participants := make([]string, len(participantsRaw))
	for i, p := range participantsRaw {
		participants[i] = p.(string)
	}

	result, err := g.groupService.ManageParticipant(ctx, domainGroup.ParticipantRequest{
		GroupID:      groupID,
		Participants: participants,
		Action:       whatsmeow.ParticipantChangeAdd,
	})
	
	if err != nil {
		return nil, err
	}

	// Build response with results for each participant
	response := "Add participants results:\n"
	for _, r := range result {
		response += fmt.Sprintf("- %s: %s (%s)\n", r.Participant, r.Status, r.Message)
	}
	
	return mcp.NewToolResultText(response), nil
}

func (g *GroupHandler) toolRemoveParticipants() mcp.Tool {
	return mcp.NewTool("whatsapp_remove_group_participants",
		mcp.WithDescription("Remove participants from a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithArray("participants",
			mcp.Required(),
			mcp.Description("Array of phone numbers to remove from the group"),
		),
	)
}

func (g *GroupHandler) handleRemoveParticipants(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	participantsRaw := request.GetArguments()["participants"].([]interface{})
	
	participants := make([]string, len(participantsRaw))
	for i, p := range participantsRaw {
		participants[i] = p.(string)
	}

	result, err := g.groupService.ManageParticipant(ctx, domainGroup.ParticipantRequest{
		GroupID:      groupID,
		Participants: participants,
		Action:       whatsmeow.ParticipantChangeRemove,
	})
	
	if err != nil {
		return nil, err
	}

	// Build response with results for each participant
	response := "Remove participants results:\n"
	for _, r := range result {
		response += fmt.Sprintf("- %s: %s (%s)\n", r.Participant, r.Status, r.Message)
	}
	
	return mcp.NewToolResultText(response), nil
}

func (g *GroupHandler) toolPromoteAdmin() mcp.Tool {
	return mcp.NewTool("whatsapp_promote_group_admin",
		mcp.WithDescription("Promote participants to admin in a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithArray("participants",
			mcp.Required(),
			mcp.Description("Array of phone numbers to promote to admin"),
		),
	)
}

func (g *GroupHandler) handlePromoteAdmin(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	participantsRaw := request.GetArguments()["participants"].([]interface{})
	
	participants := make([]string, len(participantsRaw))
	for i, p := range participantsRaw {
		participants[i] = p.(string)
	}

	result, err := g.groupService.ManageParticipant(ctx, domainGroup.ParticipantRequest{
		GroupID:      groupID,
		Participants: participants,
		Action:       whatsmeow.ParticipantChangePromote,
	})
	
	if err != nil {
		return nil, err
	}

	// Build response with results for each participant
	response := "Promote to admin results:\n"
	for _, r := range result {
		response += fmt.Sprintf("- %s: %s (%s)\n", r.Participant, r.Status, r.Message)
	}
	
	return mcp.NewToolResultText(response), nil
}

func (g *GroupHandler) toolDemoteAdmin() mcp.Tool {
	return mcp.NewTool("whatsapp_demote_group_admin",
		mcp.WithDescription("Demote admins to regular participants in a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithArray("participants",
			mcp.Required(),
			mcp.Description("Array of phone numbers to demote from admin"),
		),
	)
}

func (g *GroupHandler) handleDemoteAdmin(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	participantsRaw := request.GetArguments()["participants"].([]interface{})
	
	participants := make([]string, len(participantsRaw))
	for i, p := range participantsRaw {
		participants[i] = p.(string)
	}

	result, err := g.groupService.ManageParticipant(ctx, domainGroup.ParticipantRequest{
		GroupID:      groupID,
		Participants: participants,
		Action:       whatsmeow.ParticipantChangeDemote,
	})
	
	if err != nil {
		return nil, err
	}

	// Build response with results for each participant
	response := "Demote from admin results:\n"
	for _, r := range result {
		response += fmt.Sprintf("- %s: %s (%s)\n", r.Participant, r.Status, r.Message)
	}
	
	return mcp.NewToolResultText(response), nil
}

// ===== ADVANCED GROUP FEATURES =====

func (g *GroupHandler) toolGetGroupInfoFromLink() mcp.Tool {
	return mcp.NewTool("whatsapp_get_group_info_from_link",
		mcp.WithDescription("Get information about a WhatsApp group from an invite link."),
		mcp.WithString("link",
			mcp.Required(),
			mcp.Description("WhatsApp group invite link"),
		),
	)
}

func (g *GroupHandler) handleGetGroupInfoFromLink(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	link := request.GetArguments()["link"].(string)

	response, err := g.groupService.GetGroupInfoFromLink(ctx, domainGroup.GetGroupInfoFromLinkRequest{
		Link: link,
	})
	
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Group Info from Link:\nName: %s\nTopic: %s\nParticipants: %d\nGroup ID: %s\nCreated: %s\nLocked: %t\nAnnounce Mode: %t\nEphemeral: %t",
		response.Name, response.Topic, response.ParticipantCount, response.GroupID, 
		response.CreatedAt.Format("2006-01-02 15:04:05"), response.IsLocked, response.IsAnnounce, response.IsEphemeral)
	
	if response.Description != "" {
		result += fmt.Sprintf("\nDescription: %s", response.Description)
	}
	
	return mcp.NewToolResultText(result), nil
}

func (g *GroupHandler) toolGetGroupRequestParticipants() mcp.Tool {
	return mcp.NewTool("whatsapp_get_group_request_participants",
		mcp.WithDescription("Get list of participants requesting to join a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
	)
}

func (g *GroupHandler) handleGetGroupRequestParticipants(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)

	result, err := g.groupService.GetGroupRequestParticipants(ctx, domainGroup.GetGroupRequestParticipantsRequest{
		GroupID: groupID,
	})
	
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return mcp.NewToolResultText("No pending join requests for this group"), nil
	}

	response := "Pending join requests:\n"
	for _, r := range result {
		response += fmt.Sprintf("- %s (requested at: %s)\n", r.JID, r.RequestedAt.Format("2006-01-02 15:04:05"))
	}
	
	return mcp.NewToolResultText(response), nil
}

func (g *GroupHandler) toolManageGroupRequestParticipants() mcp.Tool {
	return mcp.NewTool("whatsapp_manage_group_request_participants",
		mcp.WithDescription("Approve or reject participants requesting to join a WhatsApp group."),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("Group ID/JID"),
		),
		mcp.WithArray("participants",
			mcp.Required(),
			mcp.Description("Array of phone numbers to approve/reject"),
		),
		mcp.WithString("action",
			mcp.Required(),
			mcp.Description("Action to take: 'approve' or 'reject'"),
		),
	)
}

func (g *GroupHandler) handleManageGroupRequestParticipants(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupID := request.GetArguments()["group_id"].(string)
	participantsRaw := request.GetArguments()["participants"].([]interface{})
	action := request.GetArguments()["action"].(string)
	
	participants := make([]string, len(participantsRaw))
	for i, p := range participantsRaw {
		participants[i] = p.(string)
	}

	// Convert string action to whatsmeow enum
	var actionEnum whatsmeow.ParticipantRequestChange
	if action == "approve" {
		actionEnum = whatsmeow.ParticipantChangeApprove
	} else if action == "reject" {
		actionEnum = whatsmeow.ParticipantChangeReject
	} else {
		return mcp.NewToolResultText(fmt.Sprintf("Invalid action: %s. Must be 'approve' or 'reject'", action)), nil
	}

	result, err := g.groupService.ManageGroupRequestParticipants(ctx, domainGroup.GroupRequestParticipantsRequest{
		GroupID:      groupID,
		Participants: participants,
		Action:       actionEnum,
	})
	
	if err != nil {
		return nil, err
	}

	// Build response with results for each participant
	response := fmt.Sprintf("Join request %s results:\n", action)
	for _, r := range result {
		response += fmt.Sprintf("- %s: %s (%s)\n", r.Participant, r.Status, r.Message)
	}
	
	return mcp.NewToolResultText(response), nil
}