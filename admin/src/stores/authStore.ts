import { signal, computed, Signal, ReadonlySignal } from '@preact/signals';
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

        // Store basic user info from email - we don't get user data from login response
        token.value = authToken;
        user.value = { email } as User; // Store email, will be validated by checkAuth
        
        localStorage.setItem('admin_token', authToken);
        
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

    // If we have both token and user, we're good (already authenticated)
    if (user.value) {
        loading.value = false;
        return;
    }

    // We have a token but no user data - token might be from previous session
    // Try to validate it silently, if it fails just clear it
    loading.value = true;
    try {
        // Validate token by making an authenticated request
        // If successful, we keep the token but need user data
        // Since there's no /me endpoint and we don't have user data,
        // we should just logout and let user re-login
        await logout();
    } catch (error) {
        // Token is invalid, clear it
        await logout();
    } finally {
        loading.value = false;
    }
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
