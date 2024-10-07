CREATE TABLE IF NOT EXISTS "songs" (
    "id" SERIAL PRIMARY KEY,
    "song" VARCHAR(255) NOT NULL,
    "group" VARCHAR(255) NOT NULL,
    "release_date" DATE,
    "link" TEXT
);

CREATE TABLE IF NOT EXISTS "lyrics" (
    "id" SERIAL PRIMARY KEY,
    "text" TEXT,
    "song_id" INT NOT NULL
);


DO $$
BEGIN
   IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_songs_by_id'
   ) THEN
        ALTER TABLE "lyrics"
        ADD CONSTRAINT "fk_songs_by_id"
        FOREIGN KEY ("song_id") REFERENCES "songs"("id")
        ON DELETE CASCADE;
   END IF;
END $$;


DO $$
BEGIN
   IF NOT EXISTS (
      SELECT 1
      FROM pg_constraint
      WHERE conname = 'unique_song_id'
   ) THEN
      ALTER TABLE "lyrics"
      ADD CONSTRAINT "unique_song_id" UNIQUE ("song_id");
   END IF;
END $$;
