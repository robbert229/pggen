-- FindBooks finds all books.
-- name: FindBooks :many
SELECT a::author, b::book
FROM book b LEFT JOIN author a ON (b.author_id = a.author_id);
-- InsertAuthor inserts an author by name and returns the ID.
-- name: InsertAuthor :one
INSERT INTO author (first_name, last_name)
VALUES (pggen.arg('FirstName'), pggen.arg('LastName'))
RETURNING author_id;

-- InsertBook inserts a book.
-- name: InsertBook :one
INSERT INTO book (title)
VALUES (pggen.arg('Title'))
RETURNING book_id;

-- AssignAuthor assigns an author to a book.
-- name: AssignAuthor :exec
UPDATE book
SET author_id = pggen.arg('AuthorID')
WHERE book_id = pggen.arg('BookID');
