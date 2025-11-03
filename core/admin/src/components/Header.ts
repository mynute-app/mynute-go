import { html } from 'htm/preact';
import { authStore } from '../stores/authStore.ts';

export default function Header() {
    const handleLogout = () => {
        if (confirm('Are you sure you want to logout?')) {
            authStore.logout();
        }
    };

    return html`
        <header class="bg-white shadow-sm">
            <div class="flex items-center justify-between px-6 py-4">
                <div>
                    <h2 class="text-xl font-semibold text-gray-800">
                        Welcome back${authStore.user.value?.name ? ', ' + authStore.user.value.name : ''}!
                    </h2>
                </div>
                
                <div class="flex items-center space-x-4">
                    <!-- User Info -->
                    <div class="text-right">
                        <p class="text-sm font-medium text-gray-900">
                            ${authStore.user.value?.name || 'Admin User'}
                        </p>
                        <p class="text-xs text-gray-500">
                            ${authStore.user.value?.email || ''}
                        </p>
                    </div>
                    
                    <!-- Logout Button -->
                    <button
                        onClick=${handleLogout}
                        class="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg transition-colors"
                    >
                        Logout
                    </button>
                </div>
            </div>
        </header>
    `;
}
