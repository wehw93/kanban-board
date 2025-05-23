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