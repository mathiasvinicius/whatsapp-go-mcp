package mcp

import (
	"context"
	"fmt"

	domainNewsletter "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/newsletter"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type NewsletterHandler struct {
	newsletterService domainNewsletter.INewsletterUsecase
}

func InitMcpNewsletter(newsletterService domainNewsletter.INewsletterUsecase) *NewsletterHandler {
	return &NewsletterHandler{
		newsletterService: newsletterService,
	}
}

func (n *NewsletterHandler) AddNewsletterTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(n.toolUnfollow(), n.handleUnfollow)
}

func (n *NewsletterHandler) toolUnfollow() mcp.Tool {
	return mcp.NewTool("whatsapp_unfollow_newsletter",
		mcp.WithDescription("Unfollow/unsubscribe from a WhatsApp newsletter."),
		mcp.WithString("newsletter_id",
			mcp.Required(),
			mcp.Description("Newsletter ID to unfollow"),
		),
	)
}

func (n *NewsletterHandler) handleUnfollow(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	newsletterID := request.GetArguments()["newsletter_id"].(string)

	err := n.newsletterService.Unfollow(ctx, domainNewsletter.UnfollowRequest{
		NewsletterID: newsletterID,
	})
	
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully unfollowed newsletter: %s", newsletterID)), nil
}