# Operation


https://app.diagrams.net/#G1-ig-auPkGtmx-x2dDd2pt-YP4oedBp2J

# Option Summary Guide


## User Management / Authentication

Configure to store users in an LDAP directory and authenticate using a Bind operation or Kerberos, or
Configure to store users and authenticate using a Postgres Database with an md5 (hash on client), bcrypt, or Argon hashed passcode or \
\
Configure to store users an authenticate using a DynamoDB Database with the same above hashes or 
\
\
Configure to use an OIDC provider and Bearer token 
\

OIDC is generally reccommended for a public facing web app, LDAP is reccomended if you have an already existing LDAP database with customers and you want to provide this service, Postgres or DynamoDB is recomended if you want to manage your own datastore but want something that scales to the range of web users, DynamoDB if one wishes not to manage their own database.

## Session Management


Configure to use OIDC provider Bearer tokens for Authentication or \

Configure to use a Redis Database with sessions using cookies and CSRF tokens or
\
\
Configure to Use simple-auth header that re-authenticates using a username/passcode for each request


## Administrative Backing Store Options

This section is to store stateful information about samba services, such as available file systems mounted on the samba servers, which backing servers exist on the cluster and at what IP, where the available file systems with available storage are that can mount a new share space, and replication regions. 

The Database Tables (Or Entities) associated with this are 

### SambaServers
    Stores Information About Which Server IDs exist at which IPs, 

### Spaces
    Stores Information About Which Samba_Share Server a space resides, how much room is allocated to it, and owner
    
### Samba_Shares
    Stores Information about Who Can Access the Samba Share, Which Space it Resides In. 

### ISCSIServerFileSystems
    In The Case you are using an ISCSI target for backing samba server, you can run a client-server program on the samba server and ISCSI server respectively that enforces a level of consistency between them. 


### Options:
1. ETCD (For Company Level Deployments)
2. Postgres (For Web-Scale Level Deployments)
3. DynamoDB (If you wish not to manage your own database)


## Noun List 


### User : Globally Defined User


### Space : Allocated Space For a File System Identified by User.id + spaceid 


### Share : File System Share


### File-System: Tracks File Systems Mounted on the Samba Share Servers (Backend Only)


### ISCSI_FS_ID : Tracks ISCSI Targets and Backup Replicas


### Samba_Server : Tracks Samba Servers, Their IP Locations and Replicas


### ISCSI_SERVER: Trcks ISCSI Targets, Their IP, ID, and Replicas


# Deployment Options

## Session Methods

Define How To Authentiate and Authorize HTTP Requests to the server

### Basic Authentication

### Cookie Based Authentication

### Bearer Token Authentication


## Authentication Methods


### Allowing Signup with LDAP or Postgres or DynamoDB

To allow creation of user store using 


### OIDC / JWKS


## Server Management Storage


### Postgres


### DynamoDB (Serverless)


### ETCD (Server Management Only)


## User Storage 

### Postgres

### DynamoDB (Serverless)




