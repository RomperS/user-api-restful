# 🚀 API RESTful para la Gestión de Usuarios

## Descripción del Proyecto

Esta es una **API RESTful** dedicada exclusivamente a la gestión del recurso **usuario**.

### Stack Tecnológico y Arquitectura

| Componente | Detalle |
| :--- | :--- |
| **Lenguaje** | **Go** |
| **Arquitectura** | **Hexagonal** (Ports and Adapters) |
| **Persistencia** | **PostgreSQL** (Base de datos relacional) |
| **ORM** | **GORM** (Go's Object-Relational Mapping) |

La arquitectura hexagonal se ha elegido para mantener una clara separación de las preocupaciones (*separation of concerns*), aislando la lógica de negocio de los detalles de infraestructura.

## Endpoints Principales (Recurso `/users`)

La API soporta las operaciones fundamentales de un CRUD (*Create, Read, Update, Delete*) para la entidad `User`.

| Método | Ruta | Resumen | Descripción | Seguridad |
| :---: | :--- | :--- | :--- | :---: |
| **GET** | `/users` | Get All Users | Recupera una lista completa de todos los usuarios registrados. | Basic Auth |
| **POST** | `/users` | Create New User | Crea un nuevo usuario. | Basic Auth |
| **GET** | `/users/{id}` | Get User by ID | Recupera un usuario específico usando su **ID (UUID)**. | Basic Auth |
| **PUT** | `/users` | Update Existing User | Actualiza los datos de un usuario existente. **Requiere el ID en el cuerpo.** | Basic Auth |
| **DELETE** | `/users/{id}` | Delete User by ID | Elimina un usuario específico usando su **ID (UUID)**. | Basic Auth |

## Seguridad

Todos los endpoints requieren autenticación.

### Basic Authentication (BasicAuth)

Se requiere el uso del esquema de autenticación **HTTP Basic** en el header de la solicitud, utilizando un nombre de usuario (`BASIC_AUTH_USER`) y una contraseña (`BASIC_AUTH_PASS`) configurados.

## Esquemas de Datos

### UserResponse (Modelo de Respuesta)

Representa la estructura completa de un usuario. El campo `id` es generado por el sistema.

| Propiedad | Tipo | Formato | Descripción | Ejemplo |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `string` | `uuid` | ID único del usuario (Generado, read-only). | `8b73f80c-7b44-482a-89a9-3d19129e9d6d` |
| `name` | `string` | | Nombre completo del usuario. | `Jane Doe` |
| `username` | `string` | | Nombre de usuario único. | `janedoe123` |
| `email` | `string` | `email` | Correo electrónico único. | `jane.doe@example.com` |

### UserCreateRequest (Para POST /users)

Datos requeridos para la creación de un nuevo usuario.

| Propiedad | Tipo | Requerido | Descripción |
| :--- | :--- | :---: | :--- |
| `name` | `string` | **Sí** | Nombre completo del usuario. |
| `username` | `string` | **Sí** | Nombre de usuario único. |
| `email` | `string` | **Sí** | Correo electrónico único. |

### UserUpdate (Para PUT /users)

Datos requeridos y opcionales para la actualización. **Se debe incluir el `id`**.

| Propiedad | Tipo | Requerido | Descripción |
| :--- | :--- | :---: | :--- |
| `id` | `string` | **Sí** | ID del usuario a actualizar. |
| `name` | `string` | No | Nuevo nombre (opcional). |
| `username` | `string` | No | Nuevo nombre de usuario único (opcional). |
| `email` | `string` | No | Nuevo correo electrónico único (opcional). |

## Entornos de Servidores

La API está disponible en los siguientes entornos:

| URL | Descripción |
| :--- | :--- |
| `http://localhost:8080` | Servidor Local (Development) |
| `https://api.prod.user.com` | Servidor de Producción (ejemplo) |

## Códigos de Respuesta Comunes

Además de los códigos de éxito (200, 201, 204), la API utiliza los siguientes para manejar errores:

| Código | Descripción | Significado |
| :---: | :--- | :--- |
| **401** | Unauthorized | Fallo de autenticación. |
| **400** | Bad Request | El cuerpo de la petición es inválido o falló la validación. |
| **404** | Not Found | El recurso (usuario) solicitado no existe. |
| **409** | Conflict | Error de duplicidad (ej. `username` o `email` ya en uso). |
| **500** | Internal Server Error | Fallo inesperado en el procesamiento. |****