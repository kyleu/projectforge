# OAuth

The **`oauth`** module provides comprehensive OAuth 2.0 authentication and session management for [Project Forge](https://projectforge.dev) applications. It supports sozens of OAuth providers and includes flexible permission systems for access control.

## Overview

This module enables secure authentication through external OAuth providers, eliminating the need for custom user management while providing fine-grained access control capabilities.

- **67+ OAuth Providers**: Support for major platforms, development tools, social networks, and international services
- **Session Management**: Secure session handling with role-based access control
- **Permission System**: Flexible rule-based permissions with domain and user-level controls
- **Progressive Enhancement**: Works without JavaScript, enhanced with it

## Key Features

### Supported Providers

**Major Platforms:**
- **Development**: GitHub, GitLab, Bitbucket, Gitea, Azure AD, Auth0, Okta
- **Business**: Google, Microsoft, Amazon, Salesforce, Slack, Shopify
- **Social**: Facebook, Twitter, Discord, LinkedIn, Instagram, TikTok
- **Gaming**: Steam, Twitch, Battlenet, Discord
- **International**: Naver, Kakao, LINE, VK, Yandex, WeChat

### Security Features
- OAuth 2.0 and OpenID Connect support
- Secure session management with configurable expiration
- CSRF protection for all OAuth flows
- Flexible redirect URL configuration
- Domain-based access control

### Permission System
- Role-based access control (RBAC)
- Path-based permissions
- Provider-specific rules (e.g., GitHub organization membership)
- Domain-based access (e.g., @company.com emails only)
- Default allow/deny policies

## Configuration

### Basic Setup

1. **Enable OAuth**: Configure provider credentials via environment variables
2. **Set Permissions**: Define access rules in your application initialization
3. **Configure Redirects**: Set up OAuth callback URLs

### Environment Variables

```bash
# GitHub OAuth
github_key=your_client_id
github_secret=your_client_secret

# Google OAuth
google_key=your_client_id
google_secret=your_client_secret

# Custom OpenID Connect
openid_connect_name="Custom Provider"
openid_connect_url=https://provider.example.com/.well-known/openid_configuration

# OAuth redirect configuration
oauth_redirect=https://yourapp.com/auth/callback
oauth_protocol=https
```

### Permission Configuration

Add permission rules to your application initialization:

```go
// Basic permission setup
user.SetPermissions(false, // default deny
    // Grant admin access to GitHub users from specific domain
    user.Perm("/admin", "github:@projectforge.dev", true),
    
    // Allow specific GitHub organization members
    user.Perm("/admin", "github:org:mycompany", true),
    
    // Deny admin access to others
    user.Perm("/admin", "*", false),
    
    // Allow authenticated users to access main app
    user.Perm("/", "*", true),
)
```

## Usage Examples

### Organization-Based Access

```go
// GitHub organization membership required
user.SetPermissions(false,
    user.Perm("/", "github:org:mycompany", true),
    user.Perm("/", "*", false),
)
```

### Domain-Based Access Control

```go
// Company email domain required
user.SetPermissions(false,
    user.Perm("/", "google:@company.com", true),
    user.Perm("/public", "*", true), // public area for all
)
```

### Multi-Provider Setup

```go
// Multiple providers with different access levels
user.SetPermissions(false,
    // Admins: GitHub org members or Google company emails
    user.Perm("/admin", "github:org:mycompany", true),
    user.Perm("/admin", "google:@company.com", true),
    
    // Users: Any authenticated user from approved providers
    user.Perm("/app", "github:*", true),
    user.Perm("/app", "google:*", true),
    user.Perm("/app", "microsoft:*", true),
)
```

## Advanced Features

### Custom Provider Setup

For OpenID Connect providers not built-in:

```bash
openid_connect_name="Custom SSO"
openid_connect_url=https://sso.company.com/.well-known/openid_configuration
openid_connect_key=your_client_id
openid_connect_secret=your_client_secret
```

### Session Customization

Configure session behavior through environment variables or programmatically:

```go
// Custom session duration
os.Setenv("session_timeout", "7200") // 2 hours

// Custom cookie settings
os.Setenv("session_secure", "true")
os.Setenv("session_same_site", "strict")
```

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/oauth
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)
- **Dependencies**: [Goth](https://github.com/markbates/goth) OAuth library

## See Also

- [User Module](user.md) - User management and profiles
- [Configuration Guide](../running.md) - Environment variable reference
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
