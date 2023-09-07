#!/bin/bash

USER = $1
PASSWORD = $2

SHARE_ID = $3


useradd -M -G $SPACE_ID $USER:$SHARE_ID

(echo $PASSWORD; echo $PASSWORD) | smbpasswd -a $USER:$SHARE_ID


#awk '{sub(/valid users = */,"valid users = $USER ")}1' samba_shell_test.txt 
#awk '{sub(/valid users = */,"valid users = poop ")}1' samba_shell_test.txt 

awk -v user=$USER '{sub(/valid users = */,"valid users = "user" ")}1' /etc/samba/smb.conf.d/$SHARE_ID

#awk -v user="poop" '{sub(/valid users = */,"valid users = "user" ")}1' samba_shell_test.txt 

