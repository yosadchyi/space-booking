CREATE TABLE launchpad
(
    id CHAR(24) NOT NULL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    status VARCHAR(32) NOT NULL
);

CREATE TABLE launch
(
    id UUID NOT NULL PRIMARY KEY,
    launchpad_id CHAR(24) NOT NULL,
    date DATE NOT NULL
);

CREATE UNIQUE INDEX launch_launchpad_id_date_idx ON launch(launchpad_id, date);

CREATE TABLE destination
(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

INSERT INTO destination (id, name) VALUES
    ('b11cadee-e907-4e55-8bef-be92aaad940c', 'Mars'),
    ('135273de-ed80-4577-8157-19af8cedbf03', 'Moon'),
    ('b805aa0a-b14b-41b2-b927-d9a6e1fa22a9', 'Pluto'),
    ('c4db5f9d-82ad-4914-9a2a-63432ed0b3a1', 'Asteriod Belt'),
    ('fe8d4bfb-8989-467d-b4c7-07091b3266ce', 'Europa'),
    ('564c14da-5717-438c-b7df-a699e0caaa6c', 'Titan'),
    ('1e255888-bb39-44c7-9208-6a673eafa264', 'Ganymede');

CREATE TABLE booking
(
    id UUID NOT NULL PRIMARY KEY,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    birthday DATE NOT NULL,
    launch_date DATE NOT NULL,
    launchpad_id CHAR(24) NOT NULL REFERENCES launchpad(id),
    destination_id UUID NOT NULL REFERENCES destination(id)
);
