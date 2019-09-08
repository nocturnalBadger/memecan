# MEMECAN

A tool for stashing and indexing memes.

## Prerequisites
Python 3.7  
Google Vision API Python library (google-cloud-vision)


## Usage
Run ingest.sh to input images:
```bash
./ingest.sh <path to memes>
```
The script will recursively copy all jpgs and pngs from the input folder into an images folder in the current working directory. They will then be assigned a file name based on their md5 sum.

A sqlite3 database is created to map each image to its file path.

Next, run get_text.py
```bash
python get_text.py
```

This will run over all the images previously ingested and check if it has an associated entry in the 'Text' table of the database. If it doesn't, it will query the text detection api and insert the resulting text.

## TODO

- [ ] gif support (at least for stashing not text recognition)
- [ ] Tagging (manual or automatic) e.g. filter by "reaction" or "comic"
- [ ] Detect / sort by meme template from knowyourmeme.com
- [ ] Some kind of interface. At least nicer cli. Maybe a UI if I'm feeling it.
