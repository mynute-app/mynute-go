import { html } from 'htm/preact';
import { useEffect } from 'preact/hooks';
import { adminStore } from '../stores/adminStore.ts';
import type { Admin } from '../types.ts';

export default function Users() {
    useEffect(() => {
        adminStore.fetchAdmins();
    }, []);

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this admin?')) return;
        
        const result = await adminStore.deleteAdmin(id);
        if (!result.success) {
            alert('Failed to delete admin: ' + result.error);
        }
    };

    return html`
        <div>
            <div class="flex justify-between items-center mb-8">
                <h1 class="text-3xl font-bold text-gray-900">Admin Users</h1>
            </div>
            
            ${adminStore.loading.value ? html`
                <div class="text-center py-12">
                    <div class="text-xl text-gray-600">Loading...</div>
                </div>
            ` : adminStore.error.value ? html`
                <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                    Error: ${adminStore.error.value}
                </div>
            ` : html`
                <div class="bg-white rounded-lg shadow overflow-hidden">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-gray-50">
                            <tr>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Name
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Email
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Role
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Status
                                </th>
                                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            ${adminStore.admins.value.length === 0 ? html`
                                <tr>
                                    <td colspan="5" class="px-6 py-12 text-center text-gray-500">
                                        No admins found. Click "Add Admin" to create one.
                                    </td>
                                </tr>
                            ` : adminStore.admins.value.map((admin: Admin) => html`
                                <tr key=${admin.id}>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="text-sm font-medium text-gray-900">
                                            ${admin.name || 'N/A'}
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="text-sm text-gray-500">${admin.email}</div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-purple-100 text-purple-800">
                                            ${admin.role || 'Admin'}
                                        </span>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                                            Active
                                        </span>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                        <button
                                            onClick=${() => handleDelete(admin.id)}
                                            class="text-red-600 hover:text-red-900"
                                        >
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                            `)}
                        </tbody>
                    </table>
                </div>
            `}
        </div>
    `;
}
