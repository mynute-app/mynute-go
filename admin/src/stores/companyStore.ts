import { signal, Signal } from '@preact/signals';
import { api } from '../utils/api.ts';

// Company model
export interface Company {
  id: string;
  legal_name: string;
  trade_name: string;
  tax_id: string;
  created_at: string;
  updated_at: string;
  // Additional fields from GetFullCompany
  branches?: Branch[];
  employees?: Employee[];
  services?: Service[];
  subdomains?: Subdomain[];
}

export interface Branch {
  id: string;
  name: string;
  company_id: string;
  address?: string;
  phone?: string;
  created_at: string;
}

export interface Employee {
  id: string;
  name: string;
  surname: string;
  email: string;
  phone?: string;
  company_id: string;
  is_owner: boolean;
  created_at: string;
}

export interface Service {
  id: string;
  name: string;
  description?: string;
  duration: number;
  price: number;
  company_id: string;
  created_at: string;
}

export interface Subdomain {
  id: string;
  name: string;
  company_id: string;
  created_at: string;
}

// Companies state
const companies: Signal<Company[]> = signal([]);
const selectedCompany: Signal<Company | null> = signal(null);
const loading: Signal<boolean> = signal(false);
const error: Signal<string | null> = signal(null);

// Actions
async function fetchCompanies(): Promise<void> {
  loading.value = true;
  error.value = null;
  
  try {
    // Note: We need an endpoint to list all companies for admin
    // This might need to be created on the backend
    const response = await api.get<Company[]>('/admin/companies');
    companies.value = Array.isArray(response) ? response : (response as any).data || [];
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to fetch companies';
    error.value = message;
  } finally {
    loading.value = false;
  }
}

async function fetchCompanyById(id: string): Promise<void> {
  loading.value = true;
  error.value = null;
  
  try {
    const response = await api.get<Company>(`/company/${id}`);
    selectedCompany.value = (response as any).data || response;
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to fetch company';
    error.value = message;
  } finally {
    loading.value = false;
  }
}

async function createCompany(companyData: Partial<Company>): Promise<{ success: boolean; data?: Company; error?: string }> {
  try {
    const response = await api.post<Company>('/company', companyData);
    const newCompany = (response as any).data || response;
    companies.value = [...companies.value, newCompany];
    return { success: true, data: newCompany };
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to create company';
    return { success: false, error: message };
  }
}

async function updateCompany(id: string, companyData: Partial<Company>): Promise<{ success: boolean; error?: string }> {
  try {
    const response = await api.put<Company>(`/company/${id}`, companyData);
    const updatedCompany = (response as any).data || response;
    companies.value = companies.value.map((company: Company) => 
      company.id === id ? updatedCompany : company
    );
    if (selectedCompany.value?.id === id) {
      selectedCompany.value = updatedCompany;
    }
    return { success: true };
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to update company';
    return { success: false, error: message };
  }
}

async function deleteCompany(id: string): Promise<{ success: boolean; error?: string }> {
  try {
    await api.delete(`/company/${id}`);
    companies.value = companies.value.filter((company: Company) => company.id !== id);
    if (selectedCompany.value?.id === id) {
      selectedCompany.value = null;
    }
    return { success: true };
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to delete company';
    return { success: false, error: message };
  }
}

// Export store
export const companyStore = {
  // State
  companies,
  selectedCompany,
  loading,
  error,
  
  // Actions
  fetchCompanies,
  fetchCompanyById,
  createCompany,
  updateCompany,
  deleteCompany,
};
