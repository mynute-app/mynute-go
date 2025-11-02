import { html } from 'htm/preact';
import { useEffect, useState } from 'preact/hooks';
import { route } from 'preact-router';
import { companyStore, Company } from '../stores/companyStore.ts';

export default function Companies() {
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        companyStore.fetchCompanies();
    }, []);

    const handleDelete = async (id: string, name: string) => {
        if (!confirm(`Are you sure you want to delete "${name}"? This will delete all associated data (branches, employees, services, etc.).`)) return;
        
        const result = await companyStore.deleteCompany(id);
        if (!result.success) {
            alert('Failed to delete company: ' + result.error);
        }
    };

    const handleView = (id: string) => {
        route(`/companies/${id}`);
    };

    const filteredCompanies = companyStore.companies.value.filter((company: Company) =>
        company.legal_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        company.trade_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        company.tax_id?.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return html`
        <div>
            <div class="flex justify-between items-center mb-8">
                <div>
                    <h1 class="text-3xl font-bold text-gray-900">Companies</h1>
                    <p class="text-gray-600 mt-1">Manage system tenants and their data</p>
                </div>
            </div>

            <!-- Search Bar -->
            <div class="mb-6">
                <input
                    type="text"
                    placeholder="Search by name or tax ID..."
                    value=${searchTerm}
                    onInput=${(e: Event) => setSearchTerm((e.target as HTMLInputElement).value)}
                    class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                />
            </div>
            
            ${companyStore.loading.value ? html`
                <div class="text-center py-12">
                    <div class="text-xl text-gray-600">Loading companies...</div>
                </div>
            ` : companyStore.error.value ? html`
                <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                    Error: ${companyStore.error.value}
                </div>
            ` : html`
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    ${filteredCompanies.length === 0 ? html`
                        <div class="col-span-full text-center py-12 text-gray-500">
                            ${searchTerm ? 'No companies found matching your search.' : 'No companies registered yet.'}
                        </div>
                    ` : filteredCompanies.map((company: Company) => html`
                        <div key=${company.id} class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow">
                            <div class="p-6">
                                <div class="flex items-start justify-between mb-4">
                                    <div class="flex-1">
                                        <h3 class="text-lg font-semibold text-gray-900 mb-1">
                                            ${company.trade_name || company.legal_name}
                                        </h3>
                                        ${company.trade_name && company.trade_name !== company.legal_name ? html`
                                            <p class="text-sm text-gray-500">${company.legal_name}</p>
                                        ` : null}
                                    </div>
                                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                        Active
                                    </span>
                                </div>
                                
                                <div class="space-y-2 mb-4">
                                    <div class="flex items-center text-sm text-gray-600">
                                        <span class="mr-2">🏢</span>
                                        <span>Tax ID: ${company.tax_id}</span>
                                    </div>
                                    <div class="flex items-center text-sm text-gray-600">
                                        <span class="mr-2">📅</span>
                                        <span>Created: ${new Date(company.created_at).toLocaleDateString()}</span>
                                    </div>
                                </div>

                                <div class="flex gap-2">
                                    <button
                                        onClick=${() => handleView(company.id)}
                                        class="flex-1 bg-primary text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition-colors text-sm"
                                    >
                                        View Details
                                    </button>
                                    <button
                                        onClick=${() => handleDelete(company.id, company.trade_name || company.legal_name)}
                                        class="px-4 py-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition-colors text-sm"
                                    >
                                        Delete
                                    </button>
                                </div>
                            </div>
                        </div>
                    `)}
                </div>
            `}
        </div>
    `;
}
