package proto_samba_management

/**
 * To Take Backup Snapshot of Content in Share...
 */
type BackupRequest struct {
	Share_id string
	Filename string
}

type BackupResponse struct {
}

var (
	/**
	 * Exchange To Use
	 */
	Exchange_Backup = "samba"
	//Backup_Exchange = "samba"

	//Name of Queue
	//Backup_Listening_queue = "samba_listening"
	Queue_Listening_Backup = "samba_listening"

	/**
	 * Routing Keys To Queues
	 */
	KeyCompressRequest = "share"

	//SambaCompressRequest = "share"

	Bucket_Backup = "backup"
)
