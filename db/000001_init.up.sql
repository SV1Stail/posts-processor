CREATE TABLE IF NOT EXISTS original_posts (
    "id" VARCHAR NOT NULL PRIMARY KEY,
    "data" JSONB NOT NULL,
    "link_new_post" VARCHAR,
    "original_channel" VARCHAR,
    "url" VARCHAR,
    "original_image_url" VARCHAR,
    "theme" VARCHAR NOT NULL
    "created_at" TIMESTAMP DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMP DEFAULT NOW() NOT NULL,
    "attempts" INT DEFAULT 0 NOT NULL CHECK("attempts" >= 0)
);

CREATE INDEX IF NOT EXISTS idx_original_posts ON original_posts("theme");
