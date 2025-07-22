# Sailtix API Documentation

## Overview

This repository contains the comprehensive API documentation for the Sailtix maritime ticketing and booking platform. The documentation covers both the current v2.0.0 API and the legacy v1.x API for reference and migration purposes.

## 🚀 Quick Start

### Current Version (v2.0.0)
- **OpenAPI Specification**: [`openapi-v2.yaml`](./openapi-v2.yaml)
- **Changelog**: [`CHANGELOG-v2.md`](./CHANGELOG-v2.md)
- **Documentation**: [View Online](https://docs.sailtix.com)

### Legacy Version (v1.x)
- **OpenAPI Specification**: [`openapi.yaml`](./openapi.yaml)
- **Documentation**: [View Legacy API](./openapi.yaml)

## 📚 Documentation Structure

```
docs-sailtix-api/
├── index.html                 # Main documentation landing page
├── openapi-v2.yaml           # Current API v2.0.0 specification
├── openapi.yaml              # Legacy API v1.x specification
├── CHANGELOG-v2.md           # Comprehensive changelog for v2.0.0
├── README.md                 # This file
├── redocly.yaml             # Redoc configuration
└── CNAME                    # Custom domain configuration
```

## 🔄 API Versioning

### Version 2.0.0 (Current)
- **Status**: ✅ Active
- **Release Date**: January 2024
- **Breaking Changes**: Yes
- **Migration Required**: Yes (from v1.x)

### Version 1.x (Legacy)
- **Status**: ⚠️ Deprecated
- **Support**: Limited support
- **Breaking Changes**: No new changes
- **Migration**: Recommended to v2.0.0

## 🚨 Breaking Changes in v2.0.0

### Authentication & Authorization
- Updated JWT token structure with enhanced security
- Completely revamped permission system
- Enhanced Google OAuth integration
- Improved session management

### Response Format
- All responses now follow a standardized format
- Enhanced error handling with detailed error codes
- Updated pagination format
- Consistent data wrapping structure

### Payment System
- Multi-gateway support (Xendit, Stripe)
- New payment method management endpoints
- Enhanced security with token-based transactions
- Improved webhook handling

### File Upload System
- Enhanced validation with better error messages
- Support for multiple storage backends
- Automatic image optimization
- Enhanced security measures

## 🛠️ Development

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/sailtix/docs-sailtix-api.git
   cd docs-sailtix-api
   ```

2. **View documentation locally**
   ```bash
   # Using Python
   python -m http.server 8000
   
   # Using Node.js
   npx serve .
   
   # Using PHP
   php -S localhost:8000
   ```

3. **Open in browser**
   ```
   http://localhost:8000
   ```

### API Specification Validation

```bash
# Validate OpenAPI specification
npx @redocly/cli lint openapi-v2.yaml

# Generate documentation
npx @redocly/cli build-docs openapi-v2.yaml -o docs.html
```

### Documentation Generation

The API documentation is automatically generated from the Go codebase using Swagger annotations. To regenerate:

```bash
# From the API repository
cd api-sailtix-com
swag init -g main.go -o ../docs-sailtix-api
```

## 📖 API Features

### Core Functionality
- **Maritime Ticketing**: Route pricing, vessel management, dock operations
- **Payment Processing**: Multi-gateway support (Xendit, Stripe)
- **Order Management**: Complete booking lifecycle
- **Authentication**: JWT-based auth with Google OAuth
- **File Management**: Image uploads and storage
- **Webhooks**: Real-time payment and order notifications

### Enhanced Features in v2.0.0
- **Real-time Order Tracking**: Live order status updates
- **Advanced Search**: Enhanced filtering and search capabilities
- **Multi-language Support**: Internationalization ready
- **Analytics Dashboard**: Comprehensive reporting
- **Mobile Optimization**: Mobile-first API design
- **Rate Limiting**: Comprehensive rate limiting
- **Caching**: Intelligent caching for performance

## 🔌 Integration Examples

### Authentication
```javascript
// Login
const response = await fetch('/landing/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});

const { data: { token, user } } = await response.json();
```

### Search Routes
```javascript
// Search for routes
const response = await fetch('/landing/route-pricings/search', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    departure_dock_id: 'dock-uuid',
    arrival_dock_id: 'dock-uuid',
    departure_date: '2024-12-25',
    adult: 2,
    child: 1,
    infant: 0
  })
});

const { data: { route_pricings, meta } } = await response.json();
```

### Create Order
```javascript
// Create booking order
const response = await fetch('/landing/orders', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    route_pricing_id: 'pricing-uuid',
    departure_date: '2024-12-25',
    adult_count: 2,
    child_count: 1,
    infant_count: 0,
    passenger_details: [...],
    contact_info: {...}
  })
});

const { data: { order } } = await response.json();
```

### Payment Processing
```javascript
// Create payment request
const response = await fetch('/landing/payments', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    order_id: 'order-uuid',
    payment_method: 'VIRTUAL_ACCOUNT',
    payment_channel: 'BCA'
  })
});

const { data: { payment_url, payment_id } } = await response.json();
```

## 🔧 Configuration

### Environment Variables
```bash
# API Base URL
API_BASE_URL=https://api.sailtix.com

# Authentication
JWT_SECRET=your-jwt-secret

# Payment Gateways
XENDIT_API_KEY=your-xendit-key
STRIPE_SECRET_KEY=your-stripe-key

# File Storage
STORAGE_TYPE=local  # local, aws_s3, gcs
STORAGE_BASE_URI=https://assets.sailtix.com
```

### CORS Configuration
```yaml
cors_allowed_origins:
  - https://sailtix.com
  - https://dash.sailtix.com
  - http://localhost:3000
```

## 📊 Performance Metrics

### API Performance (v2.0.0)
- **Response Time**: 50% improvement in average response time
- **Throughput**: 3x increase in requests per second
- **Error Rate**: 90% reduction in API errors
- **Uptime**: 99.9% uptime with improved reliability

### Database Performance
- **Query Time**: 60% improvement in average query time
- **Connection Pool**: 5x increase in concurrent connections
- **Index Efficiency**: 80% improvement in index usage
- **Storage Optimization**: 40% reduction in storage usage

## 🚀 Deployment

### Documentation Site
The documentation is deployed to GitHub Pages and can be accessed at:
- **Production**: https://docs.sailtix.com
- **Development**: https://sailtix.github.io/docs-sailtix-api

### API Endpoints
- **Production**: https://api.sailtix.com
- **Development**: http://localhost:8080
- **Staging**: https://staging-api.sailtix.com

## 🛡️ Security

### Authentication
- JWT-based authentication with enhanced security
- OAuth 2.0 integration with Google
- Session management with automatic token refresh
- Rate limiting and brute force protection

### Data Protection
- All sensitive data is encrypted at rest
- HTTPS/TLS encryption for all communications
- Input validation and sanitization
- SQL injection and XSS protection

### Payment Security
- PCI DSS compliant payment processing
- Token-based payment methods
- Secure webhook validation
- Fraud detection and prevention

## 📞 Support

### Technical Support
- **Email**: api-support@sailtix.com
- **Documentation**: https://docs.sailtix.com
- **Status Page**: https://status.sailtix.com
- **GitHub Issues**: https://github.com/sailtix/api-sailtix-com/issues

### Migration Support
- **Migration Guide**: See [CHANGELOG-v2.md](./CHANGELOG-v2.md)
- **Migration Tools**: Available in the API repository
- **Support Team**: Dedicated migration assistance
- **Training Sessions**: Available for enterprise customers

### Community
- **Developer Forum**: https://community.sailtix.com
- **Stack Overflow**: Tagged with `sailtix-api`
- **Discord**: https://discord.gg/sailtix
- **Twitter**: @sailtix_api

## 📄 License

This documentation is proprietary and confidential. All rights reserved by Sailtix.

## 🤝 Contributing

For API improvements and bug reports, please contribute to the main API repository:
https://github.com/sailtix/api-sailtix-com

For documentation improvements, please submit pull requests to this repository.

---

**© 2024 Sailtix. All rights reserved.** 