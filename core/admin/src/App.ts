import { html } from 'htm/preact';
import { Router } from 'preact-router';
import { useEffect } from 'preact/hooks';

// Store
import { authStore } from './stores/authStore.ts';

// Pages
import Login from './pages/Login.ts';
import Dashboard from './pages/Dashboard.ts';
import Users from './pages/Users.ts';
import Companies from './pages/Companies.ts';
import CompanyDetail from './pages/CompanyDetail.ts';
import Clients from './pages/Clients.ts';

// Components
import Layout from './components/Layout.ts';

export default function App() {
    useEffect(() => {
        // Check if user is already authenticated
        authStore.checkAuth();
    }, []);

    if (authStore.loading.value) {
        return html`
            <div class="flex items-center justify-center min-h-screen">
                <div class="text-xl">Loading...</div>
            </div>
        `;
    }

    if (!authStore.isAuthenticated.value) {
        return html`<${Login} />`;
    }

    return html`
        <${Layout}>
            <${Router} basepath="/admin">
                <${Dashboard} path="/" />
                <${Users} path="/users" />
                <${Companies} path="/companies" />
                <${CompanyDetail} path="/companies/:id" />
                <${Clients} path="/clients" />
            </${Router}>
        </${Layout}>
    `;
}
