import Keycloak from 'keycloak-js'

export const keycloak = new Keycloak({
  url: 'https://localhost:8443',
  realm: 'shortener',
  clientId: 'shortener-web',
})

let ready: Promise<boolean> | null = null

export function ensureAuth(): Promise<boolean> {
  if (!ready) {
    ready = keycloak.init({
      onLoad: 'login-required',
      pkceMethod: 'S256',
      checkLoginIframe: false,
    })
  }
  return ready
}
