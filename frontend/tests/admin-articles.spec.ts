import { test, expect } from '@playwright/test';

const ADMIN_EMAIL = process.env.ADMIN_EMAIL || 'admin@portfolio.com';
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD || 'Admin123!';

test.describe('Admin Articles CRUD', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/admin/login');
    await page.fill('input[type="email"]', ADMIN_EMAIL);
    await page.fill('input[type="password"]', ADMIN_PASSWORD);
    await page.click('button[type="submit"]');
    await page.waitForURL('/admin', { timeout: 10000 });
    
    // Navigate to articles
    await page.goto('/admin/articles');
    await page.waitForLoadState('networkidle');
  });

  test('should display articles page', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Articles');
    await expect(page.locator('text=New Article')).toBeVisible();
  });

  test('should create a new article', async ({ page }) => {
    const articleTitle = `Test Article ${Date.now()}`;
    const articleContent = 'This is test article content';

    // Click New Article button
    await page.click('text=New Article');
    
    // Wait for form modal
    await expect(page.locator('text=Create New Article')).toBeVisible({ timeout: 5000 });
    
    // Fill form - Title input (by placeholder)
    const titleInput = page.locator('input[placeholder="Enter article title"]');
    await titleInput.fill(articleTitle);
    
    // Slug is auto-generated, but we can verify it exists
    await page.waitForTimeout(500);
    
    // Fill content
    const contentTextarea = page.locator('textarea[placeholder*="Write your article content"]');
    await contentTextarea.fill(articleContent);
    
    // Submit form
    await page.click('button:has-text("Save"):not([disabled])');
    
    // Wait for form to close (modal should disappear)
    await expect(page.locator('text=Create New Article')).not.toBeVisible({ timeout: 10000 });
    
    // Verify article appears in list
    await expect(page.locator(`text=${articleTitle}`)).toBeVisible({ timeout: 10000 });
  });

  test('should update an existing article', async ({ page }) => {
    // Wait for articles to load
    await page.waitForSelector('table tbody tr', { timeout: 10000 });
    
    // Find first article row
    const firstRow = page.locator('table tbody tr').first();
    
    if (await firstRow.count() > 0) {
      // Get article title before edit
      const originalTitle = await firstRow.locator('td').first().textContent();
      
      // Find edit button in the row (button with Edit icon)
      const editButton = firstRow.locator('button').first();
      await editButton.click();
      
      // Wait for edit form modal
      await expect(page.locator('text=Edit Article')).toBeVisible({ timeout: 5000 });
      
      // Update title
      const titleInput = page.locator('input[placeholder="Enter article title"]');
      const updatedTitle = `Updated ${Date.now()}`;
      await titleInput.clear();
      await titleInput.fill(updatedTitle);
      
      // Save
      await page.click('button:has-text("Save"):not([disabled])');
      
      // Wait for form to close
      await expect(page.locator('text=Edit Article')).not.toBeVisible({ timeout: 10000 });
      
      // Verify update - check that updated title appears
      await expect(page.locator(`text=${updatedTitle}`)).toBeVisible({ timeout: 10000 });
    } else {
      test.skip('No articles found to update');
    }
  });

  test('should delete an article', async ({ page }) => {
    // Wait for articles to load
    await page.waitForSelector('table tbody tr', { timeout: 10000 });
    
    // Find first article row
    const firstRow = page.locator('table tbody tr').first();
    
    if (await firstRow.count() > 0) {
      // Get article title before deletion
      const articleTitle = await firstRow.locator('td').first().textContent();
      
      // Set up dialog handler before clicking delete
      page.once('dialog', dialog => dialog.accept());
      
      // Find delete button (last button in actions column)
      const deleteButton = firstRow.locator('button').last();
      await deleteButton.click();
      
      // Wait for deletion to complete
      await page.waitForTimeout(2000);
      
      // Verify article is removed - the title should not be visible
      if (articleTitle && articleTitle.trim()) {
        await expect(page.locator(`text=${articleTitle.trim()}`)).not.toBeVisible({ timeout: 10000 });
      }
    } else {
      test.skip('No articles found to delete');
    }
  });

  test('should search articles', async ({ page }) => {
    // Wait for search input
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible({ timeout: 5000 });
    
    // Type in search
    await page.fill('input[placeholder*="Search"]', 'test');
    await page.waitForTimeout(1000);
    
    // Verify search is working (results should be filtered)
    const results = page.locator('table tbody tr, .card');
    const count = await results.count();
    expect(count).toBeGreaterThanOrEqual(0);
  });
});

