create extension pgcrypto;
select * from pg_available_extensions;

CREATE TABLE session(
  session_id VARCHAR(128) PRIMARY KEY,
  account_id INTEGER NOT NULL,
  session_start VARCHAR(64) NOT NULL,
  session_end VARCHAR(64),
  FOREIGN KEY(account_id) REFERENCES ACCOUNT(account_id)
);


CREATE TABLE account(
  account_id serial PRIMARY KEY NOT NULL,
  username VARCHAR (50) UNIQUE NOT NULL,
  password VARCHAR (50) NOT NULL
);

