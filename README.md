# Sailtix API Documentation - Locked Access

This documentation is **restricted to authorized agents only**. Unauthorized access is prohibited.

## 🔐 Access Control

The API documentation is protected by authentication. Only pre-authorized agents can access the documentation.

### Authorized Agents

| Agent ID | Access Key | Role | Permissions |
|----------|------------|------|-------------|
| `sailtix-agent-001` | `access-key-2024-sailtix-secure` | Production | Read, Test |
| `sailtix-agent-002` | `access-key-2024-sailtix-secure` | Production | Read, Test |
| `sailtix-dev-001` | `dev-access-key-2024` | Development | Read, Test, Debug |
| `sailtix-admin` | `admin-access-key-2024-secure` | Admin | Read, Test, Debug, Admin |

## 🚀 How to Access

1. **Navigate to the documentation URL**
2. **Enter your Agent ID** (e.g., `sailtix-agent-001`)
3. **Enter your Access Key** (e.g., `access-key-2024-sailtix-secure`)
4. **Click "Access Documentation"**

## 🔒 Security Features

- **Session Management**: 24-hour sessions with automatic logout
- **Inactivity Timeout**: Auto-logout after 30 minutes of inactivity
- **Secure Storage**: Credentials stored securely in browser localStorage
- **Access Logging**: All access attempts are logged (future enhancement)

## 📚 Available Documentation

### API Versions
- **API v1**: Legacy version (stable)
- **API v2**: Current version with enhanced features

### Documentation Formats
- **Redoc**: Modern, responsive documentation interface
- **Swagger UI**: Classic API testing interface

## 🛠️ Managing Access

### Adding New Agents

To add a new authorized agent, update the `agents.json` file:

```json
{
  "authorized_agents": {
    "new-agent-id": {
      "access_key": "new-secure-access-key",
      "name": "New Agent Name",
      "role": "production",
      "permissions": ["read", "test"]
    }
  }
}
```

### Security Best Practices

1. **Use Strong Access Keys**: Generate cryptographically secure keys
2. **Regular Rotation**: Rotate access keys periodically
3. **Principle of Least Privilege**: Only grant necessary permissions
4. **Monitor Access**: Regularly review access logs
5. **Immediate Revocation**: Revoke access immediately when agents are decommissioned

## 🔧 Configuration

### Security Settings

```json
{
  "security_settings": {
    "session_duration_hours": 24,
    "inactivity_timeout_minutes": 30,
    "max_login_attempts": 5,
    "lockout_duration_minutes": 15
  }
}
```

### Documentation Settings

```json
{
  "documentation_settings": {
    "default_version": "v1",
    "available_versions": ["v1", "v2"],
    "enable_swagger_ui": true,
    "enable_redoc": true
  }
}
```

## 🚨 Emergency Access

In case of emergency access needed:

1. **Contact**: api-support@sailtix.com
2. **Subject**: "Emergency API Documentation Access"
3. **Include**: Justification and temporary agent credentials
4. **Response Time**: Within 1 hour during business hours

## 📞 Support

For access issues or technical support:

- **Email**: api-support@sailtix.com
- **Response Time**: 24 hours
- **Priority**: High for production agents

---

**⚠️ Security Notice**: This documentation contains sensitive API information. Do not share credentials or documentation content with unauthorized parties. 