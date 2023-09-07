#!/bin/bash
USER=$1
PASSWORD=$1

MOUNT_PATH=$1
SHARE_ID=$2
#SHARE_NAME=$3 #Name + ID Really  Add Later
OWNER=$3
OWNER_PASSWORD=$4
SPACE_ID=$5

if [ -z "$1"]
then 
    echo "No Mount Path Given"
    exit 1 
fi 


useradd -M -G $SPACE_ID $USER:$SHARE_ID

cat > /etc/samba/smb.conf.d/$SHARE_ID.conf << EOF
[$SHARE_ID]
    path = $MOUNT_PATH/$SHARE_ID
    browseable = no
    read only = no
    force create mode = 0777
    force directory mode = 0777
    create mask = 0777
    directory mask = 0777
EOF


(echo $PASSWORD; echo $PASSWORD) | smbpasswd -a $USER:$SHARE_ID


ls /etc/samba/smb.conf.d/* | sed -e 's/^/include = /' > /etc/samba/includes.conf

smbcontrol all reload-config



bob = ''' useradd -M -d /backups/$USER -s /usr/sbin/nologin -G sambashare $USER
mkdir /backups/$USER
chown $USER:sambashare /backups/$USER
chmod 2770 /backups/$USER
smbpasswd -a $USER
smbpasswd -e $USER

cat  /etc/samba/smb.conf.d/$USER.conf
[$USER]
    path = /$mount_pa/$USER
    browseable = no
    read only = no
    force create mode = 0660
    force directory mode = 2770
    valid users = $USER
EOF

ls /etc/samba/smb.conf.d/* | sed -e 's/^/include = /' > /etc/samba/includes.conf

smbcontrol all reload-config

'''