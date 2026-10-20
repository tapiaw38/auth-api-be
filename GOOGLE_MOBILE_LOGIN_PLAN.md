# Google login dentro de la app Android (Capacitor) — plan técnico

## Por qué existe este documento

El botón de Google se queda "trabado" (carga infinita) dentro de la app Android. Causa raíz verificada: `GoogleButton.vue` usa Google Identity Services como script web (`initCodeClient`, `ux_mode: 'popup'`), y Google **bloquea ese flujo dentro de cualquier WebView embebido** (`disallowed_useragent`) — es política de Google, no un bug nuestro. El callback nunca se dispara, así que `isLoading` queda en `true` para siempre.

Dos soluciones se evaluaron y se descartaron antes de llegar a esta:

1. **Esquema de URI custom** (`com.practiq.app:/oauth2redirect`) — **Google ya no lo permite en Android**, confirmado contra la documentación oficial (`developers.google.com/identity/protocols/oauth2/native-app`): "Custom URI schemes are no longer supported on Android and Chrome apps."
2. **Android App Links** (`https://app.practiq.com.ar/auth/mobile-callback` interceptado a nivel OS vía `assetlinks.json` verificado) — técnicamente viable, pero exige `AndroidManifest.xml` con `autoVerify`, hosting de `/.well-known/assetlinks.json` con el SHA-256 del certificado de firma (debug y luego release, dos fingerprints a mantener), y la verificación de Android puede tardar o fallar silenciosamente en el primer instante tras instalar. Innecesariamente frágil para lo que hace falta acá.

**La solución elegida: callback web + polling.** Sin App Links, sin assetlinks.json, sin cambios al `AndroidManifest.xml` más allá de lo que Capacitor ya genera.

## Cómo funciona

```
┌─────────────┐  1. abre browser del sistema (Chrome Custom Tabs)  ┌──────────────────┐
│  App Android │ ─────────────────────────────────────────────────▶│  accounts.google  │
│  (WebView)   │                                                    │  .com/.../auth    │
└─────────────┘                                                    └──────────────────┘
      │                                                                      │
      │ 2. empieza a hacer polling a                                        │ usuario inicia sesión,
      │    auth-api-be cada ~1.5s                                           │ Google redirige a
      │    preguntando por `state`                                         │ redirect_uri con ?code&state
      ▼                                                                      ▼
┌──────────────────┐                                            ┌───────────────────────────────┐
│   auth-api-be     │◀───────────────────────────────────────── │  app.practiq.com.ar/auth/      │
│  (sesión en       │   3. esa página web (normal, en un         │  mobile-callback                │
│   memoria,         │      tab de Chrome, no en la app)          │  (practiq-fe, nueva ruta)       │
│   keyed by state) │      hace POST {code, state}               │  lee ?code&state del querystring│
└──────────────────┘                                            └───────────────────────────────┘
      │
      │ 4. cuando el polling de la app ve status:"done",
      │    ya tiene token+refresh_token+user — cierra el
      │    browser y entra logueada. Nunca pasó por deep
      │    link ni por intercepción de OS.
      ▼
   App logueada
```

Los tokens de sesión **nunca tocan el contexto web del Custom Tab** — el tab de Chrome solo dispara el POST y muestra "ya podés volver a la app"; quien recibe los tokens reales es la app, vía su propio polling a `auth-api-be`. Esto es más seguro que meter tokens en una URL de redirect que un navegador externo podría loguear.

## Repos y ramas (todas sobre `develop`, sin tocarla)

| Repo | Rama | Worktree |
|---|---|---|
| `auth-api-be` | `feat/google-mobile-oauth` | `../auth-api-be-mobile-oauth` |
| `practiq-fe` | `feat/google-mobile-callback` | `../practiq-fe-mobile-callback` |
| `practiq-fe` (Capacitor) | `feat/capacitor-android` (ya existe) | `../practiq-fe-capacitor` |

La página de callback va en una rama **separada** de `feat/capacitor-android` porque es código web normal que se despliega en `app.practiq.com.ar` vía el deploy habitual de `develop` — no depende de que la app Android exista para tener sentido, y no debería esperar a que el trabajo de Capacitor (más grande, más lento) termine para salir.

---

## Parte 1 — `auth-api-be` (worktree: `auth-api-be-mobile-oauth`)

### Contexto que ya verifiqué en el código (no reinventar)

`internal/usecases/user/login.go` — `LoginUsecase.Execute` con `LoginInput{SsoType:"google", Code:code}` ya hace: exchange del código con Google, `GetUserInfo`, busca o crea el usuario (`googleLogin()`, líneas ~98-166), emite sesión (`issueSession`). **Esta lógica no se duplica** — se extrae a una función interna reutilizable.

`internal/adapters/web/integrations/sso/integration.go` — `ExchangeCode(ctx, code)` arma un `oauth2.Config` con `RedirectURL: "postmessage"` **hardcodeado**. Ese valor es un sentinel especial de Google que solo funciona con el flujo popup de GIS en web — no sirve para un redirect real. Hay que parametrizarlo.

### Cambios

**1. `internal/adapters/web/integrations/sso/integration.go`**

Cambiar la firma del método en la interfaz y la implementación:

```go
// antes
ExchangeCode(context.Context, string) (*oauth2.Token, error)

// después
ExchangeCode(ctx context.Context, code string, redirectURI string) (*oauth2.Token, error)
```

`initConfig` deja de setear `RedirectURL` a nivel struct — se construye el `oauth2.Config` completo (o se clona el existente y se pisa `RedirectURL`) dentro de `ExchangeCode`, usando el parámetro `redirectURI`.

**2. `internal/usecases/user/login.go`**

- El único call site existente (`googleLogin`, línea ~99) pasa a llamar `app.Integrations.SSO.ExchangeCode(ctx, input.Code, "postmessage")` — **comportamiento idéntico al actual, cero regresión en el login web**.
- Extraer el bloque de "buscar o crear usuario a partir de `userInfo`" (desde `user, appErr := app.Repositories.User.Get(...)` hasta el final de `googleLogin`, aproximadamente las líneas 108-166) a una función privada nueva, p. ej. `resolveGoogleUser(ctx context.Context, app *appcontext.Context, userInfo *sso.SocialUser) (*string, apperrors.ApplicationError)`. `googleLogin` queda como: exchange + GetUserInfo + llamar a `resolveGoogleUser`.

**3. Nuevo archivo `internal/usecases/user/google_mobile_session.go`** (o nombre similar)

Contrato: una sesión en memoria, de un solo uso, con expiración corta. No hace falta Redis ni persistencia — el volumen es bajísimo (logins interactivos) y la ventana de vida es de minutos.

```go
type GoogleMobileSessionStatus string

const (
    MobileSessionPending GoogleMobileSessionStatus = "pending"
    MobileSessionDone    GoogleMobileSessionStatus = "done"
    MobileSessionError   GoogleMobileSessionStatus = "error"
)

type GoogleMobileSessionResult struct {
    Status       GoogleMobileSessionStatus
    Token        string
    RefreshToken string
    Data         UserOutputData
    Message      string // solo si Status == error, mensaje seguro para mostrar
}
```

Store: un `map[string]storedSession` protegido por `sync.Mutex`, cada entrada con `createdAt time.Time`. TTL de 5 minutos, chequeado de forma perezosa (al leer/escribir, no hace falta goroutine de limpieza — si se quiere prolijo, un `time.AfterFunc` que borra la entrada a los 5 min también sirve, cualquiera de las dos es aceptable).

Dos funciones:

- `CompleteGoogleMobileLogin(ctx context.Context, code string, state string) apperrors.ApplicationError` — constante `const mobileRedirectURI = "https://app.practiq.com.ar/auth/mobile-callback"`. Llama `app.Integrations.SSO.ExchangeCode(ctx, code, mobileRedirectURI)`, `GetUserInfo`, `resolveGoogleUser` (la función extraída en el paso 2), busca el `user` completo igual que hace `LoginUsecase.Execute` después de `googleLogin`, llama `issueSession`, y guarda `{Status: Done, Token, RefreshToken, Data}` en el store bajo la key `state`. Si cualquier paso falla, guarda `{Status: Error, Message: "..."}` (mensaje genérico tipo "No se pudo completar el inicio de sesión con Google", no el error interno crudo) bajo la misma key, y devuelve el `apperrors.ApplicationError` igual (para que el handler HTTP pueda loguearlo server-side con detalle).
- `PollGoogleMobileSession(state string) (*GoogleMobileSessionResult, bool)` — el segundo valor es `found bool`. Si existe y no expiró: devuelve el resultado **y borra la entrada** (de un solo uso — evita que alguien reuse el `state` para repetir el login). Si no existe o expiró: `found = false`.

**4. Handlers HTTP** (seguir el patrón de los handlers existentes en `internal/adapters/web/handlers/`, ver cómo está armado el de `/auth/login` para copiar estilo)

- `POST /auth/google/mobile/callback` — body `{code, state}` (JSON). Llama `CompleteGoogleMobileLogin`. **Responde 200 con un mensaje genérico siempre** (algo como `{"message":"Ya podés volver a la app"}`), incluso si `CompleteGoogleMobileLogin` devolvió error — el detalle del error ya quedó guardado en el store para que la音 app lo lea vía polling; esta página la ve un browser público, no conviene exponer detalles ahí. Loguear el error server-side (`log.Printf` o lo que use el resto del código) antes de responder genérico.
- `GET /auth/google/mobile/poll/:state` — llama `PollGoogleMobileSession(state)`. Si `found == false`: devolver `{"status":"pending"}` con 200 (la app sigue pollenado). Si `found == true`: devolver el `GoogleMobileSessionResult` tal cual como JSON (con `token`, `refresh_token`, `data` si `status:"done"`, o `message` si `status:"error"`).

Registrar ambas rutas en el router (buscar dónde se registra `/auth/login` y poner las nuevas rutas al lado, mismo grupo/middleware que no requiera auth — son endpoints públicos por diseño, igual que `/auth/login`).

### Checklist de salida — Parte 1

- [ ] `go build ./...` pasa
- [ ] `go vet ./...` limpio
- [ ] `go test ./...` — toda la suite sigue en verde (no se tocó ningún contrato público existente aparte de la firma interna de `ExchangeCode`, que es detalle de implementación)
- [ ] El login web con Google (`LoginInput{SsoType:"google", Code}` vía `/auth/login`) sigue funcionando exactamente igual — verificar leyendo el código, no hace falta un entorno para probarlo en vivo, pero sí confirmar que el call site pasa `"postmessage"` literal.

---

## Parte 2 — `practiq-fe` (worktree: `practiq-fe-mobile-callback`)

Una sola ruta nueva, aislada, sin tocar `LoginView.vue` ni nada del flujo web existente.

**1. Nueva vista**, p. ej. `src/views/GoogleMobileCallbackView.vue`:

```vue
<script setup lang="ts">
import { onMounted, ref } from 'vue';
import axios from 'axios';

const AUTH_BASE_URL = import.meta.env.VITE_AUTH_API_URL || 'http://localhost:8082';
const done = ref(false);

onMounted(async () => {
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const state = params.get('state');
    if (code && state) {
        try {
            await axios.post(`${AUTH_BASE_URL}/auth/google/mobile/callback`, { code, state });
        } catch {
            // La app sigue pollenado y va a ver el error ahí — esta página no
            // tiene nada más que hacer si el POST falla.
        }
    }
    done.value = true;
});
</script>

<template>
    <div class="mobile-callback">
        <p>{{ done ? 'Ya podés volver a la app.' : 'Completando el inicio de sesión…' }}</p>
    </div>
</template>

<style scoped>
.mobile-callback {
    display: grid;
    place-items: center;
    min-height: 100vh;
    font-family: var(--font-body-family, sans-serif);
    text-align: center;
    padding: 24px;
}
</style>
```

No usar ningún layout existente (`StudentLayout`, `TeacherLayout`) — esta página la ve un usuario sin sesión, en un tab de Chrome, por unos segundos. Standalone.

**2. Ruta nueva** en el router (buscar `src/router/index.ts` o equivalente, agregar junto a las rutas públicas como `/login`):

```ts
{
    path: '/auth/mobile-callback',
    name: 'google-mobile-callback',
    component: () => import('@/views/GoogleMobileCallbackView.vue'),
    meta: { requiresAuth: false },
}
```

Verificar el nombre exacto de la propiedad meta que usa el resto de las rutas públicas (puede no ser `requiresAuth: false` literal, revisar cómo está hecho `/login` y calcar).

### Checklist de salida — Parte 2

- [ ] `npm run build` completo pasa (`check:css-tokens && check:fractions && vue-tsc -b && vite build`)
- [ ] La ruta `/auth/mobile-callback?code=test&state=test` carga sin crashear localmente (`npm run dev`, aunque el POST falle contra un backend no corriendo — la página debe mostrar igual el mensaje final, no quedarse colgada ni tirar un error no capturado)

---

## Parte 3 — `practiq-fe-capacitor` (worktree existente, rama `feat/capacitor-android`)

**1. Instalar dependencia:**

```bash
npm install @capacitor/browser
```

(`@capacitor/app` NO hace falta para este enfoque — no se usa deep link ni `appUrlOpen`, todo es polling HTTP.)

**2. Nuevo composable** `src/platform/useNativeGoogleLogin.ts`:

```ts
import { Browser } from '@capacitor/browser';
import type { LoginResponse } from '@/types/auth';

const AUTH_BASE_URL = import.meta.env.VITE_AUTH_API_URL || 'http://localhost:8082';
const POLL_INTERVAL_MS = 1500;
const POLL_TIMEOUT_MS = 180_000; // 3 minutos

export async function loginWithGoogleNative(): Promise<LoginResponse> {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
    const state = crypto.randomUUID();
    const redirectUri = 'https://app.practiq.com.ar/auth/mobile-callback';

    const authUrl = new URL('https://accounts.google.com/o/oauth2/v2/auth');
    authUrl.searchParams.set('client_id', clientId);
    authUrl.searchParams.set('redirect_uri', redirectUri);
    authUrl.searchParams.set('response_type', 'code');
    authUrl.searchParams.set('scope', 'openid email profile');
    authUrl.searchParams.set('state', state);
    authUrl.searchParams.set('prompt', 'select_account');

    await Browser.open({ url: authUrl.toString() });

    const startedAt = Date.now();
    while (Date.now() - startedAt < POLL_TIMEOUT_MS) {
        await new Promise((resolve) => setTimeout(resolve, POLL_INTERVAL_MS));
        const res = await fetch(`${AUTH_BASE_URL}/auth/google/mobile/poll/${state}`);
        const body = await res.json();

        if (body.status === 'done') {
            await Browser.close().catch(() => {});
            return { token: body.token, refresh_token: body.refresh_token, data: body.data };
        }
        if (body.status === 'error') {
            await Browser.close().catch(() => {});
            throw new Error(body.message || 'No se pudo iniciar sesión con Google.');
        }
        // status === 'pending' -> seguir pollenado
    }

    await Browser.close().catch(() => {});
    throw new Error('El inicio de sesión con Google tardó demasiado. Probá de nuevo.');
}
```

Ajustar el tipo `LoginResponse` al import real (`src/types/auth.ts` en el repo principal) — como este worktree es un checkout separado del mismo repo que `practiq-fe`, el archivo existe en el mismo path, usarlo tal cual.

**3. Modificar `src/components/auth/GoogleButton.vue`:**

Agregar detección de plataforma nativa y un emit nuevo, sin tocar el flujo web existente:

```ts
import { Capacitor } from '@capacitor/core';
import { loginWithGoogleNative } from '@/platform/useNativeGoogleLogin';

// agregar al defineEmits existente:
const emit = defineEmits<GoogleButtonEmits & { (e: 'session', payload: LoginResponse): void }>();
```

En `loginWithGoogle()`, bifurcar al principio:

```ts
async function loginWithGoogle() {
    if (Capacitor.isNativePlatform()) {
        isLoading.value = true;
        try {
            const session = await loginWithGoogleNative();
            emit('session', session);
        } catch (err: any) {
            // dejar que el componente padre muestre el error — emitir algo
            // que LoginView ya sepa interpretar, o agregar un emit 'error'
            // simétrico si no existe. Revisar GoogleButtonEmits actual antes
            // de decidir el nombre exacto.
        } finally {
            isLoading.value = false;
        }
        return;
    }
    // ... código existente del flujo web, sin cambios
}
```

Revisar `GoogleButtonEmits` (`GoogleButton.types.ts`) para ver si ya existe algún emit de error reutilizable antes de inventar uno nuevo.

**4. Modificar `src/views/LoginView.vue`:**

Agregar un handler nuevo para el emit `session`, que reutiliza `finalizeSession` **sin pasar por `authService.login`** (porque el polling ya completó el login server-side, ya hay tokens reales):

```ts
async function handleGoogleMobileSession(response: LoginResponse) {
    resetMessages();
    loading.value = true;
    try {
        await finalizeSession(response, profileType.value);
    } catch (err: any) {
        errorMsg.value = err.response?.data?.message || 'No se pudo iniciar con Google.';
    } finally {
        loading.value = false;
    }
}
```

Y en el template:

```html
<GoogleButton @code="handleGoogleLogin" @session="handleGoogleMobileSession" />
```

El `@code` sigue andando igual para web (nunca se dispara en nativo, porque `loginWithGoogle()` hace `return` antes de llegar a esa rama en plataforma nativa).

### Checklist de salida — Parte 3

- [ ] `npm run build` completo pasa
- [ ] `npx cap sync android` pasa
- [ ] Compilar APK (`cd android && JAVA_HOME=/usr/lib/jvm/java-21-openjdk ./gradlew assembleDebug` — ver `BUILD_ANDROID.md` para el setup completo de esta máquina) e instalar en el dispositivo conectado
- [ ] Probar el login con Google en el dispositivo físico de punta a punta, una vez que las Partes 1 y 2 estén desplegadas (avisar antes de desplegar, igual que con CORS — Parte 1 toca `auth-api-be` en producción)

---

## Orden de ejecución

1. Parte 1 (`auth-api-be`) completa, con su checklist en verde.
2. Parte 2 (`practiq-fe`) completa, con su checklist en verde. Independiente de la Parte 1, puede ir en paralelo.
3. Avisar antes de pushear/desplegar Partes 1 y 2 a producción — son las dos piezas que hacen falta desplegadas para que la Parte 3 se pueda probar de punta a punta.
4. Parte 3 (`practiq-fe-capacitor`), una vez que 1 y 2 estén en producción.
5. Build + instalación del APK + prueba real en el celular.

## Reglas de estilo (no negociables)

- **No agregar comentarios en el código.** Ni explicativos, ni de sección, ni los que aparecen de ejemplo en los snippets de este documento (esos snippets son ilustrativos del contrato que tiene que existir, no texto para copiar literal). Si el código necesita explicarse, el commit message es el lugar — nunca el código.
- **Seguir el patrón ya existente en cada repo**, no inventar uno nuevo: mismo estilo de manejo de errores (`apperrors.ApplicationError` en Go, cómo arma los mensajes el resto de `mappings/`), mismo estilo de definición de rutas/handlers que los endpoints vecinos, mismo estilo de composables/components que el resto de `src/`. Antes de escribir un archivo nuevo, mirar dos o tres archivos equivalentes ya existentes en ese mismo repo y calcar la forma, no solo la lógica.

## Qué NO hacer

- No tocar `GoogleButton.vue`'s flujo web existente (`ux_mode: 'popup'`, `initCodeClient`) — sigue funcionando para navegadores normales, no se toca.
- No implementar Android App Links / `assetlinks.json` / cambios al `AndroidManifest.xml` — se descartó esa vía, no revivirla "por las dudas".
- No usar PKCE — el exchange lo hace el backend con `ClientSecret` (cliente confidencial), no es necesario.
- No persistir las sesiones de polling en base de datos — en memoria alcanza, TTL corto, volumen bajo.
- No cambiar el contrato de `LoginUsecase.Execute` / `LoginInput` existente — toda la funcionalidad nueva es aditiva, en funciones y endpoints nuevos.
