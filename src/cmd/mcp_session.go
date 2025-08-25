package cmd

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// sessionConfig handles optional configuration from Smithery
type sessionConfig struct {
	// Add any configuration fields if needed
	// For now, we don't require any specific config
}

// parseSessionConfig extracts and decodes the base64-encoded config from URL
func parseSessionConfig(r *http.Request) (*sessionConfig, error) {
	// Parse URL to get query parameters
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return nil, err
	}
	
	// Get config parameter
	configParam := u.Query().Get("config")
	if configParam == "" {
		// No config provided, return empty config
		logrus.Debug("No config parameter provided")
		return &sessionConfig{}, nil
	}
	
	// Decode base64
	configData, err := base64.StdEncoding.DecodeString(configParam)
	if err != nil {
		// Try URL-safe base64 encoding
		configData, err = base64.URLEncoding.DecodeString(configParam)
		if err != nil {
			logrus.Warnf("Failed to decode config parameter: %v", err)
			return &sessionConfig{}, nil // Return empty config on error
		}
	}
	
	// Parse JSON
	var config sessionConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		logrus.Warnf("Failed to parse config JSON: %v", err)
		return &sessionConfig{}, nil // Return empty config on error
	}
	
	logrus.Debugf("Parsed session config: %+v", config)
	return &config, nil
}

// sessionMiddleware handles session configuration from Smithery
func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse session config if present
		config, err := parseSessionConfig(r)
		if err != nil {
			logrus.Errorf("Error parsing session config: %v", err)
		} else if config != nil {
			// Store config in request context if needed
			// For now, we just log it
			logrus.Debugf("Session config received: %+v", config)
		}
		
		// Strip config parameter from URL before passing to MCP handler
		if strings.Contains(r.URL.RawQuery, "config=") {
			u, _ := url.Parse(r.URL.String())
			q := u.Query()
			q.Del("config")
			u.RawQuery = q.Encode()
			r.URL = u
		}
		
		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}