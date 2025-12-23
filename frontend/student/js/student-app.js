// Student Dashboard Application
class StudentDashboard {
    constructor() {
        this.API_BASE = '/api/student';
    }

    async loadDashboardStats() {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/dashboard/stats`);
            if (!response || !response.ok) throw new Error('Failed to load stats');
            return await response.json();
        } catch (error) {
            console.error('Error loading dashboard stats:', error);
            return null;
        }
    }

    async loadMyRoom() {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/my-room`);
            if (!response || !response.ok) {
                if (response.status === 404) return null;
                throw new Error('Failed to load room');
            }
            return await response.json();
        } catch (error) {
            console.error('Error loading room:', error);
            return null;
        }
    }

    async loadMyPayments(params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const url = `${this.API_BASE}/my-payments${queryString ? '?' + queryString : ''}`;
            const response = await window.auth.fetchWithAuth(url);
            if (!response || !response.ok) throw new Error('Failed to load payments');
            return await response.json();
        } catch (error) {
            console.error('Error loading payments:', error);
            return [];
        }
    }

    async loadMyRepairs(params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const url = `${this.API_BASE}/my-repairs${queryString ? '?' + queryString : ''}`;
            const response = await window.auth.fetchWithAuth(url);
            if (!response || !response.ok) throw new Error('Failed to load repairs');
            return await response.json();
        } catch (error) {
            console.error('Error loading repairs:', error);
            return [];
        }
    }

    async createRepair(repairData) {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/repairs`, {
                method: 'POST',
                body: JSON.stringify(repairData)
            });
            if (!response || !response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to create repair request');
            }
            return await response.json();
        } catch (error) {
            console.error('Error creating repair:', error);
            throw error;
        }
    }

    async loadAvailableRooms() {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/rooms/available`);
            if (!response || !response.ok) throw new Error('Failed to load rooms');
            return await response.json();
        } catch (error) {
            console.error('Error loading rooms:', error);
            return [];
        }
    }

    async loadMyBookings() {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/bookings`);
            if (!response || !response.ok) throw new Error('Failed to load bookings');
            return await response.json();
        } catch (error) {
            console.error('Error loading bookings:', error);
            return [];
        }
    }

    async createBooking(bookingData) {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/bookings`, {
                method: 'POST',
                body: JSON.stringify(bookingData)
            });
            if (!response || !response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to create booking');
            }
            return await response.json();
        } catch (error) {
            console.error('Error creating booking:', error);
            throw error;
        }
    }
}

// Initialize global student dashboard instance
window.studentDashboard = new StudentDashboard();

// Helper functions
function getStatusBadge(status) {
    const badges = {
        'pending': 'bg-warning text-dark',
        'approved': 'bg-success',
        'rejected': 'bg-danger',
        'paid': 'bg-success',
        'overdue': 'bg-danger',
        'in_progress': 'bg-info',
        'completed': 'bg-success',
        'cancelled': 'bg-secondary'
    };
    return badges[status] || 'bg-secondary';
}

function formatCurrency(amount) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR'
    }).format(amount);
}

function formatDate(dateString) {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('id-ID', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}
