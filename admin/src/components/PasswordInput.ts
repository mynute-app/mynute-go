import { html } from 'htm/preact';
import { useState } from 'preact/hooks';

interface PasswordInputProps {
    value: string;
    onInput: (e: Event) => void;
    placeholder?: string;
    label?: string;
    required?: boolean;
    minLength?: number;
    helpText?: string;
    name?: string;
    testId?: string;
}

export default function PasswordInput({
    value,
    onInput,
    placeholder = '••••••••',
    label = 'Password',
    required = false,
    minLength,
    helpText,
    name,
    testId,
}: PasswordInputProps) {
    const [showPassword, setShowPassword] = useState(false);

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

    return html`
        <div>
            ${label && html`
                <label class="block text-sm font-medium text-gray-700 mb-2">
                    ${label} ${required && html`<span class="text-red-500">*</span>`}
                </label>
            `}
            <div class="relative">
                <input
                    type=${showPassword ? 'text' : 'password'}
                    value=${value}
                    onInput=${onInput}
                    name=${name}
                    data-testid=${testId}
                    class="w-full px-4 py-2 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    placeholder=${placeholder}
                    required=${required}
                    minLength=${minLength}
                />
                <button
                    type="button"
                    onClick=${togglePasswordVisibility}
                    data-testid="${testId ? testId + '-toggle' : 'password-toggle'}"
                    class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                    aria-label=${showPassword ? 'Hide password' : 'Show password'}
                >
                    ${showPassword
                        ? html`
                            <!-- Eye slash icon (hide) -->
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                                <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.065 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                            </svg>
                        `
                        : html`
                            <!-- Eye icon (show) -->
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                                <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                            </svg>
                        `
                    }
                </button>
            </div>
            ${helpText && html`
                <p class="text-xs text-gray-500 mt-1">${helpText}</p>
            `}
        </div>
    `;
}
