# Quick Fix - Frontend Data Not Showing

## Problem

Data exists in database but not showing in admin panel.

## Root Cause

Init functions in page HTML files need to be:

1. Made global (attached to window object)
2. Auto-executed when page loads

## Solution Applied

Fixed `rooms.html` by adding at the end of script:

```javascript
// Make functions global
window.initRoomsPage = initRoomsPage;
window.loadRooms = loadRooms;

// Auto-run when page loads
if (document.getElementById('roomsTableBody')) {
    initRoomsPage();
}
```

## How to Test

1. **Open browser** → Frontend admin panel
2. **Login** dengan `fauzan@example.com` / `password123`  
3. **Click "Manajemen Kamar"** di sidebar
4. **Data should appear** from database!

## If Still Not Working

**Check Browser Console** (F12):

- Look for errors
- Check network tab for API calls
- Verify token is sent in headers

**Verify API Response**:

```
GET http://localhost:8080/api/admin/rooms
Headers: Authorization: Bearer <token>
```

Should return array of rooms from database.

## Apply Same Fix to Other Pages

Need to add same pattern to:

- `dashboard.html`
- `residents.html`
- `payments.html`  
- `bookings.html`
- `repairs.html`
- `users.html`

Pattern:

```javascript
window.initPageName = initPageName;
if (document.getElementById('uniqueElementId')) {
    initPageName();
}
```
