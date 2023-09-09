#!/bin/bash

# sudo sh add_user_samba.sh "wadler-gmail.com" "password" "boba"

SAMBA_USER=$1
SAMBA_PASSWORD=$2

SHARE_ID=$3

SPACE_ID=$4


useradd -M -G $SPACE_ID $SAMBA_USER-$SHARE_ID

(echo $SAMBA_PASSWORD; echo $SAMBA_PASSWORD) | smbpasswd -a $SAMBA_USER-$SHARE_ID


#awk '{sub(/valid users = */,"valid users = $USER ")}1' samba_shell_test.txt 
#awk '{sub(/valid users = */,"valid users = poop ")}1' samba_shell_test.txt 

awk -v user=$SAMBA_USER-$SHARE_ID '{sub(/valid users = */,"valid users = "user" ")}1' /etc/samba/smb.conf.d/$SHARE_ID.conf>here.txt # >/etc/samba/smb.conf.d/$SHARE_ID.conf
cat here.txt > /etc/samba/smb.conf.d/$SHARE_ID.conf

#awk -v user="poop" '{sub(/valid users = */,"valid users = "user" ")}1' samba_shell_test.txt 

