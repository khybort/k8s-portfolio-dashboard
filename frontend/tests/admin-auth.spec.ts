import { test, expect } from '@playwright/test';

const ADMIN_EMAIL = process.env.ADMIN_EMAIL || 'admin@portfolio.com';
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD || 'Admin123!';

test.describe('Admin Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/admin/login');
  });

  test('should display login page', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Admin Login');
    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('input[type="password"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('should login successfully with valid credentials', async ({ page }) => {
    await page.fill('input[type="email"]', ADMIN_EMAIL);
    await page.fill('input[type="password"]', ADMIN_PASSWORD);
    await page.click('button[type="submit"]');
    
    // Wait for redirect to dashboard
    await page.waitForURL('/admin', { timeout: 10000 });
    await expect(page).toHaveURL('/admin');
  });

  test('should show error with invalid credentials', async ({ page }) => {
    await page.fill('input[type="email"]', 'wrong@email.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');
    
    // Wait for error message
    await expect(page.locator('text=/Login failed|Invalid credentials/i')).toBeVisible({ timeout: 5000 });
  });

  test('should redirect to login when not authenticated', async ({ page }) => {
    await page.goto('/admin');
    await page.waitForURL('/admin/login', { timeout: 5000 });
    await expect(page).toHaveURL('/admin/login');
  });
});

