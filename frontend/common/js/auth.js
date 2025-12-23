// Shared Authentication System for Dormitory Management
class AuthManager {
    constructor() {
        this.token = localStorage.getItem('token');
        this.user = JSON.parse(localStorage.getItem('user') || '{}');
        this.role = this.user.role;
        this.API_BASE = window.location.origin;
    }

    // Login with role-based redirect
    async login(email, password) {
        try {
            const response = await fetch(`${this.API_BASE}/api/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Login gagal');
            }

            this.token = data.token;
            this.user = data.user;
            this.role = data.user.role;

            localStorage.setItem('token', this.token);
            localStorage.setItem('user', JSON.stringify(this.user));

            // Redirect based on role
            if (this.role === 'admin') {
                window.location.href = '/admin';
            } else {
                window.location.href = '/student';
            }

            return data;
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    }

    // Register (student only)
    async register(name, email, password, phone) {
        try {
            const response = await fetch(`${this.API_BASE}/api/auth/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, email, password, phone })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Registrasi gagal');
            }

            return data;
        } catch (error) {
            console.error('Register error:', error);
            throw error;
        }
    }

    // Logout
    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');

        // Clear any session data
        this.token = null;
        this.user = {};
        this.role = null;

        // Redirect to login page
        window.location.href = '/login';
    }

    // Check authentication
    isAuthenticated() {
        return !!this.token;
    }

    // Check role
    isAdmin() {
        return this.role === 'admin';
    }

    isStudent() {
        return this.role === 'student';
    }

    // Fetch with authentication
    async fetchWithAuth(url, options = {}) {
        if (!this.token) {
            this.logout();
            throw new Error('Not authenticated');
        }

        // Show loading if available
        if (typeof showLoading === 'function') showLoading();

        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${this.token}`,
            'Content-Type': 'application/json'
        };

        try {
            const response = await fetch(`${this.API_BASE}${url}`, options);

            if (response.status === 401) {
                this.logout();
                return null;
            }

            return response;
        } finally {
            // Hide loading if available
            if (typeof hideLoading === 'function') hideLoading();
        }
    }

    // Update user profile
    async updateProfile(profileData) {
        const endpoint = this.isAdmin()
            ? '/api/admin/profile'
            : '/api/student/profile';

        const response = await this.fetchWithAuth(endpoint, {
            method: 'PUT',
            body: JSON.stringify(profileData)
        });

        if (response && response.ok) {
            const updatedUser = await response.json();
            this.user = updatedUser;
            localStorage.setItem('user', JSON.stringify(updatedUser));
            return updatedUser;
        }

        throw new Error('Failed to update profile');
    }

    // Change password
    async changePassword(oldPassword, newPassword) {
        const endpoint = this.isAdmin()
            ? '/api/admin/change-password'
            : '/api/student/change-password';

        const response = await this.fetchWithAuth(endpoint, {
            method: 'POST',
            body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
        });

        if (response && response.ok) {
            return await response.json();
        }

        const error = await response.json();
        throw new Error(error.error || 'Failed to change password');
    }
}

// Initialize global auth instance
window.auth = new AuthManager();

// Navigation guard
document.addEventListener('DOMContentLoaded', function () {
    const path = window.location.pathname;

    // Skip guard for public pages
    if (path.includes('/auth/') || path === '/login' || path === '/register' || path === '/') {
        return;
    }

    // If user is not logged in and trying to access protected pages
    if (!window.auth.isAuthenticated()) {
        if (path.includes('/admin') || path.includes('/student')) {
            window.location.href = '/login';
        }
        return;
    }

    // Role-based access control
    if (path.includes('/admin') && !window.auth.isAdmin()) {
        alert('Akses ditolak. Anda akan di-redirect ke dashboard mahasiswa.');
        window.location.href = '/student';
    }

    if (path.includes('/student') && !window.auth.isStudent()) {
        alert('Akses ditolak. Anda akan di-redirect ke dashboard admin.');
        window.location.href = '/admin';
    }
});

// Loading helpers
function showLoading() {
    let overlay = document.querySelector('.loading-overlay');
    if (!overlay) {
        document.body.insertAdjacentHTML('beforeend', `
            <div class="loading-overlay">
                <div class="spinner-border text-primary" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
            </div>
        `);
    }
}

function hideLoading() {
    const overlay = document.querySelector('.loading-overlay');
    if (overlay) overlay.remove();
}
