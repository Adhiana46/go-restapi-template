CREATE SEQUENCE activity_group_seq;

CREATE TABLE activity_group
(
	id INT NOT NULL DEFAULT NEXTVAL ('activity_group_seq'),
	uuid CHAR(36) NOT NULL UNIQUE,
	name VARCHAR(255),
	description TEXT,
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

INSERT INTO activity_group
(uuid, name, description)
VALUES
('b4b56351-5e98-4793-aad0-e7ed8911b91f', 'Activity 1', 'ini deskrpisio dari activity 1'),
('9a89dcac-ce5a-41e2-9337-797ca8001932', 'Activity 2', 'ini deskrpisio dari activity 2');