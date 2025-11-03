import { html } from 'htm/preact';
import { useState, useEffect } from 'preact/hooks';
import { authStore } from '../stores/authStore.ts';
import { api } from '../utils/api.ts';
import PasswordInput from '../components/PasswordInput.ts';

export default function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [hasSuperAdmin, setHasSuperAdmin] = useState<boolean | null>(null);
    const [isRegistering, setIsRegistering] = useState(false);
    
    // Registration form fields
    const [name, setName] = useState('');
    const [surname, setSurname] = useState('');
    const [regEmail, setRegEmail] = useState('');
    const [regPassword, setRegPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [registrationSuccess, setRegistrationSuccess] = useState(false);

    // Check if there are any admins on mount
    useEffect(() => {
        checkForAdmins();
    }, []);

    const checkForAdmins = async () => {
        try {
            const response = await api.get<{ has_superadmin: boolean }>('/admin/are_there_any_superadmin');
            setHasSuperAdmin(response.has_superadmin);
        } catch (err) {
            console.error('Failed to check for admins:', err);
            setHasSuperAdmin(true); // Default to login form on error
        }
    };

    const handleSubmit = async (e: Event) => {
        e.preventDefault();
        setError('');

        const result = await authStore.login(email, password);
        
        if (!result.success) {
            setError(result.error || 'Login failed');
        }
    };

    const handleRegistration = async (e: Event) => {
        e.preventDefault();
        setError('');

        // Validation
        if (regPassword !== confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        if (regPassword.length < 8) {
            setError('Password must be at least 8 characters long');
            return;
        }

        setIsRegistering(true);

        try {
            // Create first admin
            await api.post('/admin/first_superadmin', {
                name,
                surname,
                email: regEmail,
                password: regPassword,
            });

            // Send verification email
            await api.post(`/admin/send-verification-code/email/${regEmail}`);

            setRegistrationSuccess(true);
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Registration failed';
            setError(message);
        } finally {
            setIsRegistering(false);
        }
    };

    // Show loading state while checking for admins
    if (hasSuperAdmin === null) {
        return html`
            <div class="min-h-screen flex items-center justify-center bg-gray-100">
                <div class="text-xl text-gray-600">Loading...</div>
            </div>
        `;
    }

    // Show registration success message
    if (registrationSuccess) {
        return html`
            <div class="min-h-screen flex items-center justify-center bg-gray-100">
                <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
                    <div class="text-center">
                        <div class="text-6xl mb-4">✅</div>
                        <h1 class="text-3xl font-bold text-gray-900 mb-4">Registration Successful!</h1>
                        <p class="text-gray-600 mb-6">
                            A verification email has been sent to <strong>${regEmail}</strong>.
                            Please check your inbox and verify your email address before logging in.
                        </p>
                        <button
                            onClick=${() => {
                                setRegistrationSuccess(false);
                                setHasSuperAdmin(true);
                            }}
                            class="bg-primary text-white py-2 px-6 rounded-lg hover:bg-blue-600 transition-colors"
                        >
                            Go to Login
                        </button>
                    </div>
                </div>
            </div>
        `;
    }

    // Show registration form if no admin exists
    if (!hasSuperAdmin) {
        return html`
            <div class="min-h-screen flex items-center justify-center bg-gray-100 py-12">
                <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
                    <div class="text-center mb-8">
                        <h1 class="text-3xl font-bold text-gray-900 mb-2">Welcome to Mynute Admin</h1>
                        <p class="text-gray-600">Create your first admin account to get started</p>
                    </div>
                    
                    <form onSubmit=${handleRegistration} class="space-y-4">
                        ${error && html`
                            <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                                ${error}
                            </div>
                        `}
                        
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                First Name <span class="text-red-500">*</span>
                            </label>
                            <input
                                type="text"
                                value=${name}
                                onInput=${(e: Event) => setName((e.target as HTMLInputElement).value)}
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                                placeholder="John"
                                required
                            />
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                Last Name <span class="text-red-500">*</span>
                            </label>
                            <input
                                type="text"
                                value=${surname}
                                onInput=${(e: Event) => setSurname((e.target as HTMLInputElement).value)}
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                                placeholder="Doe"
                                required
                            />
                        </div>
                        
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                Email <span class="text-red-500">*</span>
                            </label>
                            <input
                                type="email"
                                value=${regEmail}
                                onInput=${(e: Event) => setRegEmail((e.target as HTMLInputElement).value)}
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                                placeholder="admin@mynute.com"
                                required
                            />
                        </div>
                        
                        <${PasswordInput}
                            value=${regPassword}
                            onInput=${(e: Event) => setRegPassword((e.target as HTMLInputElement).value)}
                            label="Password"
                            required=${true}
                            minLength=${8}
                            helpText="Must be at least 8 characters"
                            testId="registration-password"
                        />

                        <${PasswordInput}
                            value=${confirmPassword}
                            onInput=${(e: Event) => setConfirmPassword((e.target as HTMLInputElement).value)}
                            label="Confirm Password"
                            required=${true}
                            minLength=${8}
                            testId="registration-confirm-password"
                        />
                        
                        <button
                            type="submit"
                            disabled=${isRegistering}
                            class="w-full bg-primary text-white py-2 px-4 rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                        >
                            ${isRegistering ? 'Creating Account...' : 'Create Admin Account'}
                        </button>
                    </form>
                </div>
            </div>
        `;
    }

    // Show login form
    return html`
        <div class="min-h-screen flex items-center justify-center bg-gray-100">
            <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
                <h1 class="text-3xl font-bold text-center mb-8">Mynute Admin</h1>
                
                <form onSubmit=${handleSubmit} class="space-y-6">
                    ${error && html`
                        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                            ${error}
                        </div>
                    `}
                    
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">
                            Email
                        </label>
                        <input
                            type="email"
                            value=${email}
                            onInput=${(e: Event) => setEmail((e.target as HTMLInputElement).value)}
                            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                            placeholder="admin@mynute.com"
                            required
                        />
                    </div>
                    
                    <${PasswordInput}
                        value=${password}
                        onInput=${(e: Event) => setPassword((e.target as HTMLInputElement).value)}
                        label="Password"
                        required=${true}
                        testId="login-password"
                    />
                    
                    <button
                        type="submit"
                        disabled=${authStore.loading.value}
                        class="w-full bg-primary text-white py-2 px-4 rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                        ${authStore.loading.value ? 'Logging in...' : 'Login'}
                    </button>
                </form>
            </div>
        </div>
    `;
}
