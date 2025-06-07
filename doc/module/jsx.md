# JSX

The **`jsx`** module provides a lightweight JSX implementation for [Project Forge](https://projectforge.dev) applications. It enables developers to write TypeScript-based client-side code using familiar JSX syntax without the overhead of React or other heavy frameworks.

## Overview

This module delivers a minimal, performance-focused JSX implementation that:

- **Lightweight**: Zero dependencies, minimal runtime overhead
- **TypeScript Integration**: Full TypeScript support with proper type definitions
- **DOM Direct**: Works directly with DOM elements for maximum performance
- **Familiar Syntax**: Uses standard JSX patterns and conventions

## Key Features

### Performance
- Direct DOM manipulation without virtual DOM overhead
- Zero framework dependencies
- Minimal bundle size impact
- Native browser performance

### Developer Experience
- Full TypeScript support with JSX namespace declarations
- Familiar React-like JSX syntax
- Error handling for common development mistakes
- Standard attribute mapping (className → class, etc.)

### Functionality
- Standard HTML elements and attributes
- Support for event handlers and dynamic content
- `dangerouslySetInnerHTML` for raw HTML injection
- Proper handling of boolean attributes and null/undefined values

## Implementation Details

### Core JSX Function

The module exports a single `JSX` function that:
- Creates DOM elements from JSX syntax
- Handles attribute mapping and normalization
- Manages child element insertion and text nodes
- Provides comprehensive error checking

### Attribute Handling
- **`className`** → **`class`** - Standard React-style class attribute mapping
- **Boolean attributes** - Properly handles boolean values (true sets attribute, false removes it)
- **Event handlers** - Standard DOM event handling
- **Dynamic content** - Support for `dangerouslySetInnerHTML` pattern

### Child Management
- Text nodes are automatically created for string content
- Array children are properly flattened and inserted
- Comprehensive null/undefined checking with helpful error messages
- Mixed content support (elements and text)

## Usage Example

```typescript
import { JSX } from "./jsx";

// Create a simple element
const button = <button className="btn primary" onClick={handleClick}>
  Click Me
</button>;

// Create with dynamic content
const list = (
  <ul className="items">
    {items.map(item => (
      <li key={item.id}>{item.name}</li>
    ))}
  </ul>
);

// Raw HTML injection
const content = (
  <div dangerouslySetInnerHTML={{__html: htmlString}} />
);
```

## Package Structure

### Client Code
- **`client/src/jsx.ts`** - Core JSX implementation
  - JSX function for element creation
  - TypeScript namespace declarations
  - Attribute and child handling logic

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/jsx
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [TypeScript Configuration](../typescript.md) - TypeScript setup and configuration
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
