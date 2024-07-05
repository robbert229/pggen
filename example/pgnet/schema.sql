CREATE TABLE servers (
  id serial      PRIMARY KEY,
  ip_address     INET NOT NULL,
  extra_ip_address INET
);

