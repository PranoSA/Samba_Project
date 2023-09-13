CREATE TABLE StreamLinks (
    file_name VARCHAR(256),
    share_id uuid REFERENCES Samba_Shares(shareid),
    PRIMARY KEY(share_id, file_name),
    email VARCHAR(128)
);

CREATE TABLE CompressLinks (
    compressid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(128)
);

/**
    Change This Later...
*/