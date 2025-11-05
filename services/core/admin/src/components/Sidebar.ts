import { html } from 'htm/preact';
import { route } from 'preact-router';

interface NavItem {
    path: string;
    label: string;
    icon: string;
}

export default function Sidebar() {
    const currentPath = window.location.pathname;
    
    const navItems: NavItem[] = [
        { path: '/', label: 'Dashboard', icon: '📊' },
        { path: '/companies', label: 'Companies', icon: '🏢' },
        { path: '/clients', label: 'Clients', icon: '👥' },
        { path: '/users', label: 'Admin Users', icon: '🔐' },
    ];

    const handleNavigate = (path: string) => {
        route(path);
    };

    const isActive = (path: string) => {
        // Remove /admin prefix from current path for comparison
        const relativePath = currentPath.replace(/^\/admin/, '') || '/';
        if (path === '/') {
            return relativePath === '/' || relativePath === '';
        }
        return relativePath.startsWith(path);
    };

    return html`
        <div class="bg-gray-900 text-white w-64 space-y-6 py-7 px-2 flex flex-col">
            <!-- Logo -->
            <div class="px-4 mb-4">
                <h1 class="text-2xl font-bold">Mynute</h1>
                <p class="text-sm text-gray-400">Admin Panel</p>
            </div>
            
            <!-- Navigation -->
            <nav class="flex-1">
                ${navItems.map((item: NavItem) => {
                    const active = isActive(item.path);
                    return html`
                        <a
                            key=${item.path}
                            onClick=${(e: Event) => {
                                e.preventDefault();
                                handleNavigate(item.path);
                            }}
                            href=${item.path}
                            class="${active 
                                ? 'bg-gray-800 text-white' 
                                : 'text-gray-400 hover:bg-gray-800 hover:text-white'
                            } flex items-center px-4 py-3 rounded-lg mb-1 transition-colors cursor-pointer"
                        >
                            <span class="text-xl mr-3">${item.icon}</span>
                            <span>${item.label}</span>
                        </a>
                    `;
                })}
            </nav>
        </div>
    `;
}
