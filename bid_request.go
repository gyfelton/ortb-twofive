package ad_mediation

// Reference: https://github.com/marcsantiago/ortb-twofive/blob/master/bid_request.go
// This file is modified according to what Nimbus expects and what Wattpad really need
// Header required in RTB requests
const (
	OpenRTBVersionHeader  = "x-openrtb-version"
	OpenRTBVersionValue = "2.5"

	EncodingHeader = "Content-Encoding"
	EncodingValue  = "gzip"
)

// Root level object
// Request openRTB 2.5 spec
type OpenRTBRequest struct {
	Test    int        `json:"test,omitempty"    valid:"range(0|1),optional"`
	Imp     []Imp      `json:"imp"               valid:"required"`
	App     App        `json:"app"               valid:"required"`
	Device  Device     `json:"device"            valid:"required"`
	Format  Format     `json:"format"            valid:"required"`
	User    User       `json:"user"              valid:"required"`
	Tmax    int        `json:"tmax,omitempty"    valid:"-"`
	Ext     RequestExt `json:"ext,omitempty"     valid:"required"`
}

// RequestExt used to communicate the publishers api key
type RequestExt struct {
	APIKey 		string `json:"api_key" valid:"uuidv4,required"`
	SessionID 	string `json:"session_id" valid:required"`
}

// Imp describes an ad placement or impression being auctioned. A single bid request can include
// multiple Imp objects. For dual auction to happen, you can pass both Banner and Video
// but one of them must be present
type Imp struct {
	Banner               *Banner  `json:"banner,omitempty"               valid:"optional"`
	Video                *Video   `json:"video,omitempty"                valid:"optional"`
	Instl                int      `json:"instl,omitempty"                valid:"range(0|1)"` // 0 = not interstitial, 1 = interstitial
	BidFloor             float64  `json:"bidfloor,omitempty"                       valid:"optional"` // Needed if we want to changed the auction floor
}

// Banner represents the most general type of impression. Although the term “banner” may have very
// specific meaning in other contexts, here it can be many things including a simple static image, an
// expandable ad unit, or even in-banner video (refer to the Video object in Section 3.2.4 for the more
// generalized and full featured video ad units). An array of Banner objects can also appear within the
// Video to describe optional companion ads defined in the VAST specification.
type Banner struct {
	Format []Format   `json:"format,omitempty" valid:"optional"`
	Pos    int        `json:"pos,omitempty"    valid:"range(0|7),optional"`            // 0,1,2,3,4,5,6,7 -> Unknown,Above the Fold,DEPRECATED - May or may not be initially visible depending on screen size/resolution.,Below the Fold,Header,Footer,Sidebar,Full Screen
	API    []int      `json:"api,omitempty"    valid:"inintarr(1|2|3|4|5|6),optional"` // 3,5,6 -> mraid1, 2, and 3
}

// Video object represents an in-stream video impression. Many of the fields are non-essential for minimally
// viable transactions, but are included to offer fine control when needed. Video in OpenRTB generally
// assumes compliance with the VAST standard. As such, the notion of companion ads is supported by
// optionally including an array of Banner objects (refer to the Banner object in Section 3.2.3) that define
// these companion ads.
type Video struct {
	Mimes          []string  `json:"mimes,omitempty"          valid:"required"` // Should only be ["video/mp4"] for now since Apple restrict video to be mp4 only on iOS devices
	Minduration    int       `json:"minduration"              valid:"-"`
	Maxduration    int       `json:"maxduration,omitempty"    valid:"-"`
	Protocols      []int     `json:"protocols,omitempty"      valid:"inintarr(2|3|5|6),optional"` // 1,2,3,4,5,6,7,8,9,10 -> VAST 1.0,VAST 2.0,VAST 3.0,VAST 1.0 Wrapper,VAST 2.0 Wrapper,VAST 3.0 Wrapper,VAST 4.0,VAST 4.0 Wrapper,DAAST 1.0,DAAST 1.0 Wrapper
	W              int       `json:"w,omitempty"              valid:"required"`
	H              int       `json:"h,omitempty"              valid:"required"`
	Linearity      int       `json:"linearity,omitempty"      valid:"range(1|2),optional"`            // 1,2 -> linear, non linear
	Skip           int       `json:"skip"                     valid:"range(0|1),optional"`            // 0 no 1 yes
	SkipMin        int       `json:"skipmin"                  valid:"optional"`
	SkipAfter      int       `json:"skipafter"                valid:"optional"`
	Playbackmethod []int     `json:"playbackmethod,omitempty" valid:"inintarr(1|2|3|4|5|6),optional"` // 1,2,3,4,5,6 - > Initiates on Page Load with Sound On, Initiates on Page Load with Sound Off by Default, Initiates on Click with Sound On, Initiates on Mouse-Over with Sound On, Initiates on Entering Viewport with Sound On, Initiates on Entering Viewport with Sound Off by Default
	MinBitRate     int       `json:"minbitrate,omitempty"     valid:"-"`
	MaxBitRate     int       `json:"maxbitrate,omitempty"     valid:"-"`
	Pos            int       `json:"pos,omitempty"            valid:"range(0|7),optional"`            // 0,1,2,3,4,5,6,7 -> Unknown,Above the Fold,DEPRECATED - May or may not be initially visible depending on screen size/resolution.,Below the Fold,Header,Footer,Sidebar,Full Screen
}

// Format object represents an allowed size (i.e., height and width combination) for a banner impression.
// These are typically used in an array for an impression where multiple sizes are permitted.
type Format struct {
	W      int        `json:"w,omitempty"      valid:"required"`
	H      int        `json:"h,omitempty"      valid:"required"`
}

// App object should be included if the ad supported content is a non-browser application (typically in
// mobile) as opposed to a website. A bid request must not contain both an App and a Site object. At a
// minimum, it is useful to provide an App ID or bundle, but this is not strictly required
type App struct {
	Name          string    `json:"name"                 valid:"required"`
	Bundle        string    `json:"bundle"               valid:"required"`
	Domain        string    `json:"domain"               valid:"required"`
	StoreURL      string    `json:"storeurl"             valid:"required"`
	Cat           []string  `json:"cat,omitempty"        valid:"required"`
	Ver           string    `json:"ver,omitempty"        valid:"-"`
	Paid          int       `json:"paid"                 valid:"range(0|1),optional"` // free 0 paid 1
	Publisher     Publisher `json:"publisher"            valid:"required"`
}

// Publisher object describes the publisher of the media in which the ad will be displayed. The publisher is
// typically the seller in an OpenRTB transaction.
type Publisher struct {
	Name   string        `json:"name"          valid:"required"`
	Cat    []string      `json:"cat,omitempty" valid:"-"`
	Domain string        `json:"domain"        valid:"required"`
	Ext    PublisherExt `json:"ext,omitempty" valid:"-"`
}

// PublisherExt ...
type PublisherExt struct{
	FacebookAppId string `json:"facebook_app_id,omitempty" valid:"-"`
}

// Device object provides information pertaining to the device through which the user is interacting. Device
// information includes its hardware, platform, location, and carrier data. The device can refer to a mobile
// handset, a desktop computer, set top box, or other digital device.
type Device struct {
	Ua             string     `json:"ua"                        valid:"required"`
	Lmt            int        `json:"lmt"                       valid:"range(0|1),optional"` // 0 = tracking is unrestricted, 1 = tracking must be limited by commericial guidelines
	IP             string     `json:"ip"                        valid:"ipv4,required"`
	Make           string     `json:"make,omitempty"            valid:"in(Apple|Android|apple|android),required"`
	Model          string     `json:"model,omitempty"           valid:"-"`
	OS             string     `json:"os,omitempty"              valid:"-"`
	OSV            string     `json:"osv,omitempty"             valid:"-"`
	ConnectionType int        `json:"connection_type,omitempty" valid:"-"`
	Ifa            string     `json:"ifa,omitempty"             valid:"required"` // passing of invalid value is enough indication that user did not want to be tracked
}

// User object contains information known or derived about the human user of the device (i.e., the
// audience for advertising). The user id is an exchange artifact and may be subject to rotation or other
// privacy policies. However, this user ID must be stable long enough to serve reasonably as the basis for
// frequency capping and retargeting.
// Since we do not have any consent machanism, we do not have ext object for user
type User struct {
	Age        int      `json:"age,omitempty"` // outside of spec, but i'm allowing for it
	BuyerUID   string   `json:"buyeruid,omitempty"    valid:"-"`
	YOB        int      `json:"yob,omitempty"         valid:"-"`
	Gender     string   `json:"gender,omitempty"      valid:"in(M|F|O|Male|male|Female|female),optional"`
}
