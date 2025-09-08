// Utility functions for MyCoin Web App

class Utils {
    // Format number to display with proper decimals
    static formatNumber(num, decimals = 2) {
        return parseFloat(num).toFixed(decimals);
    }
    
    // Format MyCoin amount
    static formatMYC(amount, decimals = 2) {
        return `${this.formatNumber(amount, decimals)} MYC`;
    }
    
    // Format address for display (show first 6 and last 4 characters)
    static formatAddress(address, startLength = 6, endLength = 4) {
        if (!address || address.length < 10) return address;
        return `${address.substring(0, startLength)}...${address.substring(address.length - endLength)}`;
    }
    
    // Format hash for display
    static formatHash(hash, length = 8) {
        if (!hash) return 'N/A';
        return hash.substring(0, length) + '...';
    }
    
    // Format timestamp to readable date
    static formatDate(timestamp) {
        if (!timestamp) return 'N/A';
        const date = new Date(timestamp * 1000); // Convert seconds to milliseconds
        return date.toLocaleString('vi-VN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
    }
    
    // Format relative time (e.g., "2 hours ago")
    static formatRelativeTime(timestamp) {
        if (!timestamp) return 'N/A';
        const now = Date.now();
        const time = timestamp * 1000;
        const diff = now - time;
        
        const seconds = Math.floor(diff / 1000);
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);
        const days = Math.floor(hours / 24);
        
        if (days > 0) return `${days} ngày trước`;
        if (hours > 0) return `${hours} giờ trước`;
        if (minutes > 0) return `${minutes} phút trước`;
        return `${seconds} giây trước`;
    }
    
    // Validate MyCoin address
    static isValidAddress(address) {
        if (!address) return false;
        // MyCoin addresses are hex strings, typically 50-54 characters
        // Allow flexible length for different address formats
        return /^[a-fA-F0-9]{40,}$/.test(address) && address.length >= 40;
    }
    
    // Validate private key
    static isValidPrivateKey(privateKey) {
        if (!privateKey) return false;
        // Basic validation for hex string of 64 characters (32 bytes)
        return /^[a-fA-F0-9]{64}$/.test(privateKey);
    }
    
    // Validate amount
    static isValidAmount(amount) {
        const num = parseFloat(amount);
        return !isNaN(num) && num > 0;
    }
    
    // Copy text to clipboard
    static async copyToClipboard(text) {
        try {
            await navigator.clipboard.writeText(text);
            this.showToast('Đã sao chép vào clipboard', 'success');
        } catch (err) {
            console.error('Failed to copy to clipboard:', err);
            // Fallback for older browsers
            const textArea = document.createElement('textarea');
            textArea.value = text;
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            try {
                document.execCommand('copy');
                this.showToast('Đã sao chép vào clipboard', 'success');
            } catch (fallbackErr) {
                console.error('Fallback copy failed:', fallbackErr);
                this.showToast('Không thể sao chép', 'error');
            }
            document.body.removeChild(textArea);
        }
    }
    
    // Show toast notification
    static showToast(message, type = 'info', duration = CONFIG.TOAST_DURATION) {
        const container = document.getElementById('toast-container');
        if (!container) return;
        
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        toast.innerHTML = `
            ${message}
            <button class="toast-close" onclick="this.parentElement.remove()">×</button>
        `;
        
        container.appendChild(toast);
        
        // Auto remove after duration
        setTimeout(() => {
            if (toast.parentElement) {
                toast.remove();
            }
        }, duration);
    }
    
    // Show loading overlay
    static showLoading(show = true) {
        const overlay = document.getElementById('loading-overlay');
        if (overlay) {
            if (show) {
                overlay.classList.add('active');
            } else {
                overlay.classList.remove('active');
            }
        }
    }
    
    // Debounce function
    static debounce(func, delay) {
        let timeoutId;
        return function (...args) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => func.apply(this, args), delay);
        };
    }
    
    // Throttle function
    static throttle(func, delay) {
        let timeoutId;
        let lastExecTime = 0;
        
        return function (...args) {
            const currentTime = Date.now();
            
            if (currentTime - lastExecTime > delay) {
                func.apply(this, args);
                lastExecTime = currentTime;
            } else {
                clearTimeout(timeoutId);
                timeoutId = setTimeout(() => {
                    func.apply(this, args);
                    lastExecTime = Date.now();
                }, delay - (currentTime - lastExecTime));
            }
        };
    }
    
    // Local storage helpers
    static setStorage(key, value) {
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch (err) {
            console.error('Failed to save to localStorage:', err);
        }
    }
    
    static getStorage(key, defaultValue = null) {
        try {
            const item = localStorage.getItem(key);
            return item ? JSON.parse(item) : defaultValue;
        } catch (err) {
            console.error('Failed to load from localStorage:', err);
            return defaultValue;
        }
    }
    
    static removeStorage(key) {
        try {
            localStorage.removeItem(key);
        } catch (err) {
            console.error('Failed to remove from localStorage:', err);
        }
    }
    
    // Generate random color for avatars/icons
    static generateColor(seed) {
        const colors = [
            '#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7',
            '#DDA0DD', '#98D8C8', '#F7DC6F', '#BB8FCE', '#85C1E9'
        ];
        const index = Math.abs(this.hashCode(seed)) % colors.length;
        return colors[index];
    }
    
    // Simple hash function for strings
    static hashCode(str) {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            const char = str.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash = hash & hash; // Convert to 32-bit integer
        }
        return hash;
    }
    
    // Parse URL parameters
    static getURLParams() {
        const params = new URLSearchParams(window.location.search);
        const result = {};
        for (const [key, value] of params) {
            result[key] = value;
        }
        return result;
    }
    
    // Update URL without refresh
    static updateURL(params) {
        const url = new URL(window.location);
        Object.keys(params).forEach(key => {
            if (params[key] !== null && params[key] !== undefined) {
                url.searchParams.set(key, params[key]);
            } else {
                url.searchParams.delete(key);
            }
        });
        window.history.pushState({}, '', url);
    }
    
    // Format file size
    static formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }
    
    // Check if element is in viewport
    static isInViewport(element) {
        const rect = element.getBoundingClientRect();
        return (
            rect.top >= 0 &&
            rect.left >= 0 &&
            rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
            rect.right <= (window.innerWidth || document.documentElement.clientWidth)
        );
    }
    
    // Smooth scroll to element
    static scrollToElement(element, offset = 0) {
        const elementPosition = element.getBoundingClientRect().top;
        const offsetPosition = elementPosition + window.pageYOffset - offset;
        
        window.scrollTo({
            top: offsetPosition,
            behavior: 'smooth'
        });
    }
    
    // Check if user prefers dark mode
    static prefersDarkMode() {
        return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    
    // Generate QR code data URL (simple implementation)
    static generateQRCode(text, size = 200) {
        // This is a placeholder - in a real app you'd use a QR code library
        return `https://api.qrserver.com/v1/create-qr-code/?size=${size}x${size}&data=${encodeURIComponent(text)}`;
    }
    
    // Validate email
    static isValidEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }
    
    // Generate random ID
    static generateId(length = 8) {
        const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
        let result = '';
        for (let i = 0; i < length; i++) {
            result += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return result;
    }
}

// Global utility functions for backward compatibility
function copyToClipboard(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        Utils.copyToClipboard(element.value);
    }
}

function togglePasswordVisibility(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        const button = element.parentElement.querySelector('button i');
        if (element.type === 'password') {
            element.type = 'text';
            button.className = 'fas fa-eye-slash';
        } else {
            element.type = 'password';
            button.className = 'fas fa-eye';
        }
    }
}