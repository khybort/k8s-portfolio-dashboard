// Portfolio Dashboard - Main JavaScript

(function() {
  'use strict';

  // Theme Management
  const ThemeManager = {
    init: function() {
      const savedTheme = localStorage.getItem('theme');
      if (savedTheme) {
        document.documentElement.classList.toggle('light', savedTheme === 'light');
      } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
        document.documentElement.classList.add('light');
      }
      this.updateThemeButton();
    },
    
    toggle: function() {
      document.documentElement.classList.toggle('light');
      const isLight = document.documentElement.classList.contains('light');
      localStorage.setItem('theme', isLight ? 'light' : 'dark');
      this.updateThemeButton();
    },
    
    updateThemeButton: function() {
      const btn = document.querySelector('[data-theme-toggle]');
      if (btn) {
        btn.textContent = document.documentElement.classList.contains('light') 
          ? '🌙 Dark' 
          : '☀️ Light';
      }
    }
  };

  // Uptime Counter
  const UptimeCounter = {
    start: Date.now(),
    
    init: function() {
      this.tick();
      setInterval(() => this.tick(), 1000);
    },
    
    tick: function() {
      const element = document.getElementById('uptime');
      if (!element) return;
      
      const seconds = Math.floor((Date.now() - this.start) / 1000);
      const minutes = Math.floor(seconds / 60);
      const hours = Math.floor(minutes / 60);
      
      const ss = seconds % 60;
      const mm = minutes % 60;
      const hh = hours;
      
      let text = '';
      if (hh > 0) {
        text = `${hh}h ${this.pad(mm)}m ${this.pad(ss)}s`;
      } else if (mm > 0) {
        text = `${mm}m ${this.pad(ss)}s`;
      } else {
        text = `${ss}s`;
      }
      
      element.textContent = text;
    },
    
    pad: function(n) {
      return n < 10 ? '0' + n : '' + n;
    }
  };

  // Health Check
  const HealthCheck = {
    init: function() {
      this.check();
      setInterval(() => this.check(), 30000); // Her 30 saniyede bir
    },
    
    check: function() {
      fetch('/healthz')
        .then(response => {
          if (response.ok) {
            this.updateStatus('healthy', 'Running');
          } else {
            this.updateStatus('unhealthy', 'Error');
          }
        })
        .catch(() => {
          this.updateStatus('unhealthy', 'Offline');
        });
    },
    
    updateStatus: function(status, text) {
      const statusText = document.getElementById('statusText');
      const dot = document.querySelector('.dot');
      
      if (statusText) {
        statusText.textContent = text;
      }
      
      if (dot) {
        dot.className = 'dot ' + status;
      }
    }
  };

  // Form Validation
  const FormValidator = {
    init: function() {
      const forms = document.querySelectorAll('form[data-validate]');
      forms.forEach(form => {
        form.addEventListener('submit', (e) => this.handleSubmit(e, form));
        
        // Real-time validation
        const inputs = form.querySelectorAll('input, textarea');
        inputs.forEach(input => {
          input.addEventListener('blur', () => this.validateField(input));
          input.addEventListener('input', () => this.clearError(input));
        });
      });
    },
    
    handleSubmit: function(e, form) {
      e.preventDefault();
      
      const isValid = Array.from(form.querySelectorAll('input[required], textarea[required]'))
        .every(input => this.validateField(input));
      
      if (isValid) {
        this.submitForm(form);
      } else {
        this.showFormError(form, 'Lütfen tüm gerekli alanları doldurun.');
      }
    },
    
    validateField: function(field) {
      const value = field.value.trim();
      let isValid = true;
      let error = '';
      
      // Required check
      if (field.hasAttribute('required') && !value) {
        isValid = false;
        error = 'Bu alan zorunludur.';
      }
      
      // Email validation
      if (field.type === 'email' && value) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(value)) {
          isValid = false;
          error = 'Geçerli bir e-posta adresi girin.';
        }
      }
      
      // Min length
      const minLength = field.getAttribute('minlength');
      if (minLength && value.length < parseInt(minLength)) {
        isValid = false;
        error = `En az ${minLength} karakter olmalıdır.`;
      }
      
      this.setFieldError(field, isValid, error);
      return isValid;
    },
    
    setFieldError: function(field, isValid, error) {
      this.clearError(field);
      
      if (!isValid) {
        field.classList.add('error');
        const errorElement = document.createElement('span');
        errorElement.className = 'field-error';
        errorElement.textContent = error;
        field.parentNode.appendChild(errorElement);
      } else {
        field.classList.remove('error');
      }
    },
    
    clearError: function(field) {
      field.classList.remove('error');
      const errorElement = field.parentNode.querySelector('.field-error');
      if (errorElement) {
        errorElement.remove();
      }
    },
    
    showFormError: function(form, message) {
      let errorDiv = form.querySelector('.form-error');
      if (!errorDiv) {
        errorDiv = document.createElement('div');
        errorDiv.className = 'form-error';
        form.insertBefore(errorDiv, form.firstChild);
      }
      errorDiv.textContent = message;
      errorDiv.style.display = 'block';
    },
    
    submitForm: function(form) {
      // Simulate form submission
      const submitBtn = form.querySelector('button[type="submit"]');
      if (submitBtn) {
        submitBtn.disabled = true;
        submitBtn.textContent = 'Gönderiliyor...';
      }
      
      // In a real app, this would be an API call
      setTimeout(() => {
        alert('Form gönderildi! (Demo modu - gerçek bir API entegrasyonu eklenebilir)');
        form.reset();
        if (submitBtn) {
          submitBtn.disabled = false;
          submitBtn.textContent = 'Gönder';
        }
        const errorDiv = form.querySelector('.form-error');
        if (errorDiv) {
          errorDiv.style.display = 'none';
        }
      }, 1000);
    }
  };

  // Smooth Scroll
  const SmoothScroll = {
    init: function() {
      document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
          const href = this.getAttribute('href');
          if (href === '#') return;
          
          e.preventDefault();
          const target = document.querySelector(href);
          if (target) {
            target.scrollIntoView({
              behavior: 'smooth',
              block: 'start'
            });
          }
        });
      });
    }
  };

  // Initialize on DOM ready
  function init() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', init);
      return;
    }
    
    // Initialize components
    ThemeManager.init();
    UptimeCounter.init();
    HealthCheck.init();
    FormValidator.init();
    SmoothScroll.init();
    
    // Set current year
    const yearElement = document.getElementById('year');
    if (yearElement) {
      yearElement.textContent = new Date().getFullYear();
    }
    
    // Theme toggle button
    const themeToggle = document.querySelector('[data-theme-toggle]');
    if (themeToggle) {
      themeToggle.addEventListener('click', (e) => {
        e.preventDefault();
        ThemeManager.toggle();
      });
    }
  }

  // Expose toggleTheme globally for backward compatibility
  window.toggleTheme = function() {
    ThemeManager.toggle();
  };

  // Start initialization
  init();
})();

