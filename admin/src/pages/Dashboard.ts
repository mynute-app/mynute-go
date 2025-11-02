import { html } from 'htm/preact';
import { useEffect } from 'preact/hooks';
import { route } from 'preact-router';
import { adminStore } from '../stores/adminStore.ts';
import { companyStore } from '../stores/companyStore.ts';
import { clientStore } from '../stores/clientStore.ts';

interface Stat {
    label: string;
    value: string | number;
    color: string;
    icon: string;
    onClick?: () => void;
}

export default function Dashboard() {
    useEffect(() => {
        // Fetch stats on mount
        adminStore.fetchAdmins();
        companyStore.fetchCompanies();
        clientStore.fetchClients();
    }, []);

    const stats: Stat[] = [
        { 
            label: 'Total Companies', 
            value: companyStore.companies.value.length, 
            color: 'bg-purple-500',
            icon: '🏢',
            onClick: () => route('/admin/companies')
        },
        { 
            label: 'Total Clients', 
            value: clientStore.clients.value.length, 
            color: 'bg-blue-500',
            icon: '👥',
            onClick: () => route('/admin/clients')
        },
        { 
            label: 'Admin Users', 
            value: adminStore.admins.value.length, 
            color: 'bg-green-500',
            icon: '🔐',
            onClick: () => route('/admin/users')
        },
        { 
            label: 'System Status', 
            value: 'Active', 
            color: 'bg-emerald-500',
            icon: '✓'
        },
    ];

    return html`
        <div>
            <div class="mb-8">
                <h1 class="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
                <p class="text-gray-600 mt-2">Manage companies, clients, and system settings</p>
            </div>
            
            <!-- Stats Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                ${stats.map(stat => html`
                    <div 
                        key=${stat.label} 
                        class="bg-white rounded-lg shadow p-6 ${stat.onClick ? 'cursor-pointer hover:shadow-lg transition-shadow' : ''}"
                        onClick=${stat.onClick || null}
                    >
                        <div class="flex items-center">
                            <div class="${stat.color} w-14 h-14 rounded-lg flex items-center justify-center text-white text-2xl mr-4">
                                ${stat.icon}
                            </div>
                            <div>
                                <p class="text-sm text-gray-600 mb-1">${stat.label}</p>
                                <p class="text-3xl font-bold text-gray-900">${stat.value}</p>
                            </div>
                        </div>
                    </div>
                `)}
            </div>
            
            <!-- Quick Actions -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                <div class="bg-white rounded-lg shadow p-6">
                    <h2 class="text-xl font-semibold text-gray-900 mb-4">Quick Actions</h2>
                    <div class="space-y-3">
                        <button
                            onClick=${() => route('/admin/companies')}
                            class="w-full text-left px-4 py-3 rounded-lg bg-purple-50 hover:bg-purple-100 transition-colors"
                        >
                            <div class="font-medium text-purple-900">View All Companies</div>
                            <div class="text-sm text-purple-600">Manage tenants and their data</div>
                        </button>
                        <button
                            onClick=${() => route('/admin/clients')}
                            class="w-full text-left px-4 py-3 rounded-lg bg-blue-50 hover:bg-blue-100 transition-colors"
                        >
                            <div class="font-medium text-blue-900">View All Clients</div>
                            <div class="text-sm text-blue-600">Browse registered clients</div>
                        </button>
                        <button
                            onClick=${() => route('/admin/users')}
                            class="w-full text-left px-4 py-3 rounded-lg bg-green-50 hover:bg-green-100 transition-colors"
                        >
                            <div class="font-medium text-green-900">Manage Admins</div>
                            <div class="text-sm text-green-600">Add or remove admin users</div>
                        </button>
                    </div>
                </div>
                
                <!-- System Info -->
                <div class="bg-white rounded-lg shadow p-6">
                    <h2 class="text-xl font-semibold text-gray-900 mb-4">System Information</h2>
                    <div class="space-y-3">
                        <div class="flex justify-between items-center py-2 border-b">
                            <span class="text-gray-600">Database Status</span>
                            <span class="text-green-600 font-medium">Connected</span>
                        </div>
                        <div class="flex justify-between items-center py-2 border-b">
                            <span class="text-gray-600">API Version</span>
                            <span class="text-gray-900 font-medium">v1.0.0</span>
                        </div>
                        <div class="flex justify-between items-center py-2 border-b">
                            <span class="text-gray-600">Environment</span>
                            <span class="text-gray-900 font-medium">Production</span>
                        </div>
                        <div class="flex justify-between items-center py-2">
                            <span class="text-gray-600">Uptime</span>
                            <span class="text-gray-900 font-medium">99.9%</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
}
