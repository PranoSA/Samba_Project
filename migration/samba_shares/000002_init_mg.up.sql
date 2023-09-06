ALTER TABLE Samba_Shares
ADD time_created TIMESTAMP WITH TIME ZONE;

ALTER TABLE Samba_Hosts 
ADD jointoken VARCHAR(128);

ALTER TABLE Samba_File_Systems
ADD time_created TIMESTAMP WITH TIME ZONE;

ALTER TABLE Samba_Spaces
ADD time_created TIMESTAMP WITH TIME ZONE;

ADD TABLE Samba_Invites (
    inviteid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    share_id uuid REFERENCES Samba_Shares(shareid),
    email owner,
    time_created TIMESTAMP WITH TIME ZONE,
    invite_code VARCHAR(64),
    hashed_invite bytea 
);

ADD TABLE Samba_Users (
    email VARCHAR(128) PRIMARY KEY ,
    share_id uuid REFERENCES Samba_Shares(shareid)
);

