-- FindServers finds all servers.
-- name: FindServers :many
SELECT * FROM servers;

-- FindServerByIP finds a server by its ip address.
-- name: FindServerByIP :one
SELECT * FROM servers WHERE ip_address = pggen.arg('ip_address');

-- InsertServer inserts a server and returns the ID.
-- name: InsertServer :one
INSERT INTO  servers (ip_address, extra_ip_address)
VALUES (pggen.arg('ip_address'), pggen.arg('extra_ip_address'))
RETURNING id;
