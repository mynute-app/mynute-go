import { signal, computed, Signal, ReadonlySignal } from '@preact/signals';
import { api } from '../utils/api.ts';
import type { User, LoginResponse } from '../types.ts';

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
        const response = await api.post<LoginResponse>('/admin/auth/login', { email, password });
        
        token.value = response.token;
        user.value = response.user;
        
        localStorage.setItem('admin_token', response.token);
        
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
        return;
    }

    loading.value = true;
    try {
        const response = await api.get<{ user: User }>('/admin/auth/me');
        user.value = response.user;
    } catch (error) {
        // Token is invalid
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
