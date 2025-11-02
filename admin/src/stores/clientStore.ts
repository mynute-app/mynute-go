import { signal, Signal } from '@preact/signals';
import { api } from '../utils/api.ts';

// Client model
export interface Client {
  id: string;
  name: string;
  surname?: string;
  email: string;
  phone?: string;
  created_at: string;
  updated_at: string;
}

export interface ClientAppointment {
  id: string;
  client_id: string;
  employee_id: string;
  service_id: string;
  branch_id: string;
  start_time: string;
  end_time: string;
  status: string;
  created_at: string;
}

// Clients state
const clients: Signal<Client[]> = signal([]);
const selectedClient: Signal<Client | null> = signal(null);
const clientAppointments: Signal<ClientAppointment[]> = signal([]);
const loading: Signal<boolean> = signal(false);
const error: Signal<string | null> = signal(null);

// Actions
async function fetchClients(): Promise<void> {
  loading.value = true;
  error.value = null;
  
  try {
    // Note: We might need an admin endpoint to list all clients
    const response = await api.get<Client[]>('/admin/clients');
    clients.value = Array.isArray(response) ? response : (response as any).data || [];
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to fetch clients';
    error.value = message;
  } finally {
    loading.value = false;
  }
}

async function fetchClientById(id: string): Promise<void> {
  loading.value = true;
  error.value = null;
  
  try {
    const response = await api.get<Client>(`/client/${id}`);
    selectedClient.value = (response as any).data || response;
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to fetch client';
    error.value = message;
  } finally {
    loading.value = false;
  }
}

async function fetchClientAppointments(clientId: string): Promise<void> {
  loading.value = true;
  error.value = null;
  
  try {
    const response = await api.get<ClientAppointment[]>(`/client/${clientId}/appointments`);
    clientAppointments.value = Array.isArray(response) ? response : (response as any).data || [];
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to fetch appointments';
    error.value = message;
  } finally {
    loading.value = false;
  }
}

async function deleteClient(id: string): Promise<{ success: boolean; error?: string }> {
  try {
    await api.delete(`/client/${id}`);
    clients.value = clients.value.filter((client: Client) => client.id !== id);
    if (selectedClient.value?.id === id) {
      selectedClient.value = null;
    }
    return { success: true };
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to delete client';
    return { success: false, error: message };
  }
}

// Export store
export const clientStore = {
  // State
  clients,
  selectedClient,
  clientAppointments,
  loading,
  error,
  
  // Actions
  fetchClients,
  fetchClientById,
  fetchClientAppointments,
  deleteClient,
};
