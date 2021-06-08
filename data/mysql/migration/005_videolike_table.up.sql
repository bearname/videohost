CREATE TABLE if NOT EXISTS video.video_like
(
    id_video VARCHAR(255)                       NOT NULL ,
    liked    DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ,
    isLike   BOOL     DEFAULT TRUE,
    owner_id VARCHAR(255)                       NOT NULL ,
    CONSTRAINT FOREIGN KEY (owner_id) REFERENCES users (key_user),
    CONSTRAINT FOREIGN KEY (id_video) REFERENCES video (id_video),
    PRIMARY KEY (owner_id,id_video)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;