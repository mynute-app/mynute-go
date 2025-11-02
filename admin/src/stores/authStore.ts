import { signal, computed, Signal, ReadonlySignal } from '@preact/signals';
import type { User } from '../types.ts';

// Signals
const storedUser = localStorage.getItem('admin_user');
const user: Signal<User | null> = signal(storedUser ? JSON.parse(storedUser) : null);
const token: Signal<string | null> = signal(localStorage.getItem('admin_token') || null);
const loading: Signal<boolean> = signal(false);

// Computed
const isAuthenticated: ReadonlySignal<boolean> = computed(() => !!token.value && !!user.value);

// Actions
async function login(email: string, password: string): Promise<{ success: boolean; error?: string }> {
    loading.value = true;
    try {
        // Admin login returns token in X-Auth-Token header, not in response body
        const response = await fetch('/api/admin/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Login failed');
        }

        // Extract token from response header
        const authToken = response.headers.get('X-Auth-Token');
        if (!authToken) {
            throw new Error('No auth token received');
        }

        token.value = authToken;
        localStorage.setItem('admin_token', authToken);

        // Fetch admin data by email using public endpoint
        try {
            const adminResponse = await fetch(`/api/admin/email/${email}`, {
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (adminResponse.ok) {
                const adminData = await adminResponse.json();
                const userData = (adminData.data || adminData) as User;
                user.value = userData;
                localStorage.setItem('admin_user', JSON.stringify(userData));
            } else {
                // Couldn't fetch user data, but login succeeded
                const userData = { email } as User;
                user.value = userData;
                localStorage.setItem('admin_user', JSON.stringify(userData));
            }
        } catch (err) {
            // Couldn't fetch user data, but login succeeded
            const userData = { email } as User;
            user.value = userData;
            localStorage.setItem('admin_user', JSON.stringify(userData));
        }
        
        return { success: true };
    } catch (error) {
        const message = error instanceof Error ? error.message : 'Login failed';
        return { success: false, error: message };
    } finally {
        loading.value = false;
    }
}

async function logout(): Promise<void> {
    token.value = null;
    user.value = null;
    localStorage.removeItem('admin_token');
    localStorage.removeItem('admin_user');
}

async function checkAuth(): Promise<void> {
    if (!token.value) {
        loading.value = false;
        user.value = null;
        return;
    }

    // If we already have user data from localStorage, we're authenticated
    // No need to make API calls on every page load
    if (user.value && user.value.email && user.value.id) {
        loading.value = false;
        return;
    }

    // We have a token but no user data - shouldn't happen with localStorage, but clear it
    await logout();
}

// Export store
export const authStore = {
    // State
    user,
    token,
    loading,
    isAuthenticated,
    
    // Actions
    login,
    logout,
    checkAuth
};
