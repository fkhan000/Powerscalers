CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR NOT NULL,
    profile_pic VARCHAR,
    cash FLOAT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Communities (
    community_id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    picture VARCHAR,
    num_followers INTEGER DEFAULT 0,
    rate_limit INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Moderators (
    user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
    community_id INTEGER REFERENCES Communities(community_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, community_id)
);

CREATE TABLE Wagers (
    wager_id SERIAL PRIMARY KEY,
    community_id INTEGER REFERENCES Communities(community_id) ON DELETE SET NULL,
    owner_id INTEGER REFERENCES Users(user_id) ON DELETE SET NULL,
    title VARCHAR NOT NULL,
    description TEXT,
    left VARCHAR NOT NULL,
    right VARCHAR NOT NULL,
    decision VARCHAR DEFAULT "",
    explanation TEXT DEFAULT "",
    net_likes INTEGER DEFAULT 0,
    total_comments INTEGER DEFAULT 0,
    expiration_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Gambles (
    user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
    wager_id INTEGER REFERENCES Wagers(wager_id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    position VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, wager_id)
);

CREATE TABLE Comments (
    comment_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
    wager_id INTEGER REFERENCES Wagers(wager_id) ON DELETE CASCADE,
    parent_comment_id INTEGER REFERENCES Comments(comment_id) DEFAULT -1,
    description TEXT NOT NULL,
    net_likes INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Wager_Likes (
    user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
    wager_id INTEGER REFERENCES Wagers(wager_id) ON DELETE CASCADE,
    value INTEGER NOT NULL CHECK (value IN (-1, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, wager_id)
);

CREATE TABLE Comment_Likes (
    user_id INTEGER REFERENCES Users(user_id) ON DELETE CASCADE,
    comment_id INTEGER REFERENCES Comments(comment_id) ON DELETE CASCADE,
    value INTEGER NOT NULL CHECK (value IN (-1, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, comment_id)
);