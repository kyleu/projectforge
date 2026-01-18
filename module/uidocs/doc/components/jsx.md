# JSX Templating

Project Forge includes a lightweight, custom JSX engine that enables you to write component-based UI code using familiar JSX syntax. This implementation is designed to be minimal, fast, and dependency-free while providing the core benefits of JSX for dynamic content generation.

## Overview

The JSX engine allows you to:
- Write declarative UI components using JSX syntax
- Create dynamic HTML elements programmatically
- Build reusable component functions
- Generate DOM elements that can be inserted into your pages

Unlike full frameworks like React, this JSX implementation focuses on simple DOM element creation without virtual DOM, state management, or lifecycle methods. It's perfect for progressive enhancement and dynamic content generation.

## Getting Started

To use JSX in your project, create a `.tsx` file in `./client/src` and use Project Forge to import the `jsx` module into your application:

```typescript
import * as JSX from "./jsx";

export function myComponent() {
  return <div>Check out this JSX!</div>;
}
```

## Basic Usage

### Simple Elements

Create basic HTML elements with JSX syntax:

```typescript
import * as JSX from "./jsx";

// Simple div element
function createContainer() {
  return <div className="container">Hello World</div>;
}

// Element with multiple attributes
function createButton() {
  return <button type="button" className="btn btn-primary" disabled>
    Click Me
  </button>;
}

// Element with inline styles
function createStyledElement() {
  return <div style="color: red; font-weight: bold;">
    Important Message
  </div>;
}
```

### Elements with Content

JSX supports various types of content including text, other elements, and arrays:

```typescript
import * as JSX from "./jsx";

// Text content
function createHeading() {
  return <h1>Welcome to Project Forge</h1>;
}

// Nested elements
function createCard() {
  return (
    <div className="card">
      <div className="card-header">
        <h3>Card Title</h3>
      </div>
      <div className="card-body">
        <p>This is the card content.</p>
      </div>
    </div>
  );
}

// Mixed content
function createMixedContent() {
  return (
    <div>
      <h2>Features</h2>
      <ul>
        <li>Fast and lightweight</li>
        <li>No external dependencies</li>
        <li>TypeScript support</li>
      </ul>
    </div>
  );
}
```

## Advanced Features

### Dynamic Content

Generate content dynamically using JavaScript expressions:

```typescript
import * as JSX from "./jsx";

function createUserProfile(user: {name: string, email: string, isAdmin: boolean}) {
  return (
    <div className="user-profile">
      <h2>{user.name}</h2>
      <p>Email: {user.email}</p>
      {user.isAdmin && <span className="badge">Admin</span>}
    </div>
  );
}

function createItemList(items: string[]) {
  return (
    <ul>
      {items.map(item => <li key={item}>{item}</li>)}
    </ul>
  );
}
```

### Conditional Rendering

Use JavaScript expressions for conditional rendering:

```typescript
import * as JSX from "./jsx";

function createStatusMessage(isLoading: boolean, error?: string, data?: any) {
  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  if (data) {
    return <div className="success">Data loaded successfully!</div>;
  }

  return <div className="info">No data available</div>;
}

function createConditionalButton(showButton: boolean, isEnabled: boolean) {
  return (
    <div>
      {showButton && (
        <button disabled={!isEnabled} className="btn">
          {isEnabled ? "Click Me" : "Disabled"}
        </button>
      )}
    </div>
  );
}
```

### Dangerous HTML

For cases where you need to insert raw HTML (use with caution):

```typescript
import * as JSX from "./jsx";

function createRichContent(htmlContent: string) {
  return (
    <div
      className="rich-content"
      dangerouslySetInnerHTML={{__html: htmlContent}}
    />
  );
}

// Example usage
const markdownHTML = "<p>This is <strong>bold</strong> text from markdown.</p>";
const contentElement = createRichContent(markdownHTML);
```

**⚠️ Security Warning**: Only use `dangerouslySetInnerHTML` with trusted content to prevent XSS attacks.

## Component Patterns

Create reusable component functions:

```typescript
import * as JSX from "./jsx";

// Button component with variants
function Button(props: {
  text: string;
  variant?: 'primary' | 'secondary' | 'danger';
  onClick?: () => void;
  disabled?: boolean;
}) {
  const className = `btn btn-${props.variant || 'primary'}`;

  const element = (
    <button
      className={className}
      disabled={props.disabled}
      type="button"
    >
      {props.text}
    </button>
  );

  if (props.onClick) {
    (element as HTMLButtonElement).addEventListener('click', props.onClick);
  }

  return element;
}

// Card component
function Card(props: {title: string, children: Node[]}) {
  return (
    <div className="card">
      <div className="card-header">
        <h3>{props.title}</h3>
      </div>
      <div className="card-body">
        {props.children}
      </div>
    </div>
  );
}
```

## Limitations

This JSX implementation is intentionally minimal and has some limitations compared to full frameworks:

- No virtual DOM or diffing algorithm
- No built-in state management
- No lifecycle methods
- No component re-rendering
- Limited to DOM element creation

These limitations make it perfect for progressive enhancement and dynamic content generation without the overhead of a full framework.
