CREATE TABLE IF NOT EXISTS subtitle
(
    `id`       INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `video_id` VARCHAR(255) UNIQUE NOT NULL,
    CONSTRAINT FOREIGN KEY (video_id) REFERENCES video.video (id_video) ON DELETE CASCADE
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS subtitle_part
(
    `id`        INT(11) UNSIGNED AUTO_INCREMENT primary key,
    subtitle_id INT(11) UNSIGNED NOT NULL,
    `start`     INT(11) UNSIGNED NOT NULL,
    `end`       INT(11) UNSIGNED NOT NULL,
    `text`      TEXT             NOT NULL,
    CONSTRAINT FOREIGN KEY (subtitle_id) REFERENCES video.subtitle (id) ON DELETE CASCADE
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;

CREATE INDEX IN_video_id ON subtitle (video_id);