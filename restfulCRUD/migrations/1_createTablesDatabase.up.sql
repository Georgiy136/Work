
CREATE TABLE IF NOT EXISTS roles (
    ID SERIAL PRIMARY KEY,
    role character varying NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS rights (
    ID SERIAL PRIMARY KEY,
    rights character varying NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS roles_rights (
   role_id INT,
   right_id INT,
   FOREIGN KEY (role_id) REFERENCES roles (ID),
   FOREIGN KEY (right_id) REFERENCES rights (ID),
   CONSTRAINT pk_roles_rights PRIMARY KEY (role_id, right_id)
);

CREATE TABLE IF NOT EXISTS clients
(
    uuid character varying PRIMARY KEY,
    username character varying NOT NULL UNIQUE,
    login character varying NOT NULL UNIQUE,
    password character varying NOT NULL,
    role character varying NOT NULL,
    FOREIGN KEY (role) REFERENCES roles (role)
);

CREATE TABLE IF NOT EXISTS operators
(
    uuid character varying PRIMARY KEY,
    first_name character varying NOT NULL,
    last_name character varying NOT NULL,
    patronymic character varying NOT NULL,
    city character varying NOT NULL,
    phone character varying NOT NULL,
    email character varying NOT NULL,
    password character varying NOT NULL
);

CREATE TABLE IF NOT EXISTS project_types
(
    ID SERIAL PRIMARY KEY,
    project_type character varying NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS projects
(
    uuid character varying PRIMARY KEY,
    project_name character varying NOT NULL,
    project_type int NOT NULL,
    operators  character varying ARRAY,
    FOREIGN KEY (project_type) REFERENCES project_types (ID)
);


INSERT INTO roles (
	role
)
VALUES (
	'Admin'
);

INSERT INTO roles (
	role
)
VALUES (
	'User'
);



INSERT INTO rights (
	rights
)
VALUES (
	'GET'
);
INSERT INTO rights (
	rights
)
VALUES (
	'PUT'
);
INSERT INTO rights (
	rights
)
VALUES (
	'DELETE'
);
INSERT INTO rights (
	rights
)
VALUES (
	'POST'
);




INSERT INTO project_types
    (
    project_type
    )
VALUES
    (
        'Входящий'
);

INSERT INTO project_types
    (
    project_type
    )
VALUES
    (
        'Исходящий'
);

INSERT INTO project_types
    (
    project_type
    )
VALUES
    (
        'Автоинформатор'
);


INSERT INTO roles_rights VALUES (1,1),(1,2),(1,3),(1,4),(2,1);


CREATE INDEX idx_client_ID on clients (uuid);

CREATE INDEX idx_operator_ID on operators (uuid);

CREATE INDEX idx_project_ID on projects (uuid);