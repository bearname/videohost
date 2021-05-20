create table users
(
    key_user      varchar(255)                           not null
        primary key,
    username      varchar(255)                           not null,
    password      blob                                   not null,
    created       datetime     default CURRENT_TIMESTAMP null,
    role          int          default 0                 not null,
    secret        varchar(255) default 'sdmalncnjsdsmf'  not null,
    access_token  varchar(255) default ''                not null,
    refresh_token varchar(255) default '123'             not null,
    constraint username
        unique (username)
)
    charset = utf8;

INSERT INTO video.users (key_user, username, password, created, role, secret, access_token, refresh_token) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'user', 0x24326124313024503155386B746C33536B44472F57303255396C4955754442564362777A534E6E79384E764F7571732F775733667A6A6A617337564F, '2021-05-20 09:19:13', 1, 'sdmalncnjsdsmf', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjE1MDA4ODMsInJvbGUiOjEsInVzZXJJZCI6IjRiYmU3NjUzLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoidXNlciJ9.-8xVYFRYteVFGNFrVY1b3LGB2qkiBaRBsA7tB_-20j4', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE5MzY4NTE1NTMsInJvbGUiOjEsInVzZXJJZCI6IjRiYmU3NjUzLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoidXNlciJ9.FEzXwTANvfuEDlAC-YLpPqLIBNZjhj5IhpCo8oWeR9I');
INSERT INTO video.users (key_user, username, password, created, role, secret, access_token, refresh_token) VALUES ('cfaff592-b933-11eb-b073-e4e74940035b', 'mikha', 0x243261243130244854504441545036514D384D367961546366446C302E35656761636F70346956674F706E3555765761666B3334355971687058412E, '2021-05-20 09:22:54', 1, 'sdmalncnjsdsmf', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjE0OTE4OTQsInJvbGUiOjEsInVzZXJJZCI6ImNmYWZmNTkyLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoibWlraGEifQ.MqKwQHY8iEXeCXOHa_tPOHzS4AdEvUPvNq-TfTTRuJw', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE5MzY4NTE3NzQsInJvbGUiOjEsInVzZXJJZCI6ImNmYWZmNTkyLWI5MzMtMTFlYi1iMDczLWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoibWlraGEifQ.MslmmKdQtAYCAxHX3AcZXwPG2LneVYChriGEv4ns6dY');
INSERT INTO video.users (key_user, username, password, created, role, secret, access_token, refresh_token) VALUES ('d4b0c1f5-b8fd-11eb-a138-e4e74940035b', 'admin', 0x243261243130243074794550323263635230682E63722E39734E554C654C474B704A6B647959514F41497A39464255476A74416C6D6B715435737153, '2021-05-20 02:56:30', 0, 'sdmalncnjsdsmf', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjE0Njk0OTAsInJvbGUiOjEsInVzZXJJZCI6ImQ0YjBjMWY1LWI4ZmQtMTFlYi1hMTM4LWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoiYWRtaW4ifQ.gppeSF6oEd1kLgIC2zylDtiPMlfhwgxSrIP97vypSv4', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE5MzY4NDc0MDcsInJvbGUiOjEsInVzZXJJZCI6IjNlYWVlZjQ3LWI5MjgtMTFlYi1hZjU0LWU0ZTc0OTQwMDM1YiIsInVzZXJuYW1lIjoibWlraGExIn0.VHqlyEXRGbeEDm6OZWcO7b111LRw1xtLboh_JILZmpg');