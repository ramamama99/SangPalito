/**
 * JWT Verification and Authentication Utilities for SangPalito Frontend
 */

export interface JwtPayload {
  sub: string
  facility_id: string
  role: string
  full_name?: string
  exp: number
  iat: number
  jti?: string
}

/**
 * Decodes and parses a JWT token safely without external dependencies.
 * Handles base64url padding and utf-8 decoding.
 */
export function parseJwt(token: string | null): JwtPayload | null {
  if (!token || typeof token !== 'string') return null
  const parts = token.trim().split('.')
  if (parts.length !== 3) return null

  try {
    const base64Url = parts[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), '=')
    const jsonPayload = decodeURIComponent(
      atob(padded)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    return JSON.parse(jsonPayload) as JwtPayload
  } catch (err) {
    console.warn('Invalid JWT format or parsing failure:', err)
    return null
  }
}

/**
 * Verifies if the JWT token is valid (well-formed, contains required claims, and has not expired).
 */
export function isJwtValid(token: string | null): boolean {
  const payload = parseJwt(token)
  if (!payload) return false

  // Verify mandatory claims for staff authentication
  if (!payload.sub || !payload.facility_id || !payload.role) {
    return false
  }

  // Check expiration timestamp (with 5 seconds clock drift buffer)
  if (payload.exp) {
    const now = Math.floor(Date.now() / 1000)
    if (now >= payload.exp - 5) {
      return false
    }
  }

  return true
}

/**
 * Verifies staff authentication state.
 * If the access token is missing or expired, it automatically attempts to use the refresh token.
 * Returns true if valid or successfully refreshed; false if not authenticated / logged out.
 */
export async function verifyStaffAuth(): Promise<boolean> {
  const token = localStorage.getItem('antrisehat-token')
  if (isJwtValid(token)) {
    return true
  }

  // If access token is expired/missing, check if refresh token exists
  const refreshToken = localStorage.getItem('antrisehat-refresh-token')
  if (!refreshToken) {
    clearStaffAuth()
    return false
  }

  // Try refreshing session with backend
  try {
    const res = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
    if (!res.ok) {
      clearStaffAuth()
      return false
    }
    const json = await res.json()
    if (json.success && json.data?.access_token && isJwtValid(json.data.access_token)) {
      localStorage.setItem('antrisehat-token', json.data.access_token)
      return true
    }
  } catch (err) {
    console.error('Failed to refresh token during auth verification:', err)
  }

  clearStaffAuth()
  return false
}

/**
 * Clears all staff session tokens and cached user data.
 */
export function clearStaffAuth() {
  localStorage.removeItem('antrisehat-token')
  localStorage.removeItem('antrisehat-refresh-token')
  localStorage.removeItem('antrisehat-user')
}

/**
 * Returns the current staff role decoded from the JWT token.
 */
export function getStaffRoleFromJwt(): string | null {
  const token = localStorage.getItem('antrisehat-token')
  const payload = parseJwt(token)
  return payload?.role ? payload.role.toLowerCase() : null
}

/**
 * Checks whether the current staff role from JWT is authorized for Dashboard Poli (medic, head, admin).
 */
export function canAccessPoliDashboard(): boolean {
  const role = getStaffRoleFromJwt()
  return role ? ['medic', 'head', 'admin'].includes(role) : false
}

