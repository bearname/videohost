CREATE TABLE IF NOT EXISTS followers
(
    followingToUserId VARCHAR(255) NOT NULL,
    followerUserId    VARCHAR(255) NOT NULL,
    following         BOOL DEFAULT TRUE,
    CONSTRAINT FOREIGN KEY (followingToUserId) references users (key_user),
    CONSTRAINT FOREIGN KEY (followerUserId) references users (key_user),
    PRIMARY KEY (followerUserId, followingToUserId)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8;
