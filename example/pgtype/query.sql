-- FindAuthorById finds one (or zero) authors by ID.
-- name: FindAuthorByID :one
SELECT * FROM author WHERE author_id = pggen.arg('AuthorID');

-- InsertAuthor inserts an author by name and returns the ID.
-- name: InsertAuthor :one
INSERT INTO author (first_name, last_name)
VALUES (pggen.arg('FirstName'), pggen.arg('LastName'))
RETURNING author_id;

-- name: SelectInt8 :one
SELECT 1234567890123456789::int8 AS my_integer;

-- name: SelectInt4 :one
SELECT 12345::int4 AS my_integer;