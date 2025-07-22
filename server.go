package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Configuration
type Config struct {
	Port             string `json:"port"`
	SessionSecret    string `json:"session_secret"`
	SessionDuration  int    `json:"session_duration_hours"`
	MaxLoginAttempts int    `json:"max_login_attempts"`
	LockoutDuration  int    `json:"lockout_duration_minutes"`
}

// Authorized Agent
type Agent struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	AccessKey   string   `json:"access_key"`
}

// Session Data
type SessionData struct {
	AgentID     string    `json:"agent_id"`
	AgentName   string    `json:"agent_name"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	Token       string    `json:"token"`
	Expires     time.Time `json:"expires"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
}

// Login Attempt
type LoginAttempt struct {
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}

// Response structures
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Agent   *Agent `json:"agent,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Global variables
var (
	config        Config
	agents        map[string]Agent
	store         *sessions.CookieStore
	loginAttempts map[string]LoginAttempt
)

func init() {
	// Load configuration
	loadConfig()

	// Initialize session store
	store = sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.SessionDuration * 3600,
		HttpOnly: true,
		Secure:   true, // Set to false for HTTP development
		SameSite: http.SameSiteStrictMode,
	}

	// Initialize login attempts map
	loginAttempts = make(map[string]LoginAttempt)

	// Load authorized agents
	loadAgents()
}

func loadConfig() {
	config = Config{
		Port:             ":8080",
		SessionSecret:    "your-super-secret-session-key-change-this-in-production",
		SessionDuration:  24,
		MaxLoginAttempts: 5,
		LockoutDuration:  15,
	}

	// Try to load from config file
	if _, err := os.Stat("config.json"); err == nil {
		file, err := os.Open("config.json")
		if err == nil {
			defer file.Close()
			json.NewDecoder(file).Decode(&config)
		}
	}
}

func loadAgents() {
	agents = map[string]Agent{
		"sailtix-agent-001": {
			ID:          "sailtix-agent-001",
			Name:        "Primary Sailtix Agent",
			Role:        "production",
			Permissions: []string{"read", "test"},
			AccessKey:   "access-key-2024-sailtix-secure",
		},
		"sailtix-agent-002": {
			ID:          "sailtix-agent-002",
			Name:        "Secondary Sailtix Agent",
			Role:        "production",
			Permissions: []string{"read", "test"},
			AccessKey:   "access-key-2024-sailtix-secure",
		},
		"sailtix-dev-001": {
			ID:          "sailtix-dev-001",
			Name:        "Development Agent",
			Role:        "development",
			Permissions: []string{"read", "test", "debug"},
			AccessKey:   "dev-access-key-2024",
		},
		"sailtix-admin": {
			ID:          "sailtix-admin",
			Name:        "Administrator",
			Role:        "admin",
			Permissions: []string{"read", "test", "debug", "admin"},
			AccessKey:   "admin-access-key-2024-secure",
		},
	}

	// Try to load from agents file
	if _, err := os.Stat("agents.json"); err == nil {
		file, err := os.Open("agents.json")
		if err == nil {
			defer file.Close()
			var agentConfig struct {
				AuthorizedAgents map[string]Agent `json:"authorized_agents"`
			}
			if json.NewDecoder(file).Decode(&agentConfig) == nil {
				agents = agentConfig.AuthorizedAgents
			}
		}
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func checkRateLimit(ip string) error {
	attempt, exists := loginAttempts[ip]
	if !exists {
		return nil
	}

	if attempt.Count >= config.MaxLoginAttempts {
		timeDiff := time.Since(attempt.Time)
		if timeDiff < time.Duration(config.LockoutDuration)*time.Minute {
			return fmt.Errorf("too many login attempts. Please try again in %d minutes",
				int((time.Duration(config.LockoutDuration)*time.Minute - timeDiff).Minutes()))
		}
		// Reset after lockout period
		delete(loginAttempts, ip)
	}
	return nil
}

func recordLoginAttempt(ip string, success bool) {
	if success {
		delete(loginAttempts, ip)
		return
	}

	attempt := loginAttempts[ip]
	attempt.Count++
	attempt.Time = time.Now()
	loginAttempts[ip] = attempt
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	// Set security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Type", "application/json")

	// Check rate limiting
	ip := r.RemoteAddr
	if err := checkRateLimit(ip); err != nil {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Parse request
	var req struct {
		AgentID   string `json:"agentId"`
		AccessKey string `json:"accessKey"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Validate input
	if req.AgentID == "" || req.AccessKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   "Agent ID and Access Key are required",
		})
		return
	}

	// Check credentials
	agent, exists := agents[req.AgentID]
	if !exists || agent.AccessKey != req.AccessKey {
		recordLoginAttempt(ip, false)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   "Invalid Agent ID or Access Key",
		})
		return
	}

	// Authentication successful
	recordLoginAttempt(ip, true)

	// Generate session token
	token, err := generateToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	// Create session
	session, _ := store.Get(r, "sailtix_docs")
	sessionData := SessionData{
		AgentID:     agent.ID,
		AgentName:   agent.Name,
		Role:        agent.Role,
		Permissions: agent.Permissions,
		Token:       token,
		Expires:     time.Now().Add(time.Duration(config.SessionDuration) * time.Hour),
		IP:          ip,
		UserAgent:   r.UserAgent(),
	}

	session.Values["data"] = sessionData
	session.Save(r, w)

	// Log successful login
	log.Printf("Successful login: Agent ID: %s, IP: %s, Time: %s",
		agent.ID, ip, time.Now().Format("2006-01-02 15:04:05"))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "Authentication successful",
		Agent:   &agent,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sailtix_docs")
	session.Options.MaxAge = -1
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

func checkAuth(r *http.Request) (*SessionData, error) {
	session, err := store.Get(r, "sailtix_docs")
	if err != nil {
		return nil, err
	}

	data, ok := session.Values["data"].(SessionData)
	if !ok {
		return nil, fmt.Errorf("invalid session")
	}

	if time.Now().After(data.Expires) {
		return nil, fmt.Errorf("session expired")
	}

	return &data, nil
}

func apiSpecHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	_, err := checkAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Get requested file
	vars := mux.Vars(r)
	filename := vars["file"]

	// Validate filename
	allowedFiles := []string{"openapi.yaml", "openapi-v2.yaml"}
	allowed := false
	for _, allowedFile := range allowedFiles {
		if filename == allowedFile {
			allowed = true
			break
		}
	}

	if !allowed {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   "File not found",
		})
		return
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(AuthResponse{
			Success: false,
			Error:   "File not found",
		})
		return
	}

	// Serve file
	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	http.ServeFile(w, r, filename)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	sessionData, err := checkAuth(r)
	if err != nil {
		// Redirect to login page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Load and serve the documentation page
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Pass session data to template
	data := struct {
		Authenticated bool
		Agent         *SessionData
	}{
		Authenticated: true,
		Agent:         sessionData,
	}

	tmpl.Execute(w, data)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication for sensitive files
	if strings.HasSuffix(r.URL.Path, ".yaml") || strings.HasSuffix(r.URL.Path, ".yml") {
		_, err := checkAuth(r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	// Serve static files
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/auth", authenticateHandler).Methods("POST")
	r.HandleFunc("/api/logout", logoutHandler).Methods("DELETE")
	r.HandleFunc("/api/spec/{file}", apiSpecHandler).Methods("GET")

	// Documentation routes
	r.HandleFunc("/docs", docsHandler).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve login page
		http.ServeFile(w, r, "index.html")
	}).Methods("GET")

	// Static files (with authentication for sensitive files)
	r.PathPrefix("/").HandlerFunc(staticHandler)

	// Start server
	log.Printf("Starting Sailtix Documentation Server on port %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, r))
}
