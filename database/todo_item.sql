CREATE SEQUENCE todo_item_seq;

CREATE TABLE todo_item
(
	id INT NOT NULL DEFAULT NEXTVAL ('todo_item_seq'),
	uuid CHAR(36) NOT NULL UNIQUE,
    activity_id INT NOT NULL,
	name VARCHAR(255),
	description TEXT,
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (activity_id) REFERENCES activity_group(id) ON DELETE CASCADE,

	PRIMARY KEY (id)
);