CREATE TABLE ROLES
(
    ID        SERIAL PRIMARY KEY,
    ROLE_NAME VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE PROJECTS
(
    ID          SERIAL PRIMARY KEY,
    NAME        VARCHAR(50) UNIQUE NOT NULL,
    DESCRIPTION TEXT               NOT NULL
);

CREATE TABLE USERS
(
    ID       SERIAL PRIMARY KEY,
    USERNAME VARCHAR(50) UNIQUE NOT NULL,
    EMAIL    VARCHAR(50) UNIQUE NOT NULL,
    ROLE     INT                NOT NULL,
    FOREIGN KEY (ROLE) REFERENCES ROLES (ID)
);

CREATE TABLE POSTS
(
    ID      SERIAL PRIMARY KEY,
    TITLE   TEXT UNIQUE NOT NULL,
    CONTENT TEXT        NOT NULL,
    SHORT_CONTENT TEXT NOT NULL,
    OWNER   INT         NOT NULL,
    TIMESTAMP TIMESTAMP NOT NULL,
    FOREIGN KEY (OWNER) REFERENCES USERS (ID)
);

CREATE TABLE COMMENTS
(
    ID      SERIAL PRIMARY KEY,
    COMMENT TEXT NOT NULL,
    OWNER   INT  NOT NULL,
    POST    INT  NOT NULL,
    TIMESTAMP TIMESTAMP NOT NULL,
    FOREIGN KEY (POST) REFERENCES POSTS (ID),
    FOREIGN KEY (OWNER) REFERENCES USERS (ID)
);

CREATE TABLE OAUTH_CREDENTIALS
(
    ID            SERIAL PRIMARY KEY,
    USER_ID       INT          NOT NULL,
    PROVIDER      VARCHAR(50)  NOT NULL,
    PROVIDER_ID   VARCHAR(100) NOT NULL,
    ACCESS_TOKEN  TEXT         NOT NULL,
    REFRESH_TOKEN TEXT,
    EXPIRES_AT    TIMESTAMP,
    FOREIGN KEY (USER_ID) REFERENCES USERS (ID),
    UNIQUE (PROVIDER, PROVIDER_ID)
);

CREATE TABLE VOTES(
                      ID SERIAL PRIMARY KEY,
                      USER_ID INT NOT NULL,
                      POST_ID INT NOT NULL,
                      TIMESTAMP TIMESTAMP NOT NULL,
                      VALUE INT NOT NULL,
                      FOREIGN KEY (USER_ID) REFERENCES USERS (ID),
                      FOREIGN KEY (POST_ID) REFERENCES POSTS (ID)
);


-- Indexes for performance optimization
CREATE INDEX idx_users_role ON USERS (ROLE);
CREATE INDEX idx_posts_owner ON POSTS (OWNER);
CREATE INDEX idx_comments_post ON COMMENTS (POST);
CREATE INDEX idx_comments_owner ON COMMENTS (OWNER);
CREATE INDEX idx_votes_user ON VOTES (USER_ID);
CREATE INDEX idx_votes_post ON VOTES (POST_ID);
-- Indexes on TIMESTAMP columns for efficient time-based queries
CREATE INDEX idx_posts_timestamp ON POSTS (TIMESTAMP);
CREATE INDEX idx_comments_timestamp ON COMMENTS (TIMESTAMP);
CREATE INDEX idx_votes_timestamp ON VOTES (TIMESTAMP);