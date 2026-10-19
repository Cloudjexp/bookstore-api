# Arquitectura de Library API

## 1. Descripción

Library API será una API REST en Go para administrar una biblioteca. El proyecto comenzará como un **monolito modular**: una sola aplicación, con módulos y responsabilidades claramente separados.

## 2. Estructura del proyecto

```text
library-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── book.go
│   │   ├── author.go
│   │   ├── member.go
│   │   └── loan.go
│   ├── handler/
│   │   ├── book_handler.go
│   │   ├── author_handler.go
│   │   ├── member_handler.go
│   │   └── loan_handler.go
│   ├── service/
│   │   ├── book_service.go
│   │   ├── author_service.go
│   │   ├── member_service.go
│   │   └── loan_service.go
│   ├── repository/
│   │   ├── book_repository.go
│   │   ├── author_repository.go
│   │   ├── member_repository.go
│   │   ├── loan_repository.go
│   │   └── postgres/
│   │       └── connection.go
│   ├── middleware/
│   │   └── middleware.go
│   └── router/
│       └── router.go
├── migrations/
│   └── *.sql
├── pkg/
│   └── response/
│       └── response.go
├── .env.example
├── .gitignore
├── ARCHITECTURE.md
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## 3. Responsabilidad de cada directorio

### `cmd/api`

Punto de entrada de la aplicación. `main.go` cargará la configuración, creará las dependencias, configurará el router e iniciará el servidor. No debe contener lógica de negocio.

### `internal/config`

Carga y mantiene la configuración de la aplicación: entorno, puerto HTTP, URL de PostgreSQL y futuras opciones de autenticación.

### `internal/domain`

Contiene las entidades y tipos principales del negocio:

- `Book`: libro.
- `Author`: autor.
- `Member`: miembro de la biblioteca.
- `Loan`: préstamo.

El dominio no debe depender de HTTP, PostgreSQL ni de frameworks.

### `internal/handler`

Contiene los controladores HTTP. Lee requests, decodifica JSON, valida entradas básicas, llama a los servicios y construye respuestas HTTP.

Los handlers no deben ejecutar SQL ni contener reglas complejas de negocio.

### `internal/service`

Contiene los casos de uso y reglas de negocio. Por ejemplo: validar disponibilidad, impedir préstamos a miembros inactivos y coordinar una devolución.

Los servicios deben depender de interfaces de repository, no de implementaciones concretas.

### `internal/repository`

Contiene los contratos de persistencia. Define las operaciones necesarias para crear, buscar, actualizar y eliminar entidades.

### `internal/repository/postgres`

Contiene la conexión y las implementaciones de repositories para PostgreSQL, incluyendo consultas, mapeos y transacciones.

### `internal/middleware`

Contendrá funcionalidades transversales: logging, recuperación de errores, CORS, autenticación, autorización y request ID.

### `internal/router`

Configura las rutas HTTP y conecta cada ruta con su handler. La API utilizará versionado mediante `/api/v1`.

### `pkg/response`

Contendrá estructuras y helpers reutilizables para respuestas JSON y errores HTTP.

### `migrations`

Contiene los scripts SQL versionados de la base de datos.

Ejemplo:

```text
000001_create_books.up.sql
000001_create_books.down.sql
000002_create_authors.up.sql
000002_create_authors.down.sql
```

## 4. Flujo de una petición

```text
Cliente HTTP
    ↓
Middleware
    ↓
Router
    ↓
Handler
    ↓
Service
    ↓
Repository
    ↓
PostgreSQL
```

Ejemplo para crear un libro:

1. El cliente envía `POST /api/v1/books`.
2. El router dirige la petición a `BookHandler`.
3. El handler decodifica y valida el JSON.
4. El handler llama a `BookService`.
5. El service aplica las reglas de negocio.
6. El service llama a `BookRepository`.
7. El repository guarda el libro en PostgreSQL.
8. El resultado vuelve al handler.
9. El handler responde con JSON y el código HTTP correspondiente.

## 5. Módulos del dominio

### Libros

- Crear, consultar, actualizar y eliminar libros.
- Buscar por título, ISBN o autor.
- Controlar ejemplares totales y disponibles.

### Autores

- Crear y consultar autores.
- Asociar autores con libros.
- Consultar los libros de un autor.

### Miembros

- Registrar miembros.
- Actualizar sus datos.
- Activar o desactivar miembros.
- Consultar su historial de préstamos.

### Préstamos

- Crear préstamos.
- Validar disponibilidad.
- Registrar fecha límite.
- Registrar devoluciones.
- Detectar préstamos vencidos.

## 6. Endpoints previstos

### Salud

```text
GET /health
```

### Libros

```text
GET    /api/v1/books
POST   /api/v1/books
GET    /api/v1/books/:id
PATCH  /api/v1/books/:id
DELETE /api/v1/books/:id
```

### Autores

```text
GET  /api/v1/authors
POST /api/v1/authors
GET  /api/v1/authors/:id
```

### Miembros

```text
GET   /api/v1/members
POST  /api/v1/members
GET   /api/v1/members/:id
PATCH /api/v1/members/:id
```

### Préstamos

```text
POST /api/v1/loans
GET  /api/v1/loans/:id
POST /api/v1/loans/:id/return
GET  /api/v1/members/:id/loans
```

## 7. Modelo de datos previsto

### `books`

```text
id
title
isbn
description
total_copies
available_copies
created_at
updated_at
```

### `authors`

```text
id
name
created_at
updated_at
```

### `book_authors`

```text
book_id
author_id
```

### `members`

```text
id
name
email
active
created_at
updated_at
```

### `loans`

```text
id
book_id
member_id
borrowed_at
due_at
returned_at
status
```

## 8. Reglas de negocio principales

1. Un libro no puede prestarse si `available_copies` es cero.
2. Un miembro inactivo no puede solicitar préstamos.
3. Cada ISBN debe ser único.
4. Cada email de miembro debe ser único.
5. Una devolución no debe ejecutarse dos veces.
6. Crear un préstamo y reducir la disponibilidad debe ser transaccional.
7. Devolver un libro debe incrementar su disponibilidad una sola vez.
8. Los préstamos deben conservarse para mantener el historial.
9. Un libro con historial de préstamos no debería eliminarse físicamente.

## 9. Configuración

Las variables se definirán en `.env` durante el desarrollo. El archivo `.env` no debe subirse al repositorio.

```env
APP_ENV=development
HTTP_PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/library?sslmode=disable
```

`.env.example` debe mantenerse actualizado con todas las variables necesarias.

## 10. Dependencias previstas

- HTTP: comenzar con `net/http`; usar `chi` si se necesitan rutas más avanzadas.
- PostgreSQL: `pgx`.
- Migraciones: `golang-migrate`.
- Validación: `go-playground/validator`.
- Logs: `log/slog`.
- Tests: paquete estándar `testing`.

Se recomienda agregar dependencias únicamente cuando sean necesarias.

## 11. Estrategia de pruebas

### Tests unitarios

Para servicios, validaciones y reglas de negocio.

### Tests de handlers

Para comprobar códigos HTTP, respuestas JSON, validación de requests y errores.

### Tests de integración

Para comprobar repositories PostgreSQL, migraciones, transacciones y el flujo de préstamos.

## 12. Orden de implementación

### Fase 1: servidor base

1. Iniciar el servidor HTTP.
2. Crear `GET /health`.
3. Crear respuestas JSON comunes.

### Fase 2: configuración

1. Leer variables de entorno.
2. Configurar el puerto.
3. Preparar `DATABASE_URL`.

### Fase 3: libros

1. Completar la entidad `Book`.
2. Crear la interfaz del repository.
3. Crear `BookService`.
4. Crear `BookHandler`.
5. Implementar el CRUD inicialmente en memoria.

### Fase 4: persistencia

1. Agregar PostgreSQL.
2. Crear migraciones.
3. Implementar `BookRepository` para PostgreSQL.
4. Agregar transacciones y tests de integración.

### Fase 5: autores y miembros

Implementar ambos módulos siguiendo el mismo flujo utilizado para libros.

### Fase 6: préstamos

Implementar préstamos y devoluciones con transacciones y validaciones de negocio.

### Fase 7: seguridad y operación

- Autenticación.
- Roles y permisos.
- Logs estructurados.
- CORS.
- Documentación OpenAPI.
- Docker y despliegue.

## 13. Convenciones de desarrollo

- Ejecutar `gofmt` antes de confirmar cambios.
- Mantener funciones pequeñas y con una sola responsabilidad.
- No colocar SQL en handlers o services.
- No colocar reglas de negocio en repositories.
- Usar interfaces donde exista una dependencia que deba probarse o reemplazarse.
- Devolver errores explícitos y manejarlos en el límite HTTP.
- Mantener los endpoints bajo `/api/v1`.
- Escribir tests junto con cada funcionalidad nueva.

## 14. Estado actual

El proyecto contiene el scaffold inicial de directorios y paquetes. Los tipos, servicios, handlers y repositories son estructuras base pendientes de implementación.

El primer objetivo funcional es:

```text
GET /health → HTTP 200
```
