import sqlite3
import io
import sys
import os
import logging
from google.cloud import vision

API_MAX_BYTES = 10485760

logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
logger.addHandler(logging.StreamHandler(sys.stdout))


def detect_text(path):
    """Detects text in the file."""
    if os.path.getsize(path) > API_MAX_BYTES:
        logger.warning(f"Image {path} is too large to run text detection.")
        return ""

    client = vision.ImageAnnotatorClient()

    with io.open(path, 'rb') as image_file:
        content = image_file.read()

    image = vision.types.Image(content=content)
    response = client.text_detection(image=image)

    texts = response.text_annotations

    main_text = texts[0].description if texts else ""
    logger.debug(f"Text found for {path}: {main_text}")

    return main_text


def get_all_image_text():
    conn = sqlite3.connect('memes.db')
    cursor = conn.cursor()

    cursor.execute("SELECT * FROM Images")

    image_hash, image_file = cursor.fetchone()

    images = cursor.execute("SELECT * from Images")
    for image_hash, image_file in cursor.fetchall():
        logger.debug(f"Checking for text entries for {image_file}")

        cursor.execute(f"SELECT * from Text where image='{image_hash}'")
        if not cursor.fetchone():
            logger.debug(f"No entry text found for {image_file}. Querying Google Vision API.")
            image_text = detect_text(image_file)
            cursor.execute("INSERT INTO Text (image, text) VALUES (?, ?)", (image_hash, image_text))
            conn.commit()
        else:
            logger.debug(f"Found existing text for {image_file}")

    conn.close()

get_all_image_text()
