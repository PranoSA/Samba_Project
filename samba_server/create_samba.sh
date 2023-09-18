#!/bin/bash

#Test sudo sh create_samba.sh /mount/samba_server/1 "boba" "pcadler-gmail.com" "password" 1 

MOUNT_PATH=$1
SHARE_ID=$2
#SHARE_NAME=$3 #Name + ID Really  Add Later
OWNER=$3
OWNER_PASSWORD=$4
SPACE_ID=$5

if [ -z "$1" ]; then
    echo "No Mount Path Given"
    #exit 1
fi 

#addgroup samba-$SPACE_ID

status=$?

if [ $status -ne 0 ]; then
    echo "Failed To Add Group"
   # exit 1
fi

useradd -M -G $SPACE_ID $OWNER-$SHARE_ID

status=$?

if [ $status -ne 0 ]; then
    echo "Failed To Add User"
    #exit 1
fi

cat >/etc/samba/smb.conf.d/$SHARE_ID.conf <<EOF
[$SHARE_ID]
    path = $MOUNT_PATH/$SHARE_ID
    browseable = yes
    read only = no
    force create mode = 0777
    force directory mode = 0777
    create mask = 0777
    directory mask = 0777
    valid users = $OWNER-$SHARE_ID
EOF

status=$?

if [ $status -ne 0 ]; then
    echo "Failed To Write to Samba Config "
    #exit 1
fi

(
    echo $OWNER_PASSWORD
    echo $OWNER_PASSWORD
) | smbpasswd -a $OWNER-$SHARE_ID

status=$?

if [ $status -ne 0 ]; then
    echo "Failed To Set Samba password"
    #exit 1
fi

ls /etc/samba/smb.conf.d/* | sed -e 's/^/include = /' >/etc/samba/includes.conf

status=$?

if [ $status -ne 0 ]; then
    echo "Failed To Add TO Include File"
    #exit 1
fi

smbcontrol all reload-config

status=$?

if [ $status -ne 0 ]; then
    echo "Failed To Reload Samba"
    #exit 1
fi
