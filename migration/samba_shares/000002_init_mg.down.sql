ALTER TABLE Samba_Shares
DROP COLUMN time_created;


ALTER TABLE Samba_Hosts 
DROP COLUMN jointoken; 


ALTER TABLE Samba_File_Systems
DROP COLUMN time_created; 


ALTER TABLE Samba_Spaces
DROP COLUMN time_created; 