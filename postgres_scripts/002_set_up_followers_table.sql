-- Create the follower table
CREATE TABLE followers (
    user_id CHAR(64) NOT NULL,
    follower_id CHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (follower_id) REFERENCES users(id)
);