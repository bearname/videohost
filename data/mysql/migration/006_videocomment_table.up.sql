CREATE TABLE IF NOT EXISTS video_comments
(
    `id`        INT(11) UNSIGNED AUTO_INCREMENT primary key,
    `video_id`  VARCHAR(255) NOT NULL,
    `user_id`   VARCHAR(255) NOT NULL,
    `parent_id` INT(11) UNSIGNED DEFAULT NULL,
    `message`   TEXT         NOT NULL,
    `created`   TIMESTAMP        DEFAULT NOW(),
    CONSTRAINT FOREIGN KEY (parent_id) REFERENCES video_comments (id)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;