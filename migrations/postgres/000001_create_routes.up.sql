CREATE TABLE IF NOT EXISTS routes (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL CHECK (length(trim(url)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_routes_updated_at
    ON routes(updated_at);

CREATE TABLE IF NOT EXISTS route_methods (
    route_id TEXT PRIMARY KEY,
    factory_key TEXT NOT NULL CHECK (length(trim(factory_key)) > 0),
    serialized_config BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_route_methods_route_id
    FOREIGN KEY (route_id)
    REFERENCES routes(id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_route_methods_updated_at
    ON route_methods(updated_at);