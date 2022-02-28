-- drops
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS miners;
DROP TABLE IF EXISTS blacklist;

--- miner table creation
CREATE TABLE miners (
    id bigserial NOT NULL unique PRIMARY KEY,
    authorized BOOLEAN,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);

INSERT INTO miners (authorized)
    VALUES (FALSE) returning id;
INSERT INTO miners (authorized)
    VALUES (FALSE) returning id;

--- rpc entries table creation
CREATE TABLE entries (
    id bigserial NOT NULL unique PRIMARY KEY,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    miner_id bigint NOT NULL,
    method VARCHAR(255) NOT NULL,
    params VARCHAR(1023) NOT NULL,
    ip VARCHAR(255) NOT NULL,
    success BOOLEAN,
    CONSTRAINT fk_miner
        FOREIGN KEY(miner_id) 
            REFERENCES miners(id)
);

--- subscription
CREATE TABLE subscriptions (
    id bytea PRIMARY KEY unique,
    miner_id bigint NOT NULL,
    notification_name VARCHAR(255) NOT NULL,
    subscribed BOOLEAN,
    CONSTRAINT fk_miner
        FOREIGN KEY(miner_id) 
            REFERENCES miners(id)
);

--- blacklist table creation
CREATE TABLE blacklist (
    id bigserial NOT NULL unique PRIMARY KEY,
    ip VARCHAR(255) NOT NULL
);
