# Testing Guide

## Overview

This document describes the testing strategy and procedures for the Portfolio Dashboard project. The project includes both backend API tests and frontend UI tests.

## Test Structure

```
frontend/
  tests/
    admin-auth.spec.ts      # Authentication tests
    admin-articles.spec.ts  # Articles CRUD tests
    admin-projects.spec.ts  # Projects CRUD tests
    admin-portfolio.spec.ts # Portfolio update tests
```

## Prerequisites

1. **Node.js and npm** installed
2. **Playwright** installed (automatically installed with `npm install`)
3. **Backend services** running (backend, auth-service, frontend)
4. **Database** seeded with test data

## Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Install Playwright Browsers

```bash
npx playwright install
```

### 3. Start Services

Make sure all services are running:

```bash
# From project root
./scripts/start.sh
```

Or manually:
- Backend: `http://localhost:8080`
- Auth Service: `http://localhost:8081`
- Frontend: `http://localhost:5173`

## Running Tests

### Run All Tests

```bash
cd frontend
npx playwright test
```

### Run Specific Test Suite

```bash
# Authentication tests
npx playwright test admin-auth

# Articles CRUD tests
npx playwright test admin-articles

# Projects CRUD tests
npx playwright test admin-projects

# Portfolio update tests
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

### Run Tests with Specific Browser

```bash
npx playwright test --project=chromium
```

## Test Cases

### Authentication Tests (`admin-auth.spec.ts`)

1. **Display Login Page**
   - Verifies login page elements are visible
   - Checks email, password inputs, and submit button

2. **Successful Login**
   - Tests login with valid credentials
   - Verifies redirect to dashboard

3. **Invalid Credentials**
   - Tests login with wrong credentials
   - Verifies error message display

4. **Unauthenticated Access**
   - Tests redirect to login when not authenticated
   - Verifies protected routes are secured

### Articles CRUD Tests (`admin-articles.spec.ts`)

1. **Display Articles Page**
   - Verifies articles page loads correctly
   - Checks for "New Article" button

2. **Create Article**
   - Creates a new article with title, slug, and content
   - Verifies article appears in the list after creation

3. **Update Article**
   - Updates an existing article's title
   - Verifies changes are reflected in the list

4. **Delete Article**
   - Deletes an article
   - Verifies article is removed from the list

5. **Search Articles**
   - Tests search functionality
   - Verifies filtered results

### Projects CRUD Tests (`admin-projects.spec.ts`)

1. **Display Projects Page**
   - Verifies projects page loads correctly
   - Checks for "New Project" button

2. **Create Project**
   - Creates a new project with name, description, and GitHub URL
   - Verifies project appears in the list after creation

3. **Update Project**
   - Updates an existing project's name
   - Verifies changes are reflected in the list

4. **Delete Project**
   - Deletes a project
   - Verifies project is removed from the list

5. **Search Projects**
   - Tests search functionality
   - Verifies filtered results

### Portfolio Update Tests (`admin-portfolio.spec.ts`)

1. **Display Portfolio Page**
   - Verifies portfolio page loads correctly

2. **Update Portfolio**
   - Updates portfolio information (name, title, email, bio, social links)
   - Verifies changes are saved and displayed

3. **Cancel Edit**
   - Tests canceling portfolio edit
   - Verifies return to view mode

## Test Configuration

### Playwright Config (`playwright.config.ts`)

- **Base URL**: `http://localhost:5173`
- **Browser**: Chromium (Desktop Chrome)
- **Retries**: 2 in CI, 0 locally
- **Screenshots**: Only on failure
- **Traces**: On first retry

### Test Credentials

Default admin credentials used in tests:
- **Email**: `admin@portfolio.com`
- **Password**: `Admin123!`

## Error Handling

All CRUD operations include error handling:

1. **Error Display**: Errors are shown in red alert boxes
2. **Error Messages**: Specific error messages from API responses
3. **Error Dismissal**: Users can dismiss error messages
4. **Form Validation**: Client-side validation before submission

## Troubleshooting

### Tests Failing Due to Timeout

If tests fail due to timeout, increase timeout in test:

```typescript
test('my test', async ({ page }) => {
  await page.waitForSelector('.element', { timeout: 30000 });
});
```

### Tests Failing Due to Network Issues

Ensure all services are running:

```bash
# Check services
curl http://localhost:8080/healthz
curl http://localhost:8081/healthz
curl http://localhost:5173
```

### Tests Failing Due to Authentication

Ensure admin user exists in database:

```bash
# Check auth service logs
kubectl logs -n portfolio -l app=auth-service --tail=50
```

### View Test Reports

After running tests, view HTML report:

```bash
npx playwright show-report
```

## Continuous Integration

Tests can be run in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Install dependencies
  run: cd frontend && npm install

- name: Install Playwright
  run: cd frontend && npx playwright install --with-deps

- name: Run tests
  run: cd frontend && npx playwright test
```

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Tests should clean up created data when possible
3. **Wait Strategies**: Use proper wait strategies (waitForSelector, waitForURL)
4. **Assertions**: Use meaningful assertions with clear error messages
5. **Page Objects**: Consider using Page Object Model for complex pages

## Manual Testing Checklist

### Articles
- [ ] Create article with all fields
- [ ] Create article with minimal fields
- [ ] Update article title
- [ ] Update article content
- [ ] Delete article
- [ ] Search articles
- [ ] Toggle published status
- [ ] Error handling on invalid data

### Projects
- [ ] Create project with all fields
- [ ] Create project with minimal fields
- [ ] Update project name
- [ ] Update project description
- [ ] Delete project
- [ ] Search projects
- [ ] Toggle featured status
- [ ] Add/remove technologies
- [ ] Error handling on invalid data

### Portfolio
- [ ] View portfolio information
- [ ] Update name and title
- [ ] Update email
- [ ] Update bio
- [ ] Update social links
- [ ] Cancel edit
- [ ] Error handling on invalid data

## Test Coverage Goals

- **Authentication**: 100%
- **Articles CRUD**: 100%
- **Projects CRUD**: 100%
- **Portfolio Update**: 100%
- **Error Handling**: 100%

## Future Improvements

1. **API Tests**: Add backend API integration tests
2. **E2E Tests**: Add end-to-end tests for public site
3. **Performance Tests**: Add load and performance tests
4. **Visual Regression**: Add visual regression tests
5. **Accessibility Tests**: Add a11y tests

