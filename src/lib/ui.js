// Tiny formatters/helpers shared by components.

export function formatDate(iso) {
  if (!iso) return '—';
  const d = new Date(iso);
  if (isNaN(d)) return iso;
  return d.toLocaleString();
}

export function relativeTime(iso) {
  if (!iso) return '—';
  const d = new Date(iso);
  const diff = (d.getTime() - Date.now()) / 1000;
  const abs = Math.abs(diff);
  const units = [
    [60, 'second'],
    [60, 'minute'],
    [24, 'hour'],
    [7,  'day'],
    [4.3,'week'],
    [12, 'month'],
  ];
  let value = abs;
  let unit = 'second';
  for (const [factor, name] of units) {
    if (value < factor) { unit = name; break; }
    value /= factor;
    unit = name;
  }
  const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' });
  return rtf.format(Math.round(diff < 0 ? -value : value), unit);
}

export function humanBytes(n) {
  if (!n) return '0 B';
  const u = ['B', 'KB', 'MB', 'GB'];
  let i = 0;
  while (n >= 1024 && i < u.length - 1) { n /= 1024; i++; }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${u[i]}`;
}

export function guildIconUrl(g) {
  if (g && g.icon) return `https://cdn.discordapp.com/icons/${g.id}/${g.icon}.png`;
  return null;
}

export function avatarUrl(u) {
  if (u && u.avatar) return `https://cdn.discordapp.com/avatars/${u.id}/${u.avatar}.png`;
  return 'https://cdn.discordapp.com/embed/avatars/0.png';
}
