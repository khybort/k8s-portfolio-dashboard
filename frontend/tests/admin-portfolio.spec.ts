import { test, expect } from '@playwright/test';

const ADMIN_EMAIL = process.env.ADMIN_EMAIL || 'admin@portfolio.com';
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD || 'Admin123!';

test.describe('Admin Portfolio Update', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/admin/login');
    await page.fill('input[type="email"]', ADMIN_EMAIL);
    await page.fill('input[type="password"]', ADMIN_PASSWORD);
    await page.click('button[type="submit"]');
    await page.waitForURL('/admin', { timeout: 10000 });
    
    // Navigate to portfolio
    await page.goto('/admin/portfolio');
    await page.waitForLoadState('networkidle');
  });

  test('should display portfolio page', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Portfolio Settings');
  });

  test('should update portfolio information', async ({ page }) => {
    // Click Edit button
    await page.click('text=Edit Portfolio');
    
    // Wait for form
    await expect(page.locator('input[name="name"], input[placeholder*="Name"]')).toBeVisible({ timeout: 5000 });
    
    // Update fields
    const nameInput = page.locator('input[name="name"], input[placeholder*="Name"]').first();
    if (await nameInput.count() > 0) {
      await nameInput.fill(`Updated Name ${Date.now()}`);
      
      // Save
      await page.click('button:has-text("Save"), button[type="submit"]');
      await page.waitForTimeout(2000);
      
      // Verify update (should return to view mode)
      await expect(page.locator('text=/Updated Name/')).toBeVisible({ timeout: 10000 });
    }
  });

  test('should cancel portfolio edit', async ({ page }) => {
    // Click Edit button
    await page.click('text=Edit Portfolio');
    
    // Wait for form
    await page.waitForTimeout(1000);
    
    // Click Cancel
    await page.click('button:has-text("Cancel")');
    await page.waitForTimeout(1000);
    
    // Should return to view mode
    await expect(page.locator('text=Edit Portfolio')).toBeVisible();
  });
});

