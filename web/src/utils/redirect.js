export function safeRedirect(value, origin = window.location.origin) {
  if (typeof value !== 'string' || !value.startsWith('/') || value.startsWith('//')) return '/'
  try {
    const target = new URL(value, origin)
    if (target.origin !== origin || /^\/(?:admin|login|register)(?:\/|$)/.test(target.pathname)) return '/'
    return target.pathname + target.search + target.hash
  } catch {
    return '/'
  }
}
