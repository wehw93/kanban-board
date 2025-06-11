CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    encrypted_password TEXT NOT NULL
);

CREATE TABLE projects(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    id_creator BIGINT NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY(id_creator) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE columns(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    id_project BIGINT NOT NULL,
    FOREIGN KEY (id_project) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE tasks(
    id BIGSERIAL PRIMARY KEY,
    id_column BIGINT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    date_of_create DATE NOT NULL,
    date_of_execution DATE,
    id_executor BIGINT,
    id_creator BIGINT NOT NULL,
    status TEXT NOT NULL ,
    FOREIGN KEY(id_column) REFERENCES columns(id) ON DELETE CASCADE,
    FOREIGN KEY(id_executor) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(id_creator) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE logs(
    id BIGSERIAL PRIMARY KEY,
    id_task BIGINT NOT NULL,
    date_of_operation DATE NOT NULL,
    info TEXT NOT NULL,
    FOREIGN KEY(id_task) REFERENCES tasks(id) ON DELETE CASCADE
);