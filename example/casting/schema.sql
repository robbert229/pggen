CREATE TABLE author (
  author_id  serial PRIMARY KEY,
  first_name text NOT NULL,
  last_name  text NOT NULL,
  suffix text NULL
);
 
CREATE TABLE book (
  book_id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  author_id INTEGER,
  FOREIGN KEY (author_id) REFERENCES author(author_id)
);
