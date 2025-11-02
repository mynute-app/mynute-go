import { html } from 'htm/preact';
import { useEffect, useState } from 'preact/hooks';
import { route } from 'preact-router';
import { companyStore, Company, Branch, Employee, Service } from '../stores/companyStore.ts';

interface CompanyDetailProps {
    id: string;
}

export default function CompanyDetail({ id }: CompanyDetailProps) {
    const [activeTab, setActiveTab] = useState('overview');

    useEffect(() => {
        companyStore.fetchCompanyById(id);
    }, [id]);

    const company = companyStore.selectedCompany.value;

    if (companyStore.loading.value) {
        return html`
            <div class="flex items-center justify-center py-12">
                <div class="text-xl text-gray-600">Loading company details...</div>
            </div>
        `;
    }

    if (!company) {
        return html`
            <div class="text-center py-12">
                <p class="text-gray-600 mb-4">Company not found</p>
                <button
                    onClick=${() => route('/companies')}
                    class="text-primary hover:underline"
                >
                    Back to Companies
                </button>
            </div>
        `;
    }

    const tabs = [
        { id: 'overview', label: 'Overview', icon: '📋' },
        { id: 'branches', label: 'Branches', icon: '🏪', count: company.branches?.length || 0 },
        { id: 'employees', label: 'Employees', icon: '👥', count: company.employees?.length || 0 },
        { id: 'services', label: 'Services', icon: '🛎️', count: company.services?.length || 0 },
        { id: 'subdomains', label: 'Subdomains', icon: '🌐', count: company.subdomains?.length || 0 },
    ];

    return html`
        <div>
            <!-- Header -->
            <div class="mb-6">
                <button
                    onClick=${() => route('/companies')}
                    class="text-primary hover:underline mb-4 flex items-center"
                >
                    ← Back to Companies
                </button>
                <div class="flex justify-between items-start">
                    <div>
                        <h1 class="text-3xl font-bold text-gray-900">
                            ${company.trade_name || company.legal_name}
                        </h1>
                        ${company.trade_name && company.trade_name !== company.legal_name ? html`
                            <p class="text-gray-600 mt-1">${company.legal_name}</p>
                        ` : null}
                        <p class="text-sm text-gray-500 mt-2">Tax ID: ${company.tax_id}</p>
                    </div>
                    <button
                        onClick=${async () => {
                            if (!confirm('Are you sure you want to delete this company and all associated data?')) return;
                            const result = await companyStore.deleteCompany(company.id);
                            if (result.success) {
                                route('/companies');
                            } else {
                                alert('Failed to delete: ' + result.error);
                            }
                        }}
                        class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
                    >
                        Delete Company
                    </button>
                </div>
            </div>

            <!-- Tabs -->
            <div class="border-b border-gray-200 mb-6">
                <nav class="-mb-px flex space-x-8">
                    ${tabs.map(tab => html`
                        <button
                            key=${tab.id}
                            onClick=${() => setActiveTab(tab.id)}
                            class="${activeTab === tab.id 
                                ? 'border-primary text-primary' 
                                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center"
                        >
                            <span class="mr-2">${tab.icon}</span>
                            ${tab.label}
                            ${tab.count !== undefined ? html`
                                <span class="ml-2 py-0.5 px-2 rounded-full text-xs bg-gray-100 text-gray-600">
                                    ${tab.count}
                                </span>
                            ` : null}
                        </button>
                    `)}
                </nav>
            </div>

            <!-- Tab Content -->
            <div>
                ${activeTab === 'overview' ? html`<${OverviewTab} company=${company} />` : null}
                ${activeTab === 'branches' ? html`<${BranchesTab} branches=${company.branches || []} />` : null}
                ${activeTab === 'employees' ? html`<${EmployeesTab} employees=${company.employees || []} />` : null}
                ${activeTab === 'services' ? html`<${ServicesTab} services=${company.services || []} />` : null}
                ${activeTab === 'subdomains' ? html`<${SubdomainsTab} subdomains=${company.subdomains || []} />` : null}
            </div>
        </div>
    `;
}

// Overview Tab
function OverviewTab({ company }: { company: Company }) {
    return html`
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="bg-white rounded-lg shadow p-6">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">Company Information</h3>
                <dl class="space-y-3">
                    <div>
                        <dt class="text-sm font-medium text-gray-500">Legal Name</dt>
                        <dd class="mt-1 text-sm text-gray-900">${company.legal_name}</dd>
                    </div>
                    <div>
                        <dt class="text-sm font-medium text-gray-500">Trade Name</dt>
                        <dd class="mt-1 text-sm text-gray-900">${company.trade_name}</dd>
                    </div>
                    <div>
                        <dt class="text-sm font-medium text-gray-500">Tax ID</dt>
                        <dd class="mt-1 text-sm text-gray-900">${company.tax_id}</dd>
                    </div>
                    <div>
                        <dt class="text-sm font-medium text-gray-500">Created</dt>
                        <dd class="mt-1 text-sm text-gray-900">${new Date(company.created_at).toLocaleString()}</dd>
                    </div>
                    <div>
                        <dt class="text-sm font-medium text-gray-500">Last Updated</dt>
                        <dd class="mt-1 text-sm text-gray-900">${new Date(company.updated_at).toLocaleString()}</dd>
                    </div>
                </dl>
            </div>

            <div class="bg-white rounded-lg shadow p-6">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">Statistics</h3>
                <div class="space-y-4">
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600">Total Branches</span>
                        <span class="text-2xl font-bold text-gray-900">${company.branches?.length || 0}</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600">Total Employees</span>
                        <span class="text-2xl font-bold text-gray-900">${company.employees?.length || 0}</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600">Total Services</span>
                        <span class="text-2xl font-bold text-gray-900">${company.services?.length || 0}</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600">Subdomains</span>
                        <span class="text-2xl font-bold text-gray-900">${company.subdomains?.length || 0}</span>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// Branches Tab
function BranchesTab({ branches }: { branches: Branch[] }) {
    return html`
        <div class="bg-white rounded-lg shadow overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Address</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Phone</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
                        <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    ${branches.length === 0 ? html`
                        <tr>
                            <td colspan="5" class="px-6 py-12 text-center text-gray-500">
                                No branches found.
                            </td>
                        </tr>
                    ` : branches.map((branch: Branch) => html`
                        <tr key=${branch.id}>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm font-medium text-gray-900">${branch.name}</div>
                            </td>
                            <td class="px-6 py-4">
                                <div class="text-sm text-gray-500">${branch.address || '-'}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm text-gray-500">${branch.phone || '-'}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm text-gray-500">${new Date(branch.created_at).toLocaleDateString()}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                <button class="text-red-600 hover:text-red-900">Delete</button>
                            </td>
                        </tr>
                    `)}
                </tbody>
            </table>
        </div>
    `;
}

// Employees Tab
function EmployeesTab({ employees }: { employees: Employee[] }) {
    return html`
        <div class="bg-white rounded-lg shadow overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Phone</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Role</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
                        <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    ${employees.length === 0 ? html`
                        <tr>
                            <td colspan="6" class="px-6 py-12 text-center text-gray-500">
                                No employees found.
                            </td>
                        </tr>
                    ` : employees.map((employee: Employee) => html`
                        <tr key=${employee.id}>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm font-medium text-gray-900">
                                    ${employee.name} ${employee.surname}
                                    ${employee.is_owner ? html`
                                        <span class="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800">
                                            Owner
                                        </span>
                                    ` : null}
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm text-gray-500">${employee.email}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm text-gray-500">${employee.phone || '-'}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${employee.is_owner ? 'bg-purple-100 text-purple-800' : 'bg-gray-100 text-gray-800'}">
                                    ${employee.is_owner ? 'Owner' : 'Employee'}
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm text-gray-500">${new Date(employee.created_at).toLocaleDateString()}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                <button class="text-red-600 hover:text-red-900" disabled=${employee.is_owner}>
                                    ${employee.is_owner ? 'Protected' : 'Delete'}
                                </button>
                            </td>
                        </tr>
                    `)}
                </tbody>
            </table>
        </div>
    `;
}

// Services Tab  
function ServicesTab({ services }: { services: Service[] }) {
    const formatDuration = (minutes: number) => {
        const hours = Math.floor(minutes / 60);
        const mins = minutes % 60;
        if (hours > 0) {
            return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
        }
        return `${mins}m`;
    };

    const formatPrice = (price: number) => {
        return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(price);
    };

    return html`
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            ${services.length === 0 ? html`
                <div class="col-span-full text-center py-12 text-gray-500">
                    No services found.
                </div>
            ` : services.map((service: Service) => html`
                <div key=${service.id} class="bg-white rounded-lg shadow p-4">
                    <div class="flex justify-between items-start mb-3">
                        <h4 class="text-lg font-semibold text-gray-900">${service.name}</h4>
                        <button class="text-red-600 hover:text-red-900 text-sm">Delete</button>
                    </div>
                    ${service.description ? html`
                        <p class="text-sm text-gray-600 mb-3">${service.description}</p>
                    ` : null}
                    <div class="flex justify-between items-center text-sm">
                        <span class="text-gray-600">⏱️ ${formatDuration(service.duration)}</span>
                        <span class="text-lg font-bold text-primary">${formatPrice(service.price)}</span>
                    </div>
                </div>
            `)}
        </div>
    `;
}

// Subdomains Tab
function SubdomainsTab({ subdomains }: { subdomains: any[] }) {
    return html`
        <div class="bg-white rounded-lg shadow overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Subdomain</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Full URL</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
                        <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    ${subdomains.length === 0 ? html`
                        <tr>
                            <td colspan="4" class="px-6 py-12 text-center text-gray-500">
                                No subdomains found.
                            </td>
                        </tr>
                    ` : subdomains.map((subdomain: any) => html`
                        <tr key=${subdomain.id}>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm font-medium text-gray-900">${subdomain.name}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <a href="https://${subdomain.name}.mynute.com" target="_blank" class="text-sm text-primary hover:underline">
                                    ${subdomain.name}.mynute.com
                                </a>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="text-sm text-gray-500">${new Date(subdomain.created_at).toLocaleDateString()}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                <button class="text-red-600 hover:text-red-900">Delete</button>
                            </td>
                        </tr>
                    `)}
                </tbody>
            </table>
        </div>
    `;
}
