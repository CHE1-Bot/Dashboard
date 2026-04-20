// Single fetch wrapper for all dashboard API calls.
// Sends cookies, parses JSON, throws on non-2xx.

async function request(method, path, body) {
  const opts = {
    method,
    credentials: 'include',
    headers: body ? { 'Content-Type': 'application/json' } : {},
  };
  if (body !== undefined) opts.body = JSON.stringify(body);
  const res = await fetch('/api' + path, opts);
  if (res.status === 401) {
    // bubble up so callers can redirect to login
    const err = new Error('unauthorized');
    err.status = 401;
    throw err;
  }
  const text = await res.text();
  const data = text ? JSON.parse(text) : null;
  if (!res.ok) {
    const err = new Error((data && data.error) || res.statusText);
    err.status = res.status;
    throw err;
  }
  return data;
}

export const api = {
  get:    (p)      => request('GET', p),
  post:   (p, b)   => request('POST', p, b ?? {}),
  patch:  (p, b)   => request('PATCH', p, b ?? {}),
  put:    (p, b)   => request('PUT', p, b ?? {}),
  del:    (p)      => request('DELETE', p),
};

// Convenience: guild-scoped helpers
export const guildApi = (gid) => ({
  get:   (sub)       => api.get(`/guilds/${gid}${sub}`),
  post:  (sub, b)    => api.post(`/guilds/${gid}${sub}`, b),
  patch: (sub, b)    => api.patch(`/guilds/${gid}${sub}`, b),
  del:   (sub)       => api.del(`/guilds/${gid}${sub}`),
});
