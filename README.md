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
List of Objects that will be referenced in the Code and Documentation 

### User : Globally Defined User
Globally Defined User to The Web Application


### Samba User
User Defined with smbpasswd on a samba password, these credentials aren't cross-referenced anywhere, 
the user name will be the ID of the samba Share + Login Name to the web application


### Space : 
Allocated Space For a File System Identified by User.id + spaceid 


### Share :
 File System Share with an owner who allocates it through Web API


### File-System: 
Tracks File Systems Mounted on the Samba Share Servers (Backend Only)


### ISCSI_FS_ID : 
Tracks ISCSI Targets and Backup Replicas


### Samba_Server :
 Tracks Samba Servers, Their IP Locations and Replicas


### ISCSI_SERVER: (Might Remove)
Tracks ISCSI Targets, Their IP, ID, and Replicas


# Deployment Options



<br>

## Login/Signup Methods

Configuration Can Be Setup On Where to Store Usernames and their respective passwords (as well as the challenge hash to be done)
to authenticate users and return credentials to the User to map to sessions.

When using OIDC, this is done by the OIDC Provider where the client will have to use OIDC Configuration (Auth URL, ClientID or 
ClientID + ClientSecret, Redirect URIs, Token Endpoints) to receive these credentials.



### Allowing Signup with Postgres or DynamoDB

To allow creation of user store using a user table that stores a salted hash challenge row to authenticate against.
This allows Signup

### LDAP
Externally Managed LDAP search with a search base DN for users, allows authentication with attempted bind requests to the server.

### OIDC / JWKS
Externally Managed 0AUTH2 / OIDC Provider
<br>

## Session Methods

Configuration to be setup to use different Modes of Authentication for HTTP Requests to the Web-Server Application
These will be used to assume roles within the application.

The Session methods involve either Authentication Headers in the Request or Authentication Cookies with a CSRF Token 
mapped against a Redis Database. 
<br>

### Basic Authentication

Authentication Header "Basic" +" " + base64(username):base64(string)
<br>

### Cookie Based Authentication

Cookie with Name Authentication and HTTP-Only Access 
Needs to Be configured with Domain Attribute 

<br>
### Bearer Token Authentication

For OIDC Tokens from the SSO Keycloak or other OIDC provider instance

Authentication Header "Bearer" + " " + id_token
 
 <br>

## Server Management Storage

This Will Be Database Tables To Be used to store information about 
Samba Hosts
Samba File Systems
Samba File Spaces
Samba Shares / Groups / Invites
and information about accesibility per user

### Postgres

A sample migration to setup the necessary tables for the Postgres Configuration are in this repository under the migration folder.
<br>

### DynamoDB (Serverless)

This application will assume 3 DynamoDB Table Names 
1. Samba_Hosts -> Stores Information about Hosts and Their File Systems and what Spaces Have been alocated to them
2. Samba Spaces -> Stores Information about Allocated Spaces and What Samba Shares Have Been Allocated To Them
3. Samba Shares -> Stores Information About Samba Shares, What space they live on, and invites and group members to those samba shares.

The recommended indexing for these tables are as follows:
Other index names will not be searched against

#### Samba_Hosts
Primary Key Index (Only Partition Key) on column hostid and no Global Secondary Index/

#### Samba_Spaces
Primary Key Index (Only Partition Key) on column spaceid and Global Secondary Index with Only Partition Key fs_id

#### Samba_Shares
Primary Key Index (Only Partition Key) on column shareid and Global Secondary Index (Partition Key Only) on space_id


###  ETCD (Server Management Only) -> STILL NEEDS TO BE IMPLEMENTED


<br>

## User Storage (This Feature Doesn't Exist For Now )

### Postgres

### DynamoDB (Serverless)




