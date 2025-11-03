import { html } from 'htm/preact';
import { ComponentChildren } from 'preact';
import Sidebar from './Sidebar.ts';
import Header from './Header.ts';

interface LayoutProps {
    children: ComponentChildren;
}

export default function Layout({ children }: LayoutProps) {
    return html`
        <div class="flex h-screen bg-gray-100">
            <!-- Sidebar -->
            <${Sidebar} />
            
            <!-- Main Content -->
            <div class="flex-1 flex flex-col overflow-hidden">
                <!-- Header -->
                <${Header} />
                
                <!-- Page Content -->
                <main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100 p-6">
                    ${children}
                </main>
            </div>
        </div>
    `;
}
