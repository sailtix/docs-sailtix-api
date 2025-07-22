# Integrating Documentation with Existing Sailtix API

Since GitHub Pages can't handle authentication, here are the best approaches:

## 🎯 **Recommended: Integrate with Existing API**

### **Step 1: Add Documentation Endpoints to Your API**

Add these endpoints to your existing `api-sailtix-com`:

```go
// Add to your existing API routes
func (h *Handler) SetupDocsRoutes(r *mux.Router) {
    docs := r.PathPrefix("/docs").Subrouter()
    
    // Authentication endpoint
    docs.HandleFunc("/auth", h.DocsAuthHandler).Methods("POST")
    docs.HandleFunc("/logout", h.DocsLogoutHandler).Methods("DELETE")
    
    // Protected documentation endpoints
    docs.HandleFunc("/spec/{file}", h.DocsSpecHandler).Methods("GET")
    docs.HandleFunc("/", h.DocsIndexHandler).Methods("GET")
    
    // Middleware for authentication
    docs.Use(h.DocsAuthMiddleware)
}

// Authentication handler
func (h *Handler) DocsAuthHandler(c *gin.Context) {
    var req struct {
        AgentID   string `json:"agentId"`
        AccessKey string `json:"accessKey"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // Validate against your existing user/agent system
    agent, err := h.validateAgent(req.AgentID, req.AccessKey)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Create session using your existing session management
    session := sessions.Default(c)
    session.Set("docs_agent_id", agent.ID)
    session.Set("docs_agent_name", agent.Name)
    session.Set("docs_role", agent.Role)
    session.Save()
    
    c.JSON(200, gin.H{
        "success": true,
        "message": "Authentication successful",
        "agent": agent,
    })
}

// Serve OpenAPI specs
func (h *Handler) DocsSpecHandler(c *gin.Context) {
    filename := c.Param("file")
    
    // Validate filename
    if filename != "openapi.yaml" && filename != "openapi-v2.yaml" {
        c.JSON(404, gin.H{"error": "File not found"})
        return
    }
    
    // Serve the file from your API
    c.File("docs/" + filename)
}
```

### **Step 2: Update Frontend to Use Your API**

Update the documentation frontend to point to your API:

```javascript
// Update API endpoints in index.html
const API_BASE = 'https://api.sailtix.com/docs'; // Your API URL

function authenticate(agentId, accessKey) {
    return fetch(`${API_BASE}/auth`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ agentId, accessKey })
    })
    .then(response => response.json());
}

function loadDocumentation(version) {
    const specUrl = `${API_BASE}/spec/${version === 'v2' ? 'openapi-v2.yaml' : 'openapi.yaml'}`;
    Redoc.init(specUrl, { /* config */ });
}
```

### **Step 3: Deploy Documentation to Your API**

```bash
# Copy documentation files to your API
cp docs-sailtix-api/index.html api-sailtix-com/docs/
cp docs-sailtix-api/openapi.yaml api-sailtix-com/docs/
cp docs-sailtix-api/openapi-v2.yaml api-sailtix-com/docs/
```

## 🚀 **Alternative: Use Cloudflare Pages with Access**

### **Step 1: Deploy to Cloudflare Pages**
```bash
# Push to GitHub
git push origin main

# Connect to Cloudflare Pages
# - Connect your GitHub repo
# - Build command: (none, static files)
# - Publish directory: /
```

### **Step 2: Enable Cloudflare Access**
```bash
# In Cloudflare Dashboard:
# 1. Go to Access > Applications
# 2. Add Application
# 3. Set subdomain: docs.sailtix.com
# 4. Configure authentication rules
# 5. Add authorized users/emails
```

## 🔧 **Alternative: Use Netlify with Identity**

### **Step 1: Deploy to Netlify**
```bash
# Install Netlify CLI
npm install -g netlify-cli

# Deploy
netlify deploy --prod --dir=.
```

### **Step 2: Enable Netlify Identity**
```javascript
// Add to index.html
<script src="https://identity.netlify.com/v1/netlify-identity-widget.js"></script>

// Authentication logic
if (netlifyIdentity.currentUser()) {
    showDocumentation();
} else {
    netlifyIdentity.open();
}
```

## 📋 **Recommended Approach:**

**Use your existing API** because:

1. ✅ **Consistent authentication** with your main system
2. ✅ **Single source of truth** for user management
3. ✅ **Existing infrastructure** - no new services needed
4. ✅ **Better security** - server-side validation
5. ✅ **Easier maintenance** - one codebase

## 🎯 **Quick Implementation:**

1. **Add docs endpoints to your API**
2. **Update frontend to use your API**
3. **Deploy docs files to your API**
4. **Access via: `https://api.sailtix.com/docs`**

This way, you get **real security** with your existing infrastructure! 🔒 