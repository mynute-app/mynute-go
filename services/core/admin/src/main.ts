import { render } from 'preact';
import { html } from 'htm/preact';
import App from './App.ts';

const appElement = document.getElementById('app');
if (appElement) {
    render(html`<${App} />`, appElement);
}
