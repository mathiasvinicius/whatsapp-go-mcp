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
	sendHandler := mcp.InitMcpSend(sendUsecase)
	sendHandler.AddSendTools(mcpServer)

	// Get port from environment variable (Smithery sets this to 8081)
	port := os.Getenv("PORT")
	if port == "" {
		port = config.McpPort
	}

	// Create Streamable HTTP server for Smithery.ai compatibility
	streamableServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithEndpointPath("/mcp"),
		server.WithStateLess(false),
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

	// Start the HTTP server with CORS support
	addr := fmt.Sprintf("%s:%s", config.McpHost, port)
	logrus.Printf("Starting WhatsApp MCP Streamable HTTP server on %s", addr)
	logrus.Printf("MCP endpoint: http://%s/mcp", addr)
	logrus.Printf("Health endpoint: http://%s/health", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		logrus.Fatalf("Failed to start HTTP server: %v", err)
	}
}
