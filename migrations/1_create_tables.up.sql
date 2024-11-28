CREATE TABLE IF NOT EXISTS groups (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(255) UNIQUE NOT NULL,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_groups_name ON groups(name);

CREATE TABLE IF NOT EXISTS songs (
                                     id SERIAL PRIMARY KEY,
                                     group_id INT NOT NULL,
                                     song VARCHAR(255) NOT NULL,
                                     release_date DATE,
                                     lyrics TEXT NULL ,
                                     link VARCHAR(255),
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE INDEX idx_songs_group_id ON songs(group_id);
CREATE INDEX idx_songs_song ON songs(song);

INSERT INTO groups (name)
VALUES
    ('Imagine Dragons'),
    ('Arctic Monkeys')
ON CONFLICT (name) DO NOTHING;


INSERT INTO songs (group_id, song, release_date, lyrics, link)
VALUES
    (
        (SELECT id FROM groups WHERE name = 'Imagine Dragons'),
        'Believer',
        '2017-02-01',
        'Verse 1\n\nVerse 2\n\nVerse 3',
        'https://example.com/believer'
    ),
    (
        (SELECT id FROM groups WHERE name = 'Imagine Dragons'),
        'Thunder',
        '2017-04-27',
        'Verse1 \n\nverse2 \n\nverse3',
        'https://example.com/thunder'
    ),
    (
        (SELECT id FROM groups WHERE name = 'Arctic Monkeys'),
        'Do I Wanna Know?',
        '2013-06-18',
        'Verse 1: Have you got color in your cheeks?\n\nVerse 2: Do you ever get that fear that you can’t shift\n\nVerse 3: Crawling back to you...',
        'https://example.com/do-i-wanna-know'
    )
    ON CONFLICT DO NOTHING;