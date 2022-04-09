-- CreateTable
CREATE TABLE "stream" (
    "id" SERIAL NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "url" TEXT NOT NULL,

    CONSTRAINT "stream_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "track" (
    "id" SERIAL NOT NULL,
    "cast_id" INTEGER NOT NULL,
    "started_at" TIMESTAMP(3) NOT NULL,
    "ended_at" TIMESTAMP(3) NOT NULL,
    "title" TEXT NOT NULL,
    "listeners" INTEGER NOT NULL,

    CONSTRAINT "track_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "stream_url_key" ON "stream"("url");

-- CreateIndex
CREATE UNIQUE INDEX "track_cast_id_title_started_at_key" ON "track"("cast_id", "title", "started_at");

-- AddForeignKey
ALTER TABLE "track" ADD CONSTRAINT "track_cast_id_fkey" FOREIGN KEY ("cast_id") REFERENCES "stream"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
