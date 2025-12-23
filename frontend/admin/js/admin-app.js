// Admin Dashboard Application
class AdminDashboard {
    constructor() {
        this.API_BASE = '/api/admin';
    }

    async loadDashboardStats() {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/dashboard/stats`);
            if (!response || !response.ok) {
                throw new Error('Failed to load stats');
            }
            const stats = await response.json();
            return stats;
        } catch (error) {
            console.error('Error loading dashboard stats:', error);
            return null;
        }
    }

    async loadRooms(params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const url = `${this.API_BASE}/rooms${queryString ? '?' + queryString : ''}`;
            const response = await window.auth.fetchWithAuth(url);
            if (!response || !response.ok) throw new Error('Failed to load rooms');
            return await response.json();
        } catch (error) {
            console.error('Error loading rooms:', error);
            return [];
        }
    }

    async createRoom(roomData) {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/rooms`, {
                method: 'POST',
                body: JSON.stringify(roomData)
            });
            if (!response || !response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to create room');
            }
            return await response.json();
        } catch (error) {
            console.error('Error creating room:', error);
            throw error;
        }
    }

    async updateRoom(id, roomData) {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/rooms/${id}`, {
                method: 'PUT',
                body: JSON.stringify(roomData)
            });
            if (!response || !response.ok) throw new Error('Failed to update room');
            return await response.json();
        } catch (error) {
            console.error('Error updating room:', error);
            throw error;
        }
    }

    async deleteRoom(id) {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/rooms/${id}`, {
                method: 'DELETE'
            });
            if (!response || !response.ok) throw new Error('Failed to delete room');
            return await response.json();
        } catch (error) {
            console.error('Error deleting room:', error);
            throw error;
        }
    }

    async loadBookings(params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const url = `${this.API_BASE}/bookings${queryString ? '?' + queryString : ''}`;
            const response = await window.auth.fetchWithAuth(url);
            if (!response || !response.ok) throw new Error('Failed to load bookings');
            return await response.json();
        } catch (error) {
            console.error('Error loading bookings:', error);
            return [];
        }
    }

    async updateBookingStatus(id, status) {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/bookings/${id}/status`, {
                method: 'PUT',
                body: JSON.stringify({ status })
            });
            if (!response || !response.ok) throw new Error('Failed to update booking');
            return await response.json();
        } catch (error) {
            console.error('Error updating booking:', error);
            throw error;
        }
    }

    async loadResidents(params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const url = `${this.API_BASE}/residents${queryString ? '?' + queryString : ''}`;
            const response = await window.auth.fetchWithAuth(url);
            if (!response || !response.ok) throw new Error('Failed to load residents');
            return await response.json();
        } catch (error) {
            console.error('Error loading residents:', error);
            return [];
        }
    }

    async loadPayments(params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const url = `${this.API_BASE}/payments${queryString ? '?' + queryString : ''}`;
            const response = await window.auth.fetchWithAuth(url);
            if (!response || !response.ok) throw new Error('Failed to load payments');
            return await response.json();
        } catch (error) {
            console.error('Error loading payments:', error);
            return [];
        }
    }

    async updatePaymentStatus(id, status, paymentDate = null) {
        try {
            const body = { status };
            if (paymentDate) body.payment_date = paymentDate;

            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/payments/${id}/status`, {
                method: 'PUT',
                body: JSON.stringify(body)
            });
            if (!response || !response.ok) throw new Error('Failed to update payment');
            return await response.json();
        } catch (error) {
            console.error('Error updating payment:', error);
            throw error;
        }
    }

    async loadRepairs(params = {}) {
        try {
            const queryString = new URLSearchParams(params).toString();
            const url = `${this.API_BASE}/repairs${queryString ? '?' + queryString : ''}`;
            const response = await window.auth.fetchWithAuth(url);
            if (!response || !response.ok) throw new Error('Failed to load repairs');
            return await response.json();
        } catch (error) {
            console.error('Error loading repairs:', error);
            return [];
        }
    }

    async updateRepairStatus(id, status) {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/repairs/${id}/status`, {
                method: 'PUT',
                body: JSON.stringify({ status })
            });
            if (!response || !response.ok) throw new Error('Failed to update repair');
            return await response.json();
        } catch (error) {
            console.error('Error updating repair:', error);
            throw error;
        }
    }

    async loadUsers() {
        try {
            const response = await window.auth.fetchWithAuth(`${this.API_BASE}/users`);
            if (!response || !response.ok) throw new Error('Failed to load users');
            return await response.json();
        } catch (error) {
            console.error('Error loading users:', error);
            return [];
        }
    }
}

// Initialize global admin dashboard instance
window.adminDashboard = new AdminDashboard();

// Helper functions for status badges
function getStatusBadge(status) {
    const badges = {
        'pending': 'bg-warning text-dark',
        'approved': 'bg-success',
        'rejected': 'bg-danger',
        'paid': 'bg-success',
        'overdue': 'bg-danger',
        'available': 'bg-success',
        'occupied': 'bg-primary',
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
