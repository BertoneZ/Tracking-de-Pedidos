-- Extensión para manejar coordenadas GPS
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Usuarios (Para Login)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) CHECK (role IN ('driver', 'customer')) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Repartidores (Datos espaciales)
CREATE TABLE IF NOT EXISTS drivers (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT false,
    last_location GEOGRAPHY(Point, 4326), 
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Pedidos
CREATE TYPE order_status AS ENUM ('PENDING', 'ASSIGNED', 'PICKED_UP', 'DELIVERED', 'CANCELLED');

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES users(id),
    driver_id UUID REFERENCES users(id),
    status order_status DEFAULT 'PENDING',
    origin GEOGRAPHY(Point, 4326) NOT NULL,
    destination GEOGRAPHY(Point, 4326) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índice para búsquedas geográficas rápidas
CREATE INDEX IF NOT EXISTS idx_drivers_location ON drivers USING GIST(last_location);