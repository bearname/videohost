create table if not exists video.users
(
    key_user      varchar(255)                           not null
        primary key,
    username      varchar(255)                           not null,
    password      blob                                   not null,
    email         varchar(255)                           not null,
    created       datetime     default CURRENT_TIMESTAMP null,
    role          int          default 0                 not null,
    secret        varchar(255) default 'sdmalncnjsdsmf'  not null,
    access_token  varchar(255) default ''                not null,
    refresh_token varchar(255) default '123'             not null,
    isSubscribed  tinyint(1)                             not null,
    constraint username
        unique (username)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

INSERT INTO video.users (key_user, username, password, email, created, role, secret, access_token, refresh_token,
                         isSubscribed)
VALUES ('2b3d1b76-bc6a-11eb-8b34-e4e74940035b', 'admin',
        0x24326124313024686970515337564A614353332F7033544B74734269755865574E76344E30474A61574356353964633359654C46637A2F642E78652E,
        'admin@gmail.com', '2021-05-24 11:29:34', 0, 'sdmalncnjsdsmf',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjIzMjcyMTIsInJvbGUiOjAsInVzZXJJZCI6IjJiM2QxYjc2LWJjNmEtMTFlYi04YjM0LWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoiYWRtaW4ifQ.jdjyqazXmuklpHGyv351lljlqPnUTQXs3jaZrHdGe5U',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE5MzcyMDQ5NzQsInJvbGUiOjEsInVzZXJJZCI6IjJiM2QxYjc2LWJjNmEtMTFlYi04YjM0LWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoiYWRtaW4ifQ.AIJnoxwVudeqFuI8EkB61kp5Nawisl_y8AL8H5Seij8',
        0);
INSERT INTO video.users (key_user, username, password, email, created, role, secret, access_token, refresh_token,
                         isSubscribed)
VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'user',
        0x24326124313024503155386B746C33536B44472F57303255396C4955754442564362777A534E6E79384E764F7571732F775733667A6A6A617337564F,
        '', '2021-05-20 09:19:13', 1, 'sdmalncnjsdsmf',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjE2Nzg3MTEsInJvbGUiOjEsInVzZXJJZCI6IjRiYmU3NjUzLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoidXNlciJ9.isBL0rFPidyKxPIlBUEtnQ8_wk_0GAtEbIrOiviP7ag',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE5MzY4NTE1NTMsInJvbGUiOjEsInVzZXJJZCI6IjRiYmU3NjUzLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoidXNlciJ9.FEzXwTANvfuEDlAC-YLpPqLIBNZjhj5IhpCo8oWeR9I',
        0);
INSERT INTO video.users (key_user, username, password, email, created, role, secret, access_token, refresh_token,
                         isSubscribed)
VALUES ('cfaff592-b933-11eb-b073-e4e74940035b', 'mikha',
        0x243261243130244854504441545036514D384D367961546366446C302E35656761636F70346956674F706E3555765761666B3334355971687058412E,
        'mihail12russ@gmail.com', '2021-05-20 09:22:54', 1, 'sdmalncnjsdsmf',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjIzMjU4NjIsInJvbGUiOjEsInVzZXJJZCI6ImNmYWZmNTkyLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoibWlraGEifQ.p2fJRU5hIPSR0Im7bF3NdaZr15JlaW7ltJOru6hi52g',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE5MzY4NTE3NzQsInJvbGUiOjEsInVzZXJJZCI6ImNmYWZmNTkyLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoibWlraGEifQ.MslmmKdQtAYCAxHX3AcZXwPG2LneVYChriGEv4ns6dY',
        1);