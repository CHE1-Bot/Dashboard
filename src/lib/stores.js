import { writable, derived, get } from 'svelte/store';
import { api } from './api.js';

export const user = writable(null);            // current Discord user, or null
export const guilds = writable([]);            // all guilds the user is in
export const currentGuildId = writable(
  localStorage.getItem('che1.currentGuildId') || null
);

currentGuildId.subscribe((v) => {
  if (v) localStorage.setItem('che1.currentGuildId', v);
  else localStorage.removeItem('che1.currentGuildId');
});

export const currentGuild = derived(
  [guilds, currentGuildId],
  ([$gs, $id]) => $gs.find((g) => g.id === $id) || null
);

export async function loadMe() {
  try {
    const me = await api.get('/auth/me');
    user.set(me);
    return me;
  } catch (e) {
    user.set(null);
    throw e;
  }
}

export async function loadGuilds() {
  const list = await api.get('/guilds');
  guilds.set(list || []);
  // Auto-select a manageable guild if none chosen
  const cur = get(currentGuildId);
  const valid = (list || []).some((g) => g.id === cur && g.bot_present);
  if (!valid) {
    const first = (list || []).find((g) => g.bot_present);
    currentGuildId.set(first ? first.id : null);
  }
  return list;
}

export async function logout() {
  try { await api.post('/auth/logout'); } catch {}
  user.set(null);
  guilds.set([]);
  currentGuildId.set(null);
}

export function loginUrl() { return '/api/auth/login'; }
