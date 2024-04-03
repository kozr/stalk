-- Insert mock users
INSERT INTO users (id, username, created_at, updated_at) VALUES ('1', 'user1', NOW(), NOW());
INSERT INTO users (id, username, created_at, updated_at) VALUES ('2', 'user2', NOW(), NOW());
INSERT INTO users (id, username, created_at, updated_at) VALUES ('3', 'user3', NOW(), NOW());
-- Insert mock followers
INSERT INTO followers (user_id, follower_id, created_at) VALUES ('1', '2', NOW());
INSERT INTO followers (user_id, follower_id, created_at) VALUES ('1', '3', NOW());
INSERT INTO followers (user_id, follower_id, created_at) VALUES ('2', '1', NOW());