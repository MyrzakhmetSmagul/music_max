CREATE TABLE "songs" (
    "id" SERIAL PRIMARY KEY,
    "song" VARCHAR(255) NOT NULL,
    "group" VARCHAR(255) NOT NULL,
    "release_date" DATE,
    "link" TEXT
);

CREATE TABLE "lyrics" (
    "id" SERIAL PRIMARY KEY,
    "text" TEXT,
    "song_id" INT NOT NULL
);

ALTER TABLE "lyrics"
ADD CONSTRAINT "fk_songs_by_id"
FOREIGN KEY ("song_id") REFERENCES "songs"("id")
ON DELETE CASCADE;

ALTER TABLE "lyrics"
ADD CONSTRAINT "unique_song_id" UNIQUE("song_id");