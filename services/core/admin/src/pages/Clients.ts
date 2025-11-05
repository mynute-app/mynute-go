import { html } from 'htm/preact';
import { useEffect, useState } from 'preact/hooks';
import { clientStore, Client } from '../stores/clientStore.ts';

export default function Clients() {
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedClientId, setSelectedClientId] = useState<string | null>(null);

    useEffect(() => {
        clientStore.fetchClients();
    }, []);

    useEffect(() => {
        if (selectedClientId) {
            clientStore.fetchClientAppointments(selectedClientId);
        }
    }, [selectedClientId]);

    const handleDelete = async (id: string, name: string) => {
        if (!confirm(`Are you sure you want to delete client "${name}"?`)) return;
        
        const result = await clientStore.deleteClient(id);
        if (!result.success) {
            alert('Failed to delete client: ' + result.error);
        }
    };

    const filteredClients = clientStore.clients.value.filter((client: Client) =>
        client.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        client.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        client.phone?.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return html`
        <div>
            <div class="mb-8">
                <h1 class="text-3xl font-bold text-gray-900">Clients</h1>
                <p class="text-gray-600 mt-1">Manage registered clients across all companies</p>
            </div>

            <!-- Search Bar -->
            <div class="mb-6">
                <input
                    type="text"
                    placeholder="Search by name, email, or phone..."
                    value=${searchTerm}
                    onInput=${(e: Event) => setSearchTerm((e.target as HTMLInputElement).value)}
                    class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                />
            </div>
            
            ${clientStore.loading.value ? html`
                <div class="text-center py-12">
                    <div class="text-xl text-gray-600">Loading clients...</div>
                </div>
            ` : clientStore.error.value ? html`
                <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                    Error: ${clientStore.error.value}
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
                                    Phone
                                </th>
                                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Registered
                                </th>
                                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            ${filteredClients.length === 0 ? html`
                                <tr>
                                    <td colspan="5" class="px-6 py-12 text-center text-gray-500">
                                        ${searchTerm ? 'No clients found matching your search.' : 'No clients registered yet.'}
                                    </td>
                                </tr>
                            ` : filteredClients.map((client: Client) => html`
                                <tr key=${client.id} class="hover:bg-gray-50">
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="flex items-center">
                                            <div class="flex-shrink-0 h-10 w-10 rounded-full bg-primary flex items-center justify-center text-white font-semibold">
                                                ${client.name.charAt(0).toUpperCase()}
                                            </div>
                                            <div class="ml-4">
                                                <div class="text-sm font-medium text-gray-900">
                                                    ${client.name} ${client.surname || ''}
                                                </div>
                                            </div>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="text-sm text-gray-500">${client.email}</div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="text-sm text-gray-500">${client.phone || '-'}</div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="text-sm text-gray-500">
                                            ${new Date(client.created_at).toLocaleDateString()}
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                                        <button
                                            onClick=${() => setSelectedClientId(client.id)}
                                            class="text-primary hover:text-blue-900"
                                        >
                                            View Details
                                        </button>
                                        <button
                                            onClick=${() => handleDelete(client.id, client.name)}
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

            <!-- Client Details Modal -->
            ${selectedClientId ? html`
                <${ClientDetailsModal} 
                    clientId=${selectedClientId} 
                    onClose=${() => setSelectedClientId(null)} 
                />
            ` : null}
        </div>
    `;
}

// Client Details Modal
function ClientDetailsModal({ clientId, onClose }: { clientId: string; onClose: () => void }) {
    const client = clientStore.clients.value.find((c: Client) => c.id === clientId);
    const appointments = clientStore.clientAppointments.value;

    if (!client) return null;

    return html`
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick=${onClose}>
            <div class="bg-white rounded-lg shadow-xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto" onClick=${(e: Event) => e.stopPropagation()}>
                <div class="p-6 border-b border-gray-200">
                    <div class="flex justify-between items-start">
                        <div>
                            <h2 class="text-2xl font-bold text-gray-900">
                                ${client.name} ${client.surname || ''}
                            </h2>
                            <p class="text-gray-600 mt-1">${client.email}</p>
                        </div>
                        <button onClick=${onClose} class="text-gray-400 hover:text-gray-600">
                            <span class="text-2xl">×</span>
                        </button>
                    </div>
                </div>

                <div class="p-6">
                    <!-- Client Info -->
                    <div class="mb-6">
                        <h3 class="text-lg font-semibold text-gray-900 mb-3">Client Information</h3>
                        <div class="grid grid-cols-2 gap-4">
                            <div>
                                <dt class="text-sm font-medium text-gray-500">Email</dt>
                                <dd class="mt-1 text-sm text-gray-900">${client.email}</dd>
                            </div>
                            <div>
                                <dt class="text-sm font-medium text-gray-500">Phone</dt>
                                <dd class="mt-1 text-sm text-gray-900">${client.phone || '-'}</dd>
                            </div>
                            <div>
                                <dt class="text-sm font-medium text-gray-500">Registered</dt>
                                <dd class="mt-1 text-sm text-gray-900">${new Date(client.created_at).toLocaleString()}</dd>
                            </div>
                            <div>
                                <dt class="text-sm font-medium text-gray-500">Last Updated</dt>
                                <dd class="mt-1 text-sm text-gray-900">${new Date(client.updated_at).toLocaleString()}</dd>
                            </div>
                        </div>
                    </div>

                    <!-- Appointments -->
                    <div>
                        <h3 class="text-lg font-semibold text-gray-900 mb-3">
                            Appointments (${appointments.length})
                        </h3>
                        ${appointments.length === 0 ? html`
                            <p class="text-gray-500 text-center py-8">No appointments found</p>
                        ` : html`
                            <div class="space-y-3">
                                ${appointments.map((apt: any) => html`
                                    <div key=${apt.id} class="border border-gray-200 rounded-lg p-4">
                                        <div class="flex justify-between items-start">
                                            <div>
                                                <div class="font-medium text-gray-900">
                                                    ${new Date(apt.start_time).toLocaleString()}
                                                </div>
                                                <div class="text-sm text-gray-500 mt-1">
                                                    Status: ${apt.status}
                                                </div>
                                            </div>
                                            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                                                apt.status === 'confirmed' ? 'bg-green-100 text-green-800' :
                                                apt.status === 'cancelled' ? 'bg-red-100 text-red-800' :
                                                'bg-gray-100 text-gray-800'
                                            }">
                                                ${apt.status}
                                            </span>
                                        </div>
                                    </div>
                                `)}
                            </div>
                        `}
                    </div>
                </div>

                <div class="p-6 border-t border-gray-200 flex justify-end">
                    <button
                        onClick=${onClose}
                        class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
                    >
                        Close
                    </button>
                </div>
            </div>
        </div>
    `;
}
