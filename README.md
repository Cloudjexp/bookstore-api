# Library API

API para la gestión de una biblioteca.

## Estructura

- `cmd/api`: punto de entrada de la aplicación.
- `internal/domain`: entidades del dominio.
- `internal/handler`: controladores HTTP.
- `internal/service`: casos de uso y lógica de aplicación.
- `internal/repository`: contratos de persistencia.
- `internal/repository/postgres`: implementación PostgreSQL.
- `internal/middleware`: middlewares HTTP.
- `internal/router`: configuración de rutas.
- `migrations`: migraciones de base de datos.
- `pkg/response`: utilidades compartidas para respuestas.

## Estado

Scaffold inicial. La lógica de la API se implementará posteriormente.
