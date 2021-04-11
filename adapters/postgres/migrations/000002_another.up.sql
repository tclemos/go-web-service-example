BEGIN;

CREATE TABLE another_things (
    id INT NOT NULL,
    code UUID NOT NULl,
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY (id)
);

END;