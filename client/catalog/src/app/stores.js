import { writable } from 'svelte/store';

export const navStore = writable({});
export const terminalStore = writable({});
export const dirtyStore = writable({'clean' : true});
