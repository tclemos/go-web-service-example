BEGIN;

CREATE TABLE things (
    code UUID NOT NULl,
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY (code)
);

END;