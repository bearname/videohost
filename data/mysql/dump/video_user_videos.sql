create table user_videos
(
    key_user varchar(255) not null,
    id_video varchar(255) not null,
    constraint user_videos_ibfk_1
        foreign key (key_user) references users (key_user)
            on delete cascade,
    constraint user_videos_ibfk_2
        foreign key (id_video) references video (id_video)
            on delete cascade
)
    charset = utf8;

INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'b2b8ba78-b86f-11eb-a0d6-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('cfaff592-b933-11eb-b073-e4e74940035b', 'f9648da9-b865-11eb-a0d6-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'e016b2cf-b5d6-11eb-b729-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('cfaff592-b933-11eb-b073-e4e74940035b', 'cf561af0-b5f0-11eb-a7d7-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'c5b34027-b791-11eb-a175-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'b591502a-b91c-11eb-a2fc-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', '57655d3e-b92d-11eb-989d-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'f3c685ea-b942-11eb-a79d-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'c6c04ca0-b943-11eb-a79d-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'cdc1a0b3-b943-11eb-a79d-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'd4065d39-b943-11eb-a79d-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'de1fa90f-b946-11eb-a79d-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', '34bc4791-b947-11eb-be57-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', 'ac7545ec-b947-11eb-be57-e4e74940035b');
INSERT INTO video.user_videos (key_user, id_video) VALUES ('4bbe7653-b933-11eb-b073-e4e74940035b', '4bebbff8-b948-11eb-be57-e4e74940035b');