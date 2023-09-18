CREATE TABLE StreamLinks (
    file_name VARCHAR(256),
    share_id uuid REFERENCES Samba_Shares(shareid),
    PRIMARY KEY(share_id, file_name),
    email VARCHAR(128)
);

CREATE TABLE CompressLinks (
    share_id uuid REFERENCES Samba_Shares(shareid),
    time_backed TIMESTAMP WITH TIME ZONE DEFAULT now(),
    creator VARCHAR(128),
    PRIMARY KEY(share_id, time_backed)
);

/**
    compressid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    Change This Later...
*/