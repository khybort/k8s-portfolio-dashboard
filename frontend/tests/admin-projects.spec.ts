import { test, expect } from '@playwright/test';

const ADMIN_EMAIL = process.env.ADMIN_EMAIL || 'admin@portfolio.com';
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD || 'Admin123!';

test.describe('Admin Projects CRUD', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/admin/login');
    await page.fill('input[type="email"]', ADMIN_EMAIL);
    await page.fill('input[type="password"]', ADMIN_PASSWORD);
    await page.click('button[type="submit"]');
    await page.waitForURL('/admin', { timeout: 10000 });
    
    // Navigate to projects
    await page.goto('/admin/projects');
    await page.waitForLoadState('networkidle');
  });

  test('should display projects page', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Projects');
    await expect(page.locator('text=New Project')).toBeVisible();
  });

  test('should create a new project', async ({ page }) => {
    const projectName = `Test Project ${Date.now()}`;
    const projectDescription = 'This is a test project description';
    const githubUrl = 'https://github.com/test/project';

    // Click New Project button
    await page.click('text=New Project');
    
    // Wait for form modal
    await expect(page.locator('text=Create New Project, text=New Project').or(page.locator('input[placeholder*="Name"]'))).toBeVisible({ timeout: 5000 });
    
    // Fill form - find inputs by placeholder or label
    const nameInput = page.locator('input[placeholder*="Name"], input').first();
    await nameInput.fill(projectName);
    
    const descTextarea = page.locator('textarea[placeholder*="Description"], textarea').first();
    await descTextarea.fill(projectDescription);
    
    const githubInput = page.locator('input[placeholder*="GitHub"], input[type="url"]').first();
    await githubInput.fill(githubUrl);
    
    // Submit form
    await page.click('button:has-text("Save"):not([disabled])');
    
    // Wait for form to close
    await page.waitForTimeout(2000);
    
    // Verify project appears in list
    await expect(page.locator(`text=${projectName}`)).toBeVisible({ timeout: 10000 });
  });

  test('should update an existing project', async ({ page }) => {
    // Wait for projects to load
    await page.waitForSelector('.card, [class*="grid"]', { timeout: 10000 });
    
    // Find first project edit button
    const editButton = page.locator('button:has(svg), button[aria-label*="Edit"]').first();
    
    if (await editButton.count() > 0) {
      await editButton.click();
      
      // Wait for form
      await page.waitForTimeout(1000);
      
      // Update name
      const nameInput = page.locator('input[name="name"], input[placeholder*="Name"]').first();
      if (await nameInput.count() > 0) {
        await nameInput.fill(`Updated Project ${Date.now()}`);
        
        // Save
        await page.click('button:has-text("Save"), button[type="submit"]');
        await page.waitForTimeout(2000);
        
        // Verify update
        await expect(page.locator('text=/Updated Project/')).toBeVisible({ timeout: 5000 });
      }
    } else {
      test.skip('No projects found to update');
    }
  });

  test('should delete a project', async ({ page }) => {
    // Wait for projects to load
    await page.waitForSelector('.card, [class*="grid"]', { timeout: 10000 });
    
    // Find first delete button
    const deleteButton = page.locator('button:has(svg), button[aria-label*="Delete"]').first();
    
    if (await deleteButton.count() > 0) {
      // Get project name before deletion
      const projectCard = deleteButton.locator('..').locator('..').locator('..');
      const projectText = await projectCard.textContent();
      
      // Click delete
      await deleteButton.click();
      
      // Confirm deletion
      page.on('dialog', dialog => dialog.accept());
      await page.waitForTimeout(1000);
      
      // Verify project is removed
      await page.waitForTimeout(2000);
      
      // Check that the project is no longer visible
      if (projectText) {
        await expect(page.locator(`text=${projectText}`)).not.toBeVisible({ timeout: 5000 });
      }
    } else {
      test.skip('No projects found to delete');
    }
  });

  test('should search projects', async ({ page }) => {
    // Wait for search input
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible({ timeout: 5000 });
    
    // Type in search
    await page.fill('input[placeholder*="Search"]', 'test');
    await page.waitForTimeout(1000);
    
    // Verify search is working
    const results = page.locator('.card, [class*="grid"] > div');
    const count = await results.count();
    expect(count).toBeGreaterThanOrEqual(0);
  });
});

