CREATE TABLE Images (
    hash char(32) NOT NULL PRIMARY KEY,
    file_path varchar(255)
);

CREATE TABLE Text (
    id INTEGER PRIMARY KEY,
    image char(32),
    text VARCHAR(65535),
    FOREIGN KEY (image) REFERENCES Images(hash)
);
