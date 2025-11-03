// Type declarations for CDN imports
// These modules are loaded via importmap in index.html

declare module 'preact' {
  export interface VNode<P = {}> {
    type: any;
    props: P;
    key: any;
  }
  
  export type ComponentChildren = VNode<any> | string | number | null | undefined | boolean | ComponentChildren[];
  export type FunctionComponent<P = {}> = (props: P) => VNode<any> | null;
  
  export function render(vnode: VNode<any>, container: Element | Document | ShadowRoot | DocumentFragment): void;
  export function h(type: any, props: any, ...children: any[]): VNode<any>;
}

declare module 'preact/hooks' {
  export function useState<T>(initialState: T | (() => T)): [T, (value: T | ((prevState: T) => T)) => void];
  export function useEffect(effect: () => (void | (() => void)), deps?: any[]): void;
  export function useRef<T>(initialValue: T): { current: T };
  export function useMemo<T>(factory: () => T, deps: any[]): T;
  export function useCallback<T extends (...args: any[]) => any>(callback: T, deps: any[]): T;
}

declare module 'htm/preact' {
  import { VNode } from 'preact';
  export function html(strings: TemplateStringsArray, ...values: any[]): VNode<any>;
}

declare module '@preact/signals' {
  export interface Signal<T = any> {
    value: T;
  }
  
  export interface ReadonlySignal<T = any> {
    readonly value: T;
  }
  
  export function signal<T>(value: T): Signal<T>;
  export function computed<T>(compute: () => T): ReadonlySignal<T>;
  export function effect(fn: () => void): () => void;
  export function batch(fn: () => void): void;
}

declare module 'preact-router' {
  import { VNode, ComponentChildren } from 'preact';
  
  export interface RouterProps {
    children?: ComponentChildren;
    onChange?: (event: { url: string; previous?: string }) => void;
  }
  
  export interface RouteProps<Props = any> {
    path: string;
    component?: any;
    children?: ComponentChildren;
  }
  
  export function Router(props: RouterProps): VNode<any>;
  export function Route<T = any>(props: RouteProps<T> & T): VNode<any> | null;
  export function route(url: string, replace?: boolean): boolean;
  export function getCurrentUrl(): string;
}
