package proto_samba_management

/**
 *
 * Message To Publish To Exchange With The Dash Published Routing Key
 * Tells the web_manager server to transform using dash script
 */
type DashMessage struct {
	Share_id    string
	Filename    string //Where In Bucket Resides...
	Resolutions []struct {
		Width    int
		Height   int
		Bit_Rate int
	}
}

/**
 * Do We Even Need This ????
 */

type DashRequest struct {
	Share_id    string
	Filename    string
	Resolutions []struct {
		Width    int
		Height   int
		Bit_Rate int
	}
}

type DashFinished struct {
	Share_id string
	URL      string
	Filename string
}

var (
	Correct        = 0
	Full           = 1
	Upstream_error = 2
)

var (
	/**
	 * Exchange For Dash Requests (maybe Samba Would Be Better)
	 */
	Dash_exchange = "samba" //Name Of Exchange Used or Dash Messages
	/**
	 * Routing Keys
	 */
	Dash_Request   = "dash request"
	Dash_published = "http_dash_new"      //Used To Signify A New Request for a Dash Transformation Has Been Made
	Dash_complete  = "http_dash_complete" //Used To Signify a file has been Dashed
	/**
	 * Name For QUeues That BindTO Exchange AT Routing Keys
	 *
	 */
	Dash_Queue_Requests = "dash queue requests"

	/**
	 * Name of Buckets ForPublsihing and storing temporarily
	 */

	PublishBucket = "frontend"

	TempBucket = "backend"
)
