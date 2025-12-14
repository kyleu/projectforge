# Client (TypeScript Frontend)

This directory contains the TypeScript client code for Project Forge, providing progressive enhancement for this application's UI.

## Requirements

- Node.js `>=18.18`
- npm (lockfile included)

## Project Structure

```
client/
├── src/                   # TypeScript source files
│   ├── client.css         # Main stylesheet, includes other files, included in almost all pages
│   ├── client.ts          # Main entry point, includes other files, referenced in almost all pages
│   ├── app.ts             # Custom code for this application (the rest is provided)
│   ├── *.ts               # Individual modules, included in `client`
│   ├── style/             # CSS files, add custom rules in `client.css`
│   └── svg/               # SVG icons, usually processed by Project Forge
```

## Build System

The project uses **esbuild** for fast bundling:

- **Entry point**: `src/client.ts`
- **Output**: `../assets/client.js` (consumed by Go application)
- **Bundle**: Single minified file with sourcemap
- **Target**: ES2021 for broad browser compatibility

### Build Configuration

- **TypeScript**: Strict mode, ES2021 target, JSX support
- **Module system**: ESM (bundled by esbuild; `tsc` is typecheck-only)
- **JSX**: Custom JSX factory for lightweight React-like syntax

## Development Workflow

### Common Commands

```bash
# Install dependencies (recommended)
npm ci

# Dev build (watch for changes)
npm run dev

# Build (typecheck + bundle)
npm run build

# Bundle only (skips typecheck)
npm run build:bundle

# Run tests
npm run test

# Format code
npm run format

# Check formatting
npm run format:check

# Lint code
npm run lint
```

### Development Integration

The client build is integrated with the main Go development workflow:

To build the client application, watching for changes, run:

```bash
./bin/build/client-watch.sh
```

If the main Golang application was started with `./bin/dev.sh`, the server will automatically rebuild when the client code changes.

## Code Quality

### Formatting & Linting

- **Prettier**: Consistent code formatting
- **ESLint**: TypeScript-specific linting rules
- **TypeScript**: Strict type checking with comprehensive compiler options

### Configuration Files

- `.prettierrc`: Code formatting rules
- `eslint.config.js`: Linting configuration
- `tsconfig.json`: TypeScript compiler settings

## Integration with Go Application

The TypeScript code builds to `../assets/client.js`, which is:

- Embedded in the Go binary
- Served at `/assets/client.js`
- Included in HTML templates

The CSS imported by `src/client.ts` is emitted as `../assets/client.css`.

## Performance

- **Bundle size**: Optimized with esbuild minification
- **Loading**: Single file reduces HTTP requests
- **Execution**: Lightweight runtime with minimal dependencies
- **Caching**: Long-term caching via content-based URLs
