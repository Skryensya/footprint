# Releasing fp

Este documento describe cómo crear nuevas releases de `fp`.

## Flujo de Release

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. Preparar release                                             │
│    - Actualizar CHANGELOG.md                                    │
│    - Verificar que tests pasen                                  │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Crear tag                                                    │
│    git tag -a v0.1.0 -m "Release v0.1.0"                       │
│    git push origin v0.1.0                                       │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. GitHub Actions se dispara automáticamente                    │
│    - Compila binarios para cada OS/arch                         │
│    - Crea GitHub Release                                        │
│    - Sube binarios como assets                                  │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Release publicada                                            │
│    - Usuarios pueden descargar binarios                         │
│    - `fp update` detecta nueva versión                          │
└─────────────────────────────────────────────────────────────────┘
```

## Paso a paso

### 1. Preparar la release

```bash
# Asegúrate de estar en main y actualizado
git checkout main
git pull origin main

# Verifica que los tests pasen
make test

# Actualiza CHANGELOG.md con los cambios de esta versión
# Sigue el formato Keep a Changelog (https://keepachangelog.com)
```

### 2. Crear el tag

```bash
# Formato de versión: vMAJOR.MINOR.PATCH (semantic versioning)
VERSION="v0.1.0"

# Crear tag anotado
git tag -a $VERSION -m "Release $VERSION"

# Push del tag (esto dispara el workflow)
git push origin $VERSION
```

### 3. Verificar el release

1. Ve a GitHub → Actions → verifica que el workflow "Release" esté corriendo
2. Una vez completado, ve a Releases → verifica que los binarios estén disponibles
3. Prueba la actualización:
   ```bash
   fp update
   ```

## Plataformas soportadas

El workflow genera binarios para:

| OS      | Arquitectura | Archivo                    |
|---------|--------------|----------------------------|
| macOS   | Apple Silicon| `fp_darwin_arm64.tar.gz`   |
| macOS   | Intel        | `fp_darwin_amd64.tar.gz`   |
| Linux   | x86_64       | `fp_linux_amd64.tar.gz`    |
| Linux   | ARM64        | `fp_linux_arm64.tar.gz`    |
| Windows | x86_64       | `fp_windows_amd64.zip`     |

## Versionado

Usamos [Semantic Versioning](https://semver.org/):

- **MAJOR**: Cambios incompatibles con versiones anteriores
- **MINOR**: Nueva funcionalidad compatible hacia atrás
- **PATCH**: Correcciones de bugs compatibles hacia atrás

Ejemplos:
- `v0.1.0` → Primera versión beta
- `v0.1.1` → Corrección de bug
- `v0.2.0` → Nueva funcionalidad
- `v1.0.0` → Primera versión estable

## Estructura del workflow

El archivo `.github/workflows/release.yml` hace lo siguiente:

1. **Trigger**: Se activa cuando se pushea un tag `v*`
2. **Build Matrix**: Compila para múltiples OS/arquitecturas en paralelo
3. **Artifacts**: Crea archivos `.tar.gz` (Unix) o `.zip` (Windows)
4. **Release**: Crea el GitHub Release y sube los binarios

## Hotfixes

Para releases urgentes:

```bash
# Crear branch desde el tag
git checkout -b hotfix/security-fix v0.1.0

# Hacer el fix
# ...

# Crear nuevo tag
git tag -a v0.1.1 -m "Security fix"
git push origin v0.1.1

# Merge a main
git checkout main
git merge hotfix/security-fix
git push origin main
```

## Troubleshooting

### El workflow falló

1. Ve a Actions → click en el workflow fallido
2. Revisa los logs para identificar el error
3. Si es un problema de compilación, arréglalo y crea un nuevo tag:
   ```bash
   git tag -d v0.1.0              # Elimina tag local
   git push origin :v0.1.0        # Elimina tag remoto
   # Arregla el problema
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

### Los binarios no se generaron

Verifica que el workflow tenga permisos de escritura:
- Settings → Actions → General → Workflow permissions → "Read and write permissions"

### `fp update` no detecta la nueva versión

El cache de actualización dura 24 horas. Para forzar:
```bash
fp config unset update_last_check
fp update
```
