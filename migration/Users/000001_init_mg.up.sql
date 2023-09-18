CREATE TABLE IF NOT EXISTS Samba_Share_Users (
    user_id VARCHAR(128),
    share_id uuid REFERENCES Samba_Shares(shareid),
    privilege INTEGER,
    PRIMARY KEY(user_id, share_id)
);