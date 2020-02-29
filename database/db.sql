CREATE TABLE users
(
    id       SERIAL PRIMARY KEY UNIQUE NOT NULL,
    username varchar(50) UNIQUE        NOT NULL,
    password varchar(97)               NOT NULL
);

INSERT INTO users(username, password)
VALUES ('admin',
        '$argon2id$v=19$m=65536,t=3,p=2$UDQbl6KSK7h0nzhb/eAAtw$e3ZT9jPEprmoP1A40wD04a5s0J5/2otW1gOdny5M9/E'); -- admin, testadminpassword