CREATE INDEX IN_video_comments_user_id
    ON video_comments (user_id);
CREATE INDEX IN_video_comments_video_id
    ON video_comments (video_id);
CREATE INDEX IN_video_comments_parent_id
    ON video_comments (parent_id);
CREATE INDEX IN_video_comments_created
    ON video_comments (created);