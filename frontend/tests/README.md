# UI Tests

This directory contains Playwright end-to-end tests for the admin panel.

## Test Files

- `admin-auth.spec.ts` - Authentication tests (login, logout, protected routes)
- `admin-articles.spec.ts` - Articles CRUD tests (create, read, update, delete, search)
- `admin-projects.spec.ts` - Projects CRUD tests (create, read, update, delete, search)
- `admin-portfolio.spec.ts` - Portfolio update tests

## Running Tests

### Prerequisites

1. All services must be running:
   - Backend: `http://localhost:8080`
   - Auth Service: `http://localhost:8081`
   - Frontend: `http://localhost:5173`

2. Admin user must exist:
   - Email: `admin@portfolio.com`
   - Password: `Admin123!`

### Run All Tests

```bash
cd frontend
npm install
npx playwright install
npx playwright test
```

### Run Specific Test Suite

```bash
# Authentication tests
npx playwright test admin-auth

# Articles tests
npx playwright test admin-articles

# Projects tests
npx playwright test admin-projects

# Portfolio tests
npx playwright test admin-portfolio
```

### Run Tests in UI Mode

```bash
npx playwright test --ui
```

### Run Tests in Debug Mode

```bash
npx playwright test --debug
```

### View Test Report

```bash
npx playwright show-report
```

## Test Structure

Each test file follows this structure:

1. **Before Each Hook**: Logs in and navigates to the relevant page
2. **Test Cases**: Individual test scenarios
3. **Assertions**: Verify expected behavior

## Common Issues

### Tests Fail with Timeout

- Ensure all services are running
- Check network connectivity
- Increase timeout if needed

### Tests Fail with Authentication Error

- Verify admin user exists in database
- Check auth service is running
- Verify credentials are correct

### Tests Fail to Find Elements

- Check if UI has changed
- Verify selectors are still valid
- Use Playwright Inspector to debug: `npx playwright test --debug`

## Writing New Tests

When adding new tests:

1. Use descriptive test names
2. Wait for elements before interacting
3. Use proper selectors (prefer text, role, or data-testid)
4. Clean up test data when possible
5. Add error handling for edge cases

## CI/CD Integration

Tests can be integrated into CI/CD pipelines:

```yaml
- name: Run E2E Tests
  run: |
    cd frontend
    npm install
    npx playwright install --with-deps
    npx playwright test
```

