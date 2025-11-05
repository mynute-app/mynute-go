// Type definitions for the application

export interface User {
    id: string;
    name: string;
    email: string;
    role?: string;
    createdAt?: string;
    updatedAt?: string;
}

export interface Admin extends User {
    permissions?: string[];
}

export interface LoginCredentials {
    email: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    user: User;
}

export interface ApiResponse<T = any> {
    success: boolean;
    data?: T;
    error?: string;
    message?: string;
}

export interface AuthStore {
    user: Signal<User | null>;
    token: Signal<string | null>;
    loading: Signal<boolean>;
    isAuthenticated: ReadonlySignal<boolean>;
    login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
    logout: () => Promise<void>;
    checkAuth: () => Promise<void>;
}

export interface AdminStore {
    admins: Signal<Admin[]>;
    loading: Signal<boolean>;
    error: Signal<string | null>;
    fetchAdmins: () => Promise<void>;
    createAdmin: (adminData: Partial<Admin>) => Promise<{ success: boolean; data?: Admin; error?: string }>;
    updateAdmin: (id: string, adminData: Partial<Admin>) => Promise<{ success: boolean; error?: string }>;
    deleteAdmin: (id: string) => Promise<{ success: boolean; error?: string }>;
}

// Signal types from @preact/signals
export type Signal<T> = {
    value: T;
};

export type ReadonlySignal<T> = {
    readonly value: T;
};
