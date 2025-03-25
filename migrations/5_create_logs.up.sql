CREATE TABLE logs(
    id BIGSERIAL PRIMARY KEY,
    id_task BIGINT NOT NULL,
    date_of_operation DATE NOT NULL,
    info TEXT NOT NULL,
    FOREIGN KEY(id_task) REFERENCES tasks(id) ON DELETE CASCADE
);