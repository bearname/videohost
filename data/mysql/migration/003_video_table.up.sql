create table video.video
(
    id_video      varchar(255)                                not null primary key,
    title         varchar(255)                                not null,
    uploaded      datetime     default CURRENT_TIMESTAMP      null,
    duration      int          default -1                     null,
    status        int          default 1                      null,
    thumbnail_url varchar(255) default 'content\\default.jpg' null,
    url           varchar(255)                                not null,
    description   text                                        null,
    quality       varchar(255) default ''                     null,
    views         int          default 0                      null,
    owner_id      varchar(255)                                not null,
    constraint video_users_key_user_fk
        foreign key (owner_id) references users (key_user)
)
    charset = utf8;

create fulltext index IN_video_title on video.video (title);

INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('1', 'Black Mirror trailer season 4', '2021-05-27 10:47:19', -1, 2, 'content\\default.jpg', '1',
        'In an abstrusely dystopian future, several individuals grapple with the manipulative effects of cutting edge technology in their personal lives and behaviours.',
        ',2160', 0, 'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('3624f091-bcac-11eb-afa0-e4e74940035b', 'xmen adaptive', '2021-05-24 19:22:20', 320, 3,
        'content\\3624f091-bcac-11eb-afa0-e4e74940035b\\screen.jpg',
        'content\\3624f091-bcac-11eb-afa0-e4e74940035b\\index.mp4', 'pexels-polina-tankilevitch-6646568', '', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('53f81f3b-c07f-11eb-bc58-e4e74940035b', 'black mirror season 4 episode titles', '2021-05-29 16:11:07', 50, 3,
        'content\\53f81f3b-c07f-11eb-bc58-e4e74940035b\\screen.jpg',
        'content\\53f81f3b-c07f-11eb-bc58-e4e74940035b\\index.mp4', 'black mirror season 4 episode titles',
        ',1080,720,480,360', 0, 'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('6da13e9b-bca7-11eb-afa0-e4e74940035b', 'andy barbourasas', '2021-05-28 07:50:00', 15, 3,
        'content\\2894ca86-bf70-11eb-b944-e4e74940035b\\screen.jpg',
        'content\\2894ca86-bf70-11eb-b944-e4e74940035b\\index.mp4', 'andy barbour smile', ',1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('8cbb79f0-bc76-11eb-afc7-e4e74940035b', '1234 pexels ', '2021-05-24 12:58:12', 320, 3,
        'content\\8cbb79f0-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\8cbb79f0-bc76-11eb-afc7-e4e74940035b\\index.mp4', '1234 pexels ', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('9bb38b2c-bcac-11eb-afa0-e4e74940035b', 'pexels-polina-tankilevitch-6646568', '2021-05-24 19:25:11', 320, 3,
        'content\\9bb38b2c-bcac-11eb-afa0-e4e74940035b\\screen.jpg',
        'content\\9bb38b2c-bcac-11eb-afa0-e4e74940035b\\index.mp4', 'pexels-polina-tankilevitch-6646568',
        '1080,720,480,360', 0, 'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('9be6d063-bc76-11eb-afc7-e4e74940035b', 'Pexels Videos 1443653', '2021-05-24 12:58:38', 22, 3,
        'content\\9be6d063-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\9be6d063-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'Pexels Videos 1443653', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('9f46f80b-c07f-11eb-bc58-e4e74940035b', 'black mirror season 4 episode titles', '2021-05-29 16:13:13', 50, 3,
        'content\\9f46f80b-c07f-11eb-bc58-e4e74940035b\\screen.jpg',
        'content\\9f46f80b-c07f-11eb-bc58-e4e74940035b\\index.mp4', 'black mirror season 4 episode titles',
        ',1080,720,480,360', 0, 'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('a131307e-bcaa-11eb-afa0-e4e74940035b', 'xmen adaptive', '2021-05-24 19:11:00', 13, 3,
        'content\\a131307e-bcaa-11eb-afa0-e4e74940035b\\screen.jpg',
        'content\\a131307e-bcaa-11eb-afa0-e4e74940035b\\index.mp4', 'pexels-polina-tankilevitch-6646568', '1080', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('a7e608d9-bc76-11eb-afc7-e4e74940035b', 'Pexels Videos 1583096', '2021-05-24 12:58:57', 13, 3,
        'content\\a7e608d9-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\a7e608d9-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'Pexels Videos 1583096', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('aa6f3609-befc-11eb-a50d-e4e74940035b', 'andy barbour', '2021-05-27 18:03:17', 15, 3,
        'content\\aa6f3609-befc-11eb-a50d-e4e74940035b\\screen.jpg',
        'content\\aa6f3609-befc-11eb-a50d-e4e74940035b\\index.mp4', 'andy barbour smile', ',1440,1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('afe67012-bc76-11eb-afc7-e4e74940035b', 'Black Mirror trailer', '2021-05-24 12:59:11', 150, 3,
        'content\\afe67012-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\afe67012-bc76-11eb-afc7-e4e74940035b\\index.mp4',
        'In an abstrusely dystopian future, several individuals grapple with the manipulative effects of cutting edge technology in their personal lives and behaviours.',
        '1080,720,480,360', 0, 'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('b6302ea4-bc76-11eb-afc7-e4e74940035b', 'pexels-andy-barbour-5337311', '2021-05-24 12:59:22', 15, 3,
        'content\\b6302ea4-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\b6302ea4-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'pexels-andy-barbour-5337311', '1080,720,480,360',
        0, 'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('bc6a1c33-bc76-11eb-afc7-e4e74940035b', 'Pexels Videos 2046575', '2021-05-24 12:59:32', 30, 3,
        'content\\bc6a1c33-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\bc6a1c33-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'Pexels Videos 2046575', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('bf915b4f-c07f-11eb-bc58-e4e74940035b', 'black mirror season 4 episode titles', '2021-05-29 16:14:07', -1, 2,
        'content\\default.jpg', 'content\\bf915b4f-c07f-11eb-bc58-e4e74940035b\\index.mp4',
        'black mirror season 4 episode titles', ',1080,720,480,360', 0, 'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('c4246118-bc76-11eb-afc7-e4e74940035b', 'Pexels Videos 4708', '2021-05-24 12:59:45', 30, 3,
        'content\\c4246118-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\c4246118-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'Pexels Videos 4708', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('d263d26b-bc76-11eb-afc7-e4e74940035b', 'video', '2021-05-24 13:00:09', 10, 3,
        'content\\d263d26b-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\d263d26b-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'video', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('dbcc692c-bc76-11eb-afc7-e4e74940035b', 'Pexels Videos 2019791', '2021-05-24 13:00:25', 20, 3,
        'content\\dbcc692c-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\dbcc692c-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'Pexels Videos 2019791', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('e712ec28-bc76-11eb-afc7-e4e74940035b', 'Pexels Videos 2019781', '2021-05-24 13:00:44', 19, 3,
        'content\\e712ec28-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\e712ec28-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'Pexels Videos 2019781', '1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('e7303069-bf6d-11eb-b944-e4e74940035b', 'andy barbourasas', '2021-05-28 07:33:52', 15, 3,
        'content\\e7303069-bf6d-11eb-b944-e4e74940035b\\screen.jpg',
        'content\\e7303069-bf6d-11eb-b944-e4e74940035b\\index.mp4', 'andy barbour smile', ',1080', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('ecc30c6b-bedd-11eb-9bd3-e4e74940035b', 'andy barbour', '2021-05-27 14:23:14', 15, 3,
        'content\\ecc30c6b-bedd-11eb-9bd3-e4e74940035b\\screen.jpg',
        'content\\ecc30c6b-bedd-11eb-9bd3-e4e74940035b\\index.mp4', 'andy barbour smile', ',1440,1080,720,480,360', 0,
        'cfaff592-b933-11eb-b073-e4e74940035b');
INSERT INTO video.video (id_video, title, uploaded, duration, status, thumbnail_url, url, description, quality, views,
                         owner_id)
VALUES ('fcd589a1-bc76-11eb-afc7-e4e74940035b', 'Black Mirror _ Season 4 Episode Titles _ Netflix',
        '2021-05-24 13:01:20', 50, 3, 'content\\fcd589a1-bc76-11eb-afc7-e4e74940035b\\screen.jpg',
        'content\\fcd589a1-bc76-11eb-afc7-e4e74940035b\\index.mp4', 'Black Mirror _ Season 4 Episode Titles _ Netflix',
        '1080,720,480,360', 0, 'cfaff592-b933-11eb-b073-e4e74940035b');