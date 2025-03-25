CREATE TABLE projects(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    id_creator BIGINT NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY(id_creator) REFERENCES users(id) ON DELETE CASCADE
);