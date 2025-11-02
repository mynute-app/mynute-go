import { signal, Signal } from '@preact/signals';
import { api } from '../utils/api.ts';
import type { Admin } from '../types.ts';

// Admin users state
const admins: Signal<Admin[]> = signal([]);
const loading: Signal<boolean> = signal(false);
const error: Signal<string | null> = signal(null);

// Actions
async function fetchAdmins(): Promise<void> {
    loading.value = true;
    error.value = null;
    
    try {
        const response = await api.get<Admin[]>('/admin/users');
        admins.value = Array.isArray(response) ? response : (response as any).data || [];
    } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to fetch admins';
        error.value = message;
    } finally {
        loading.value = false;
    }
}

async function createAdmin(adminData: Partial<Admin>): Promise<{ success: boolean; data?: Admin; error?: string }> {
    try {
        const response = await api.post<Admin>('/admin/users', adminData);
        const newAdmin = (response as any).data || response;
        admins.value = [...admins.value, newAdmin];
        return { success: true, data: newAdmin };
    } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to create admin';
        return { success: false, error: message };
    }
}

async function updateAdmin(id: string, adminData: Partial<Admin>): Promise<{ success: boolean; error?: string }> {
    try {
        const response = await api.put<Admin>(`/admin/users/${id}`, adminData);
        const updatedAdmin = (response as any).data || response;
        admins.value = admins.value.map((admin: Admin) => 
            admin.id === id ? updatedAdmin : admin
        );
        return { success: true };
    } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to update admin';
        return { success: false, error: message };
    }
}

async function deleteAdmin(id: string): Promise<{ success: boolean; error?: string }> {
    try {
        await api.delete(`/admin/users/${id}`);
        admins.value = admins.value.filter((admin: Admin) => admin.id !== id);
        return { success: true };
    } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to delete admin';
        return { success: false, error: message };
    }
}

// Export store
export const adminStore = {
    // State
    admins,
    loading,
    error,
    
    // Actions
    fetchAdmins,
    createAdmin,
    updateAdmin,
    deleteAdmin
};
