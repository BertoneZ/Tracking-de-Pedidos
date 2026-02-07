# Tracking de Pedidos - API Logística
API desarrollada en Go para el seguimiento de pedidos en tiempo real, utilizando PostgreSQL (PostGIS) para datos geoespaciales y Redis para el tracking de alta frecuencia.

## Tecnologías
Backend: Go  (Gin Gonic)

Base de Datos: PostgreSQL + PostGIS

Cache/Tracking: Redis (Geo commands)

Documentación: Swagger

Infraestructura: Docker & Docker Compose

## Instalación y Despliegue
No es necesario instalar Go o bases de datos localmente. El sistema está completamente contenedorizado.

Clonar el repositorio.

Asegurarse de tener Docker Desktop activo.

Ejecutar el siguiente comando en la terminal:

Bash
docker-compose up --build

## Documentación
Una vez levantado el servicio, podés acceder a la interfaz de Swagger para probar los endpoints: http://localhost:8080/swagger/index.html#/

## Seguridad
Autenticación basada en JWT.

Control de acceso por roles (RBAC): customer y driver.

Validación de propiedad: un cliente solo puede trackear sus propios pedidos.
