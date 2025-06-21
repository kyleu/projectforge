# Link Augments

Project Forge provides progressive enhancement for HTML links, adding JavaScript functionality while maintaining full accessibility for users without JavaScript. The link augmentation system focuses on improving user experience through confirmation dialogs and other interactive behaviors.

## Overview

Link augments work by:
- Adding CSS classes to standard HTML links
- Automatically detecting and enhancing these links when JavaScript loads
- Providing fallback behavior for users without JavaScript
- Maintaining semantic HTML structure and accessibility

## Confirmation Links

### Basic Usage

Create links that require user confirmation before navigation by adding the `link-confirm` class:

```html
<a href="/dangerous-action" class="link-confirm">Delete Item</a>
```

This will show a browser confirmation dialog with the default message "Are you sure?" before allowing navigation.

### Custom Confirmation Messages

Provide custom confirmation messages using the `data-message` attribute:

```html
<a href="/delete-user/123" class="link-confirm" data-message="Are you sure you want to delete this user? This action cannot be undone.">
  Delete User
</a>

<a href="/reset-database" class="link-confirm" data-message="This will reset all data. Continue?">
  Reset Database
</a>

<a href="/logout" class="link-confirm" data-message="Are you sure you want to log out?">
  Sign Out
</a>
```

### Progressive Enhancement

The confirmation functionality is implemented as a progressive enhancement:

**With JavaScript enabled:**
- User clicks the link
- Confirmation dialog appears with the specified message
- If user clicks "OK", navigation proceeds
- If user clicks "Cancel", navigation is prevented

**Without JavaScript:**
- Link works as a normal HTML link
- User navigates directly to the target URL
- No confirmation dialog appears
