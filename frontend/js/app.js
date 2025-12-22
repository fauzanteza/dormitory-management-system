// Global app utilities
class DormitoryApp {
    constructor() {
        this.token = localStorage.getItem('token');
        this.user = JSON.parse(localStorage.getItem('user') || '{}');
        this.baseURL = '/api';
    }

    async fetchWithAuth(url, options = {}) {
        if (!this.token) {
            this.redirectToLogin();
            throw new Error('No token found');
        }

        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${this.token}`,
            'Content-Type': 'application/json'
        };

        const response = await fetch(url, options);

        if (response.status === 401) {
            this.logout();
            return null;
        }

        return response;
    }

    redirectToLogin() {
        window.location.href = '/login';
    }

    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        this.redirectToLogin();
    }

    formatCurrency(amount) {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0
        }).format(amount);
    }

    formatDate(dateString) {
        return new Date(dateString).toLocaleDateString('id-ID', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    }

    showToast(message, type = 'success') {
        // Simple toast implementation
        const toast = document.createElement('div');
        toast.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
        toast.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 1050;
            min-width: 300px;
        `;
        toast.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;

        document.body.appendChild(toast);

        setTimeout(() => {
            toast.remove();
        }, 5000);
    }

    async loadUserProfile() {
        try {
            const response = await this.fetchWithAuth(`${this.baseURL}/profile`);
            if (response) {
                const user = await response.json();
                localStorage.setItem('user', JSON.stringify(user));
                this.user = user;
                return user;
            }
        } catch (error) {
            console.error('Error loading profile:', error);
        }
        return null;
    }
}

// Initialize global app instance
window.app = new DormitoryApp();

// Navigation guard
document.addEventListener('DOMContentLoaded', function () {
    if (!window.app.token && !window.location.pathname.includes('/login') && window.location.pathname !== '/') {
        window.app.redirectToLogin();
    }
});
