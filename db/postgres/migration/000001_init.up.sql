CREATE TABLE IF NOT EXISTS original_posts (
    "id" VARCHAR NOT NULL PRIMARY KEY,
    "data" JSONB NOT NULL,
    "link_original_post" VARCHAR,
    "original_channel" VARCHAR,
    "original_image_url" VARCHAR,
    "theme" VARCHAR NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMP DEFAULT NOW() NOT NULL,
    "attempts" INT DEFAULT 0 NOT NULL CHECK("attempts" >= 0)
);

CREATE INDEX IF NOT EXISTS idx_original_posts_theme ON original_posts("theme");

CREATE TABLE IF NOT EXISTS system_prompts (
    "id" SERIAL PRIMARY KEY,
    "owner_name" VARCHAR NOT NULL,
    "title_name" VARCHAR NOT NULL,
    "message" JSONB,
    "created" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
