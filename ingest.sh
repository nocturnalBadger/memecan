#!/usr/bin/env bash

INPUT_FOLDER="$(realpath "$1")"
IMAGES_FOLDER=$PWD/images
DB_FILE=memes.db
DB_TABLE=Images

mkdir -p $IMAGES_FOLDER

cat schema.sql | sqlite3 $DB_FILE

for extension in .png .PNG .jpg .JPG .jpeg .JPEG; do
    cp -v "$INPUT_FOLDER"/*$extension "$INPUT_FOLDER"/**/*$extension "$IMAGES_FOLDER"
done;

for f in "$IMAGES_FOLDER"/*; do
    extension="${f##*.}"
    newExt=$(echo $extension | awk '{print tolower($0)}' | sed 's/jpeg/jpg/g')

    fileHash=$(md5sum "$f" | cut -d ' ' -f 1)
    newFilename="$IMAGES_FOLDER"/$fileHash.$newExt
    mv -v "$f" $newFilename

    insertQuery="INSERT INTO $DB_TABLE (hash, file_path) VALUES ('$fileHash', '$newFilename')
                 ON CONFLICT DO NOTHING"
    sqlite3 $DB_FILE "$insertQuery"
done;
