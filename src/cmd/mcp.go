package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/aldinokemal/go-whatsapp-web-multidevice/config"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/ui/mcp"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/ui/rest/helpers"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start WhatsApp MCP server using SSE/HTTP streaming",
	Long:  `Start a WhatsApp MCP (Model Context Protocol) server using Server-Sent Events (SSE) transport for HTTP streaming. This allows AI agents and Smithery.ai to interact with WhatsApp through a standardized protocol.`,
	Run:   mcpServer,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
	mcpCmd.Flags().StringVar(&config.McpPort, "port", "8081", "Port for the MCP server")
	mcpCmd.Flags().StringVar(&config.McpHost, "host", "0.0.0.0", "Host for the MCP server")
}

func mcpServer(_ *cobra.Command, _ []string) {
	// Set auto reconnect to whatsapp server after booting
	go helpers.SetAutoConnectAfterBooting(appUsecase)
	// Set auto reconnect checking
	go helpers.SetAutoReconnectChecking(whatsappCli)

	// Create MCP server with capabilities
	mcpServer := server.NewMCPServer(
		"WhatsApp Web Multidevice MCP Server",
		config.AppVersion,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// Add all WhatsApp tools
	// App tools (QR login, devices, etc.)
	appHandler := mcp.InitMcpApp(appUsecase)
	appHandler.AddAppTools(mcpServer)
	
	// Send tools (messages, media, etc.)
	sendHandler := mcp.InitMcpSend(sendUsecase)
	sendHandler.AddSendTools(mcpServer)
	
	// User tools (info, avatar, privacy)
	userHandler := mcp.InitMcpUser(userUsecase)
	userHandler.AddUserTools(mcpServer)
	
	// Message tools (react, delete, mark as read)
	messageHandler := mcp.InitMcpMessage(messageUsecase)
	messageHandler.AddMessageTools(mcpServer)
	
	// Group tools (create, manage, participants)
	groupHandler := mcp.InitMcpGroup(groupUsecase)
	groupHandler.AddGroupTools(mcpServer)
	
	// Chat tools (list, archive, delete)
	chatHandler := mcp.InitMcpChat(chatUsecase)
	chatHandler.AddChatTools(mcpServer)
	
	// Newsletter tools (unfollow)
	newsletterHandler := mcp.InitMcpNewsletter(newsletterUsecase)
	newsletterHandler.AddNewsletterTools(mcpServer)

	// Get port from environment variable (Smithery sets this to 8081)
	port := os.Getenv("PORT")
	if port == "" {
		port = config.McpPort
	}

	// Create Streamable HTTP server for Smithery.ai compatibility
	// Use stateless mode for simpler integration with Smithery
	streamableServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithEndpointPath("/mcp"),
		server.WithStateLess(true), // Enable stateless mode for Smithery
	)

	// Create HTTP server with CORS and session middleware
	mux := http.NewServeMux()
	mux.Handle("/mcp", corsMiddleware(sessionMiddleware(streamableServer)))
	
	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"whatsapp-mcp"}`))
	})
	
	// Add tools info endpoint for debugging
	mux.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		tools := `{
			"total": 50,
			"note": "Complete API coverage with all advanced features implemented for MCP AI agents",
			"categories": {
				"app": ["whatsapp_get_qr", "whatsapp_login_with_code", "whatsapp_logout", "whatsapp_reconnect", "whatsapp_get_devices"],
				"send": ["whatsapp_send_text", "whatsapp_send_image", "whatsapp_send_audio", "whatsapp_send_video", "whatsapp_send_file", "whatsapp_send_contact", "whatsapp_send_link", "whatsapp_send_location", "whatsapp_send_poll", "whatsapp_send_presence"],
				"message": ["whatsapp_get_messages", "whatsapp_mark_as_read", "whatsapp_react_message", "whatsapp_delete_message", "whatsapp_update_message", "whatsapp_revoke_message", "whatsapp_star_message", "whatsapp_unstar_message", "whatsapp_download_media"],
				"group": ["whatsapp_create_group", "whatsapp_leave_group", "whatsapp_get_group_info", "whatsapp_join_group_link", "whatsapp_get_invite_link", "whatsapp_set_group_name", "whatsapp_set_group_locked", "whatsapp_set_group_announce", "whatsapp_set_group_topic", "whatsapp_add_group_participants", "whatsapp_remove_group_participants", "whatsapp_promote_group_admin", "whatsapp_demote_group_admin", "whatsapp_get_group_info_from_link", "whatsapp_get_group_request_participants", "whatsapp_manage_group_request_participants"],
				"user": ["whatsapp_get_user_info", "whatsapp_check_phone", "whatsapp_get_business_profile", "whatsapp_get_avatar", "whatsapp_change_avatar", "whatsapp_change_push_name", "whatsapp_get_my_groups", "whatsapp_get_my_newsletters", "whatsapp_get_my_contacts", "whatsapp_get_my_privacy"],
				"chat": ["whatsapp_get_chat_list", "whatsapp_archive_chat", "whatsapp_mark_chat_as_read", "whatsapp_delete_chat"],
				"newsletter": ["whatsapp_unfollow_newsletter"]
			}
		}`
		w.Write([]byte(tools))
	})

	// Start the HTTP server with CORS support
	addr := fmt.Sprintf("%s:%s", config.McpHost, port)
	logrus.Printf("Starting WhatsApp MCP Streamable HTTP server on %s", addr)
	logrus.Printf("MCP endpoint: http://%s/mcp", addr)
	logrus.Printf("Health endpoint: http://%s/health", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		logrus.Fatalf("Failed to start HTTP server: %v", err)
	}
}
