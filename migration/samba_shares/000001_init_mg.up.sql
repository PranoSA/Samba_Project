CREATE TABLE IF NOT EXISTS Samba_Servers (
    serverid INTEGER PRIMARY KEY,
    lastip VARCHAR(128),
    hostname VARCHAR(128)
);

CREATE TABLE Samba_File_Systems (
    fsid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    server_id INTEGER NOT NULL REFERENCES Samba_Servers(serverid),
    device VARCHAR(128),
    mnt_point VARCHAR(255),
    capacity INTEGER NOT NULL
);


CREATE TABLE Samba_Spaces (
    spaceid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    fs_id uuid NOT NULL REFERENCES Samba_File_Systems(fsid),
    owner VARCHAR(128),
    alloc_size INTEGER NOT NULL
);

CREATE INDEX space_owner ON Samba_Spaces USING btree (
    owner    
);  

CREATE TABLE Samba_Shares (
    shareid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id uuid NOT NULL REFERENCES Samba_Spaces(spaceid),
    owner VARCHAR(128)
);

CREATE INDEX share_owner ON Samba_Shares USING btree (
    owner
);


 