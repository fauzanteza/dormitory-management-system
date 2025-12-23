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
        if (!this.token) return null;
        try {
            const response = await this.fetchWithAuth('/api/profile');
            if (response.ok) {
                const user = await response.json();
                this.user = user;
                localStorage.setItem('user', JSON.stringify(user));
                this.updateUI();
                this.updateSidebar();
                this.checkAccess(); // Check page access
                return user;
            }
        } catch (error) {
            console.error('Error loading profile:', error);
        }
        return null;
    }

    updateUI() {
        // Update common UI elements if they exist
        const nameEl = document.getElementById('userName');
        const roleEl = document.getElementById('userRole');
        if (nameEl && this.user) nameEl.textContent = this.user.name;
        if (roleEl && this.user) roleEl.textContent = this.user.role === 'admin' ? 'Administrator' : 'Mahasiswa';
    }

    updateSidebar() {
        if (!this.user) return;

        const role = this.user.role;

        // Admin Specifics
        if (role === 'admin') {
            const reportBtn = document.getElementById('menu-reports');
            if (reportBtn) reportBtn.style.display = 'block';
        }

        // Student Specifics
        if (role !== 'admin') {
            // Hide Admin Menus
            const adminMenus = ['menu-residents', 'menu-users', 'menu-reports'];
            adminMenus.forEach(id => {
                const el = document.getElementById(id);
                if (el) el.style.display = 'none';
            });

            // Rename 'Kamar' to 'Kamar Saya'
            const kamarLabel = document.getElementById('label-rooms');
            if (kamarLabel) kamarLabel.textContent = 'Kamar Saya';
        }
    }

    checkAccess() {
        const path = window.location.pathname;
        if (this.user.role !== 'admin') {
            // Block access to admin-only pages
            if (path.includes('/residents') || path.includes('/users')) {
                alert('Akses ditolak. Halaman ini hanya untuk Administrator.');
                window.location.href = '/dashboard';
            }
        }
    }
}

// Initialize global app instance
const app = new DormitoryApp();
window.app = app; // Expose to window

// Navigation guard
document.addEventListener('DOMContentLoaded', function () {
    if (!window.app.token && !window.location.pathname.includes('/login') && window.location.pathname !== '/') {
        window.app.redirectToLogin();
    }
});
