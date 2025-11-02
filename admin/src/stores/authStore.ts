import { signal, computed, Signal, ReadonlySignal } from '@preact/signals';
import { api } from '../utils/api.ts';
import type { User } from '../types.ts';

// Signals
const user: Signal<User | null> = signal(null);
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
                user.value = (adminData.data || adminData) as User;
            } else {
                // Couldn't fetch user data, but login succeeded
                user.value = { email } as User;
            }
        } catch (err) {
            // Couldn't fetch user data, but login succeeded
            user.value = { email } as User;
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
}

async function checkAuth(): Promise<void> {
    if (!token.value) {
        loading.value = false;
        user.value = null;
        return;
    }

    // If we already have user data, just validate the token is still good
    if (user.value && user.value.email) {
        loading.value = true;
        try {
            // Re-fetch user data to validate token and refresh user info
            const response = await api.get<any>(`/admin/email/${user.value.email}`);
            user.value = (response.data || response) as User;
        } catch (error) {
            // Token is invalid
            await logout();
        } finally {
            loading.value = false;
        }
        return;
    }

    // We have a token but no user - shouldn't happen, but clear it
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
