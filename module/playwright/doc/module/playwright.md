# Playwright

The **`playwright`** module provides comprehensive end-to-end testing capabilities for [Project Forge](https://projectforge.dev) applications using [Playwright](https://playwright.dev). It enables cross-browser testing with support for multiple device configurations, accessibility scenarios, and progressive enhancement validation.

## Overview

This module adds a complete Playwright testing setup that validates your application across:

- **Multiple Browsers**: Chrome, Firefox, Safari, and Edge
- **Device Types**: Desktop and mobile configurations  
- **Accessibility Features**: Reduced motion and JavaScript-disabled modes
- **Theme Support**: Light and dark mode testing
- **Progressive Enhancement**: No-JavaScript functionality validation

## Key Features

### Cross-Browser Testing
- Chrome, Firefox, Safari, and Edge support
- Multiple viewport sizes and device emulation
- Consistent test execution across environments

### Accessibility Testing
- JavaScript-disabled mode validation
- Reduced motion preference testing
- Screen reader compatibility verification
- Progressive enhancement validation

### Visual Testing
- Full-page screenshot capture
- Visual regression detection
- Mobile-responsive layout verification
- Theme-specific visual testing

### Performance Testing
- Page load performance metrics
- Network request monitoring
- Core Web Vitals measurement
- Resource loading validation

## Test Configuration

### Browser Matrix
The module includes comprehensive browser and device coverage:

- **Desktop Chrome**: Standard, no-JS, reduced motion, dark mode variants
- **Desktop Edge**: Cross-browser compatibility testing
- **Desktop Firefox**: Mozilla engine testing
- **Desktop Safari**: WebKit engine testing
- **Mobile Safari**: iPhone 12 portrait and landscape modes
- **Mobile Chrome**: Pixel 5 device emulation

### Test Categories
- **Core Functionality**: Basic application features
- **Progressive Enhancement**: No-JavaScript fallbacks
- **Accessibility**: WCAG compliance and user preferences
- **Visual Consistency**: Cross-browser visual validation
- **Performance**: Load times and resource efficiency

## Package Structure

### Test Files

- **`test/playwright/`** - Playwright test suite root
  - `pages.spec.ts` - Core page functionality tests
  - `playwright.config.ts` - Test configuration and browser matrix
  - `package.json` - Node.js dependencies for testing

### Configuration

- **Multi-Browser Setup**: Automated testing across major browsers
- **Device Emulation**: Mobile and desktop viewport testing
- **Visual Regression**: Screenshot-based visual testing
- **CI/CD Integration**: GitHub Actions and continuous integration support

## Usage

### Running Tests

```bash
# Navigate to the test directory
cd test/playwright

# Install dependencies
npm install

# Run all tests
npx playwright test

# Run specific browser tests
npx playwright test --project=chrome
npx playwright test --project=safari.mobile

# Run with visual output
npx playwright test --headed

# Generate HTML report
npx playwright show-report
```

### Test Development

Create new test files following the existing pattern:

```typescript
import { expect, test } from '@playwright/test';

test.describe('feature name', () => {
  test('should validate functionality', async ({ page, browserName }, testInfo) => {
    await page.goto('/your-feature');
    
    // Take screenshots for visual validation
    const screenshot = await page.screenshot({ fullPage: true });
    await testInfo.attach(`feature/${browserName}`, { 
      body: screenshot, 
      contentType: 'image/png' 
    });
    
    // Perform assertions
    await expect(page).toHaveTitle(/Expected Title/);
  });
});
```

### Accessibility Testing

The configuration automatically tests:

- **No JavaScript**: Validates progressive enhancement
- **Reduced Motion**: Tests motion preference compliance
- **Dark Mode**: Ensures theme accessibility
- **Mobile Views**: Tests touch-friendly interfaces

## Integration

### CI/CD Pipeline
The module integrates with GitHub Actions for:
- Automated test execution on pull requests
- Cross-browser compatibility verification
- Visual regression detection
- Performance monitoring

### Local Development
- Automatic application server startup
- Live test execution during development
- Visual debugging with headed browser mode
- Detailed error reporting and screenshots

## Configuration

### Custom Configuration
Modify `playwright.config.ts` to:
- Add new browser configurations
- Adjust viewport sizes
- Configure additional test reporters
- Set custom timeout values

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/playwright
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Playwright Documentation](https://playwright.dev) - Complete Playwright guide
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation  
