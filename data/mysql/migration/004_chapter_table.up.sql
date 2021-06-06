CREATE TABLE IF NOT EXISTS video.video_chapter
(
    id       INT AUTO_INCREMENT PRIMARY KEY,
    id_video VARCHAR(255) NOT NULL,
    title    VARCHAR(255) NOT NULL,
    start    INT          NOT NULL,
    end      INT          NOT NULL,
    CONSTRAINT FOREIGN KEY (id_video) REFERENCES video (id_video)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;