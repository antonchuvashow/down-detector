import os

FEATURE_FLAGS = {"EMBEDDED_SUPERSET": True}
SESSION_COOKIE_SAMESITE = None

ENABLE_CORS = True
CORS_OPTIONS = {
    "supports_credentials": True,
    "allow_headers": "*",
    "expose_headers": "*",
    "resources": "*",
    "origins": os.environ.get("CORS_OPTIONS_ORIGINS", "*").split(",")
}
print(CORS_OPTIONS)
GUEST_ROLE_NAME = "EmbedGuest"
# GUEST_TOKEN_JWT_SECRET = "YOUR_SECURE_RANDOM_JWT_SECRET_KEY"
# GUEST_TOKEN_JWT_ALGO = "HS256"
# GUEST_TOKEN_HEADER_NAME = "X-GuestToken"
# GUEST_TOKEN_JWT_EXP_SECONDS = 300

OVERRIDE_HTTP_HEADERS = {
    'X-Frame-Options': 'ALLOWALL'
}
TALISMAN_ENABLED = False

SECRET_KEY = os.environ["SUPERSET_SECRET_KEY"]

SQLALCHEMY_DATABASE_URI = os.environ["SQLALCHEMY_DATABASE_URI"]

RATELIMIT_STORAGE_URI = os.environ.get(
    "RATELIMIT_STORAGE_URI",
    "redis://redis:6379/4",
)

CACHE_CONFIG = {
    "CACHE_TYPE": "RedisCache",
    "CACHE_REDIS_HOST": os.environ.get("REDIS_HOST", "redis"),
    "CACHE_REDIS_PORT": int(os.environ.get("REDIS_PORT", 6379)),
    "CACHE_REDIS_DB": int(os.environ.get("CACHE_REDIS_DB", 1)),
    "CACHE_DEFAULT_TIMEOUT": 300,
    "CACHE_KEY_PREFIX": "superset_cache_",
}

DATA_CACHE_CONFIG = CACHE_CONFIG

FILTER_STATE_CACHE_CONFIG = {
    "CACHE_TYPE": "RedisCache",
    "CACHE_REDIS_HOST": os.environ.get("REDIS_HOST", "redis"),
    "CACHE_REDIS_PORT": int(os.environ.get("REDIS_PORT", 6379)),
    "CACHE_REDIS_DB": 2,
    "CACHE_DEFAULT_TIMEOUT": 86400,
    "CACHE_KEY_PREFIX": "superset_filter_",
}

EXPLORE_FORM_DATA_CACHE_CONFIG = {
    "CACHE_TYPE": "RedisCache",
    "CACHE_REDIS_HOST": os.environ.get("REDIS_HOST", "redis"),
    "CACHE_REDIS_PORT": int(os.environ.get("REDIS_PORT", 6379)),
    "CACHE_REDIS_DB": 3,
    "CACHE_DEFAULT_TIMEOUT": 86400,
    "CACHE_KEY_PREFIX": "superset_explore_",
}