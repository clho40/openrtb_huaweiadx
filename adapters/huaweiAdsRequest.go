package adapters

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prebid/openrtb/v17/native1"
	nativeRequests "github.com/prebid/openrtb/v17/native1/request"
	openrtb2 "github.com/prebid/openrtb/v17/openrtb2"
	constants "main.go/utils"
)

// creative type
const (
	text                   int32 = 1
	bigPicture             int32 = 2
	bigPicture2            int32 = 3
	gif                    int32 = 4
	videoText              int32 = 6
	smallPicture           int32 = 7
	threeSmallPicturesText int32 = 8
	video                  int32 = 9
	iconText               int32 = 10
	videoWithPicturesText  int32 = 11
)

// interaction type
const (
	appPromotion int32 = 3
)

// ads type
const (
	banner       int32 = 8
	native       int32 = 3
	roll         int32 = 60
	interstitial int32 = 12
	rewarded     int32 = 7
	splash       int32 = 1
	magazinelock int32 = 2
	audio        int32 = 17
)

const huaweiAdxApiVersion = "3.4"
const defaultCountryName = "ZA"
const defaultUnknownNetworkType = 0
const timeFormat = "2006-01-02 15:04:05.000"
const defaultTimeZone = "+0200"
const defaultModelName = "HUAWEI"
const defaultEndpoint = "https://adx-dre.op.hicloud.com/ppsadx/getResult"
const chineseSiteEndPoint = "https://acd.op.hicloud.com/ppsadx/getResult"
const europeanSiteEndPoint = "https://adx-dre.op.hicloud.com/ppsadx/getResult"
const asianSiteEndPoint = "https://adx-dra.op.hicloud.com/ppsadx/getResult"
const russianSiteEndPoint = "https://adx-drru.op.hicloud.com/ppsadx/getResult"

type HuaweiAdsRequest struct {
	Version           string     `json:"version"`
	Multislot         []adslot30 `json:"multislot"`
	App               app        `json:"app"`
	Device            device     `json:"device"`
	Network           network    `json:"network,omitempty"`
	Regs              regs       `json:"regs,omitempty"`
	Geo               geo        `json:"geo,omitempty"`
	Consent           string     `json:"consent,omitempty"`
	ClientAdRequestId string     `json:"clientAdRequestId,omitempty"`
}

type adslot30 struct {
	Slotid                   string   `json:"slotid"`
	Adtype                   int32    `json:"adtype"`
	Test                     int32    `json:"test"`
	TotalDuration            int32    `json:"totalDuration,omitempty"`
	Orientation              int32    `json:"orientation,omitempty"`
	W                        int64    `json:"w,omitempty"`
	H                        int64    `json:"h,omitempty"`
	Format                   []format `json:"format,omitempty"`
	DetailedCreativeTypeList []string `json:"detailedCreativeTypeList,omitempty"`
}

type format struct {
	W int64 `json:"w,omitempty"`
	H int64 `json:"h,omitempty"`
}

type app struct {
	Version string `json:"version,omitempty"`
	Name    string `json:"name,omitempty"`
	Pkgname string `json:"pkgname"`
	Lang    string `json:"lang,omitempty"`
	Country string `json:"country,omitempty"`
}

type device struct {
	Type                int32   `json:"type,omitempty"`
	Useragent           string  `json:"useragent,omitempty"`
	Os                  string  `json:"os,omitempty"`
	Version             string  `json:"version,omitempty"`
	Maker               string  `json:"maker,omitempty"`
	Model               string  `json:"model,omitempty"`
	Width               int32   `json:"width,omitempty"`
	Height              int32   `json:"height,omitempty"`
	Language            string  `json:"language,omitempty"`
	BuildVersion        string  `json:"buildVersion,omitempty"`
	Dpi                 int32   `json:"dpi,omitempty"`
	Pxratio             float32 `json:"pxratio,omitempty"`
	Imei                string  `json:"imei,omitempty"`
	Oaid                string  `json:"oaid,omitempty"`
	IsTrackingEnabled   string  `json:"isTrackingEnabled,omitempty"`
	EmuiVer             string  `json:"emuiVer,omitempty"`
	LocaleCountry       string  `json:"localeCountry"`
	BelongCountry       string  `json:"belongCountry"`
	GaidTrackingEnabled string  `json:"gaidTrackingEnabled,omitempty"`
	Gaid                string  `json:"gaid,omitempty"`
	ClientTime          string  `json:"clientTime"`
	Ip                  string  `json:"ip,omitempty"`
}

type network struct {
	Type     int32      `json:"type"`
	Carrier  int32      `json:"carrier,omitempty"`
	CellInfo []cellInfo `json:"cellInfo,omitempty"`
}

type regs struct {
	Coppa int32 `json:"coppa,omitempty"`
}

type geo struct {
	Lon      float32 `json:"lon,omitempty"`
	Lat      float32 `json:"lat,omitempty"`
	Accuracy int32   `json:"accuracy,omitempty"`
	Lastfix  int32   `json:"lastfix,omitempty"`
}

type cellInfo struct {
	Mcc string `json:"mcc,omitempty"`
	Mnc string `json:"mnc,omitempty"`
}

type huaweiAdsResponse struct {
	Retcode int32  `json:"retcode"`
	Reason  string `json:"reason"`
	Multiad []ad30 `json:"multiad"`
}

type ad30 struct {
	AdType    int32     `json:"adtype"`
	Slotid    string    `json:"slotid"`
	Retcode30 int32     `json:"retcode30"`
	Content   []content `json:"content"`
}

type content struct {
	Contentid       string    `json:"contentid"`
	Interactiontype int32     `json:"interactiontype"`
	Creativetype    int32     `json:"creativetype"`
	MetaData        metaData  `json:"metaData"`
	Monitor         []monitor `json:"monitor"`
	Cur             string    `json:"cur"`
	Price           float64   `json:"price"`
}

type metaData struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ImageInfo   []imageInfo `json:"imageInfo"`
	Icon        []icon      `json:"icon"`
	ClickUrl    string      `json:"clickUrl"`
	Intent      string      `json:"intent"`
	VideoInfo   videoInfo   `json:"videoInfo"`
	ApkInfo     apkInfo     `json:"apkInfo"`
	Duration    int64       `json:"duration"`
	MediaFile   mediaFile   `json:"mediaFile"`
}

type imageInfo struct {
	Url       string `json:"url"`
	Height    int64  `json:"height"`
	FileSize  int64  `json:"fileSize"`
	Sha256    string `json:"sha256"`
	ImageType string `json:"imageType"`
	Width     int64  `json:"width"`
}

type icon struct {
	Url       string `json:"url"`
	Height    int64  `json:"height"`
	FileSize  int64  `json:"fileSize"`
	Sha256    string `json:"sha256"`
	ImageType string `json:"imageType"`
	Width     int64  `json:"width"`
}

type videoInfo struct {
	VideoDownloadUrl string  `json:"videoDownloadUrl"`
	VideoDuration    int32   `json:"videoDuration"`
	VideoFileSize    int32   `json:"videoFileSize"`
	Sha256           string  `json:"sha256"`
	VideoRatio       float32 `json:"videoRatio"`
	Width            int32   `json:"width"`
	Height           int32   `json:"height"`
}

type apkInfo struct {
	Url         string `json:"url"`
	FileSize    int64  `json:"fileSize"`
	Sha256      string `json:"sha256"`
	PackageName string `json:"packageName"`
	SecondUrl   string `json:"secondUrl"`
	AppName     string `json:"appName"`
	VersionName string `json:"versionName"`
	AppDesc     string `json:"appDesc"`
	AppIcon     string `json:"appIcon"`
}

type mediaFile struct {
	Mime     string `json:"mime"`
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	FileSize int64  `json:"fileSize"`
	Url      string `json:"url"`
	Sha256   string `json:"sha256"`
}

type monitor struct {
	EventType string   `json:"eventType"`
	Url       []string `json:"url"`
}

type adapter struct {
	endpoint  string
	extraInfo ExtraInfo
}

type ExtraInfo struct {
	PkgNameConvert              []pkgNameConvert `json:"pkgNameConvert,omitempty"`
	CloseSiteSelectionByCountry string           `json:"closeSiteSelectionByCountry,omitempty"`
}

type pkgNameConvert struct {
	ConvertedPkgName           string   `json:"convertedPkgName,omitempty"`
	UnconvertedPkgNames        []string `json:"unconvertedPkgNames,omitempty"`
	UnconvertedPkgNameKeyWords []string `json:"unconvertedPkgNameKeyWords,omitempty"`
	UnconvertedPkgNamePrefixs  []string `json:"unconvertedPkgNamePrefixs,omitempty"`
	ExceptionPkgNames          []string `json:"exceptionPkgNames,omitempty"`
}

type PublishersCredential struct {
	SlotId              string `json:"slotid"`
	Adtype              string `json:"adtype"`
	PublisherId         string `json:"publisherid"`
	SignKey             string `json:"signkey"`
	KeyId               string `json:"keyid"`
	IsTestAuthorization string `json:"isTestAuthorization,omitempty"`
}

type ExtUserDataHuaweiAds struct {
	Data ExtUserDataDeviceIdHuaweiAds `json:"data,omitempty"`
}

type ExtUserDataDeviceIdHuaweiAds struct {
	Imei       []string `json:"imei,omitempty"`
	Oaid       []string `json:"oaid,omitempty"`
	Gaid       []string `json:"gaid,omitempty"`
	ClientTime []string `json:"clientTime,omitempty"`
}

type ExtImpBidder struct {
	// Prebid *openrtb_ext.ExtImpPrebid `json:"prebid"`
	Bidder json.RawMessage `json:"bidder"`
	// AuctionEnvironment openrtb_ext.AuctionEnvironmentType `json:"ae,omitempty"`
}

type ExtUser struct {
	// Consent is a GDPR consent string. See "Advised Extensions" of
	// https://iabtechlab.com/wp-content/uploads/2018/02/OpenRTB_Advisory_GDPR_2018-02.pdf
	Consent                          string                         `json:"consent,omitempty"`
	ConsentedProvidersSettings       *ConsentedProvidersSettingsIn  `json:"ConsentedProvidersSettings,omitempty"`
	ConsentedProvidersSettingsParsed *ConsentedProvidersSettingsOut `json:"consented_providers_settings,omitempty"`
	Eids                             []openrtb2.EID                 `json:"eids,omitempty"`
}

type ConsentedProvidersSettingsIn struct {
	ConsentedProvidersString string `json:"consented_providers,omitempty"`
}

type ConsentedProvidersSettingsOut struct {
	ConsentedProvidersList []int `json:"consented_providers,omitempty"`
}

type empty struct{}

type RequestData struct {
	Method  string
	Uri     string
	Body    []byte
	Headers http.Header
}

func MakeRequest(openRTBRequest *openrtb2.BidRequest) (*HuaweiAdsRequest, error) {
	var huaweiAdsRequest HuaweiAdsRequest
	var multislot []adslot30
	var publishersCredential *PublishersCredential
	for _, imp := range openRTBRequest.Imp {
		var err error
		publishersCredential, err = GetPublishersCredentials(&imp)
		if err != nil {
			return nil, err
		}

		if publishersCredential == nil {
			return nil, errors.New("publishers credentials is not complete!")
		}

		adslot30, err1 := getReqAdslot30(publishersCredential, &imp)
		if err != nil {
			return nil, err1
		}

		multislot = append(multislot, adslot30)
	}
	huaweiAdsRequest.Multislot = multislot
	huaweiAdsRequest.ClientAdRequestId = openRTBRequest.ID
	countryCode, err := getReqJson(&huaweiAdsRequest, openRTBRequest)
	if err != nil {
		return nil, err
	}
	reqJSON, err := json.Marshal(huaweiAdsRequest)
	if err != nil {
		return nil, err
	}
	var isTestAuthorization = false
	if publishersCredential != nil && publishersCredential.IsTestAuthorization == "true" {
		isTestAuthorization = true
	}
	header := getHeaders(publishersCredential, openRTBRequest, isTestAuthorization)
	bidRequest := RequestData{
		Method:  http.MethodPost,
		Uri:     getFinalEndPoint(countryCode),
		Body:    reqJSON,
		Headers: header,
	}

	log.Println(bidRequest)
	return &huaweiAdsRequest, nil
}

func GetPublishersCredentials(openRTBImp *openrtb2.Imp) (*PublishersCredential, error) {
	// var bidderExt ExtImpBidder
	huaweiAdsImpExt := PublishersCredential{
		SlotId:      "1",
		Adtype:      "2",
		PublisherId: "3",
		SignKey:     "4",
		KeyId:       "5",
	}

	// if err := json.Unmarshal(openRTBImp.Ext, &bidderExt); err != nil {
	// 	return nil, errors.New("Unmarshal: openRTBImp.Ext -> bidderExt failed")
	// }
	// if err := json.Unmarshal(bidderExt.Bidder, &huaweiAdsImpExt); err != nil {
	// 	return nil, errors.New("Unmarshal: bidderExt.Bidder -> huaweiAdsImpExt failed")
	// }
	// if huaweiAdsImpExt.SlotId == "" {
	// 	return nil, errors.New("ExtImpHuaweiAds: slotid is empty.")
	// }
	// if huaweiAdsImpExt.Adtype == "" {
	// 	return nil, errors.New("ExtImpHuaweiAds: adtype is empty.")
	// }
	// if huaweiAdsImpExt.PublisherId == "" {
	// 	return nil, errors.New("ExtHuaweiAds: publisherid is empty.")
	// }
	// if huaweiAdsImpExt.SignKey == "" {
	// 	return nil, errors.New("ExtHuaweiAds: signkey is empty.")
	// }
	// if huaweiAdsImpExt.KeyId == "" {
	// 	return nil, errors.New("ExtImpHuaweiAds: keyid is empty.")
	// }
	return &huaweiAdsImpExt, nil
}

func getReqAdslot30(publishersCredential *PublishersCredential, openRTBImp *openrtb2.Imp) (adslot30, error) {
	adtype := GetAdtype(publishersCredential.Adtype)
	testStatus := GetTestStatus(publishersCredential.IsTestAuthorization)
	var adslot30 = adslot30{
		Slotid: publishersCredential.SlotId,
		Adtype: adtype,
		Test:   testStatus,
	}
	if err := checkAndExtractOpenrtbFormat(&adslot30, adtype, publishersCredential.Adtype, openRTBImp); err != nil {
		return adslot30, err
	}
	return adslot30, nil
}

func GetAdtype(adtype string) int32 {
	switch strings.ToLower(adtype) {
	case "banner":
		return banner
	case "native":
		return native
	case "rewarded":
		return rewarded
	case "interstitial":
		return interstitial
	case "roll":
		return roll
	case "splash":
		return splash
	case "magazinelock":
		return magazinelock
	case "audio":
		return audio
	default:
		return banner
	}
}

func GetTestStatus(testStatus string) int32 {
	if testStatus == "true" {
		return 1
	}
	return 0
}

func checkAndExtractOpenrtbFormat(adslot30 *adslot30, adtype int32, yourAdtype string, openRTBImp *openrtb2.Imp) error {
	log.Println("a")
	if openRTBImp.Banner != nil {
		log.Println("b")
		if adtype != banner && adtype != interstitial {
			return errors.New("check openrtb format: request has banner, doesn't correspond to huawei adtype " + yourAdtype)
		}
		getBannerFormat(adslot30, openRTBImp)
	} else if openRTBImp.Native != nil {
		log.Println("c")
		if adtype != native {
			return errors.New("check openrtb format: request has native, doesn't correspond to huawei adtype " + yourAdtype)
		}
		if err := getNativeFormat(adslot30, openRTBImp); err != nil {
			return err
		}
	} else if openRTBImp.Video != nil {
		log.Println("d")
		if adtype != banner && adtype != interstitial && adtype != rewarded && adtype != roll {
			return errors.New("check openrtb format: request has video, doesn't correspond to huawei adtype " + yourAdtype)
		}
		if err := getVideoFormat(adslot30, adtype, openRTBImp); err != nil {
			return err
		}
	} else if openRTBImp.Audio != nil {
		log.Println("e")
		return errors.New("check openrtb format: request has audio, not currently supported")
	} else {
		return errors.New("check openrtb format: please choose one of our supported type banner, native, or video")
	}
	return nil
}

func getReqJson(request *HuaweiAdsRequest, openRTBRequest *openrtb2.BidRequest) (countryCode string, err error) {
	request.Version = huaweiAdxApiVersion
	if countryCode, err = getReqAppInfo(request, openRTBRequest); err != nil {
		return "", err
	}
	if err = getReqDeviceInfo(request, openRTBRequest); err != nil {
		return "", err
	}
	getReqNetWorkInfo(request, openRTBRequest)
	getReqRegsInfo(request, openRTBRequest)
	getReqGeoInfo(request, openRTBRequest)
	getReqConsentInfo(request, openRTBRequest)
	return countryCode, nil
}

func getReqAppInfo(request *HuaweiAdsRequest, openRTBRequest *openrtb2.BidRequest) (countryCode string, err error) {
	var app app
	if openRTBRequest.App != nil {
		if openRTBRequest.App.Ver != "" {
			app.Version = openRTBRequest.App.Ver
		}
		if openRTBRequest.App.Name != "" {
			app.Name = openRTBRequest.App.Name
		}

		// bundle cannot be empty, we need package name.
		if openRTBRequest.App.Bundle != "" {
			app.Pkgname = openRTBRequest.App.Bundle
		} else {
			return "", errors.New("generate HuaweiAds AppInfo failed: openrtb BidRequest.App.Bundle is empty.")
		}

		if openRTBRequest.App.Content != nil && openRTBRequest.App.Content.Language != "" {
			app.Lang = openRTBRequest.App.Content.Language
		} else {
			app.Lang = "en"
		}
	}
	countryCode = getCountryCode(openRTBRequest)
	app.Country = countryCode
	request.App = app
	return countryCode, nil
}

// getReqDeviceInfo: get device information for HuaweiAds request
func getReqDeviceInfo(request *HuaweiAdsRequest, openRTBRequest *openrtb2.BidRequest) (err error) {
	var device device
	if openRTBRequest.Device != nil {
		device.Type = int32(openRTBRequest.Device.DeviceType)
		device.Useragent = openRTBRequest.Device.UA
		device.Os = openRTBRequest.Device.OS
		device.Version = openRTBRequest.Device.OSV
		device.Maker = openRTBRequest.Device.Make
		device.Model = openRTBRequest.Device.Model
		if device.Model == "" {
			device.Model = defaultModelName
		}
		device.Height = int32(openRTBRequest.Device.H)
		device.Width = int32(openRTBRequest.Device.W)
		device.Language = openRTBRequest.Device.Language
		device.Pxratio = float32(openRTBRequest.Device.PxRatio)
		var country = getCountryCode(openRTBRequest)
		device.BelongCountry = country
		device.LocaleCountry = country
		device.Ip = openRTBRequest.Device.IP
		device.Gaid = openRTBRequest.Device.IFA
	}

	// get oaid gaid imei in openRTBRequest.User.Ext.Data
	if err = getDeviceIDFromUserExt(&device, openRTBRequest); err != nil {
		return err
	}

	// IsTrackingEnabled = 1 - DNT
	if openRTBRequest.Device != nil && openRTBRequest.Device.DNT != nil {
		if device.Oaid != "" {
			device.IsTrackingEnabled = strconv.Itoa(1 - int(*openRTBRequest.Device.DNT))
		}
		if device.Gaid != "" {
			device.GaidTrackingEnabled = strconv.Itoa(1 - int(*openRTBRequest.Device.DNT))
		}
	}

	request.Device = device
	return nil
}

func getCountryCode(openRTBRequest *openrtb2.BidRequest) string {
	if openRTBRequest.Device != nil && openRTBRequest.Device.Geo != nil && openRTBRequest.Device.Geo.Country != "" {
		return convertCountryCode(openRTBRequest.Device.Geo.Country)
	} else if openRTBRequest.User != nil && openRTBRequest.User.Geo != nil && openRTBRequest.User.Geo.Country != "" {
		return convertCountryCode(openRTBRequest.User.Geo.Country)
	} else if openRTBRequest.Device != nil && openRTBRequest.Device.MCCMNC != "" {
		return getCountryCodeFromMCC(openRTBRequest.Device.MCCMNC)
	} else {
		return defaultCountryName
	}
}

// convertCountryCode: ISO 3166-1 Alpha3 -> Alpha2, Some countries may use
func convertCountryCode(country string) (out string) {
	if country == "" {
		return defaultCountryName
	}
	var mapCountryCodeAlpha3ToAlpha2 = map[string]string{"AND": "AD", "AGO": "AO", "AUT": "AT", "BGD": "BD",
		"BLR": "BY", "CAF": "CF", "TCD": "TD", "CHL": "CL", "CHN": "CN", "COG": "CG", "COD": "CD", "DNK": "DK",
		"GNQ": "GQ", "EST": "EE", "GIN": "GN", "GNB": "GW", "GUY": "GY", "IRQ": "IQ", "IRL": "IE", "ISR": "IL",
		"KAZ": "KZ", "LBY": "LY", "MDG": "MG", "MDV": "MV", "MEX": "MX", "MNE": "ME", "MOZ": "MZ", "PAK": "PK",
		"PNG": "PG", "PRY": "PY", "POL": "PL", "PRT": "PT", "SRB": "RS", "SVK": "SK", "SVN": "SI", "SWE": "SE",
		"TUN": "TN", "TUR": "TR", "TKM": "TM", "UKR": "UA", "ARE": "AE", "URY": "UY"}
	if mappedCountry, exists := mapCountryCodeAlpha3ToAlpha2[country]; exists {
		return mappedCountry
	}

	if len(country) >= 3 {
		return country[0:2]
	}

	return defaultCountryName
}

func getCountryCodeFromMCC(MCC string) (out string) {
	var countryMCC = strings.Split(MCC, "-")[0]
	intVar, err := strconv.Atoi(countryMCC)

	if err != nil {
		return defaultCountryName
	} else {
		if result, found := constants.MccList[intVar]; found {
			return strings.ToUpper(result)
		} else {
			return defaultCountryName
		}
	}
}

// getDeviceID include oaid gaid imei. In prebid mobile, use TargetingParams.addUserData("imei", "imei-test");
// When ifa: gaid exists, other device id can be passed by TargetingParams.addUserData("oaid", "oaid-test");
func getDeviceIDFromUserExt(device *device, openRTBRequest *openrtb2.BidRequest) (err error) {
	var userObjExist = true
	if openRTBRequest.User == nil || openRTBRequest.User.Ext == nil {
		userObjExist = false
	}
	if userObjExist {
		var extUserDataHuaweiAds ExtUserDataHuaweiAds
		if err := json.Unmarshal(openRTBRequest.User.Ext, &extUserDataHuaweiAds); err != nil {
			return errors.New("get gaid from openrtb Device.IFA failed, and get device id failed: Unmarshal openRTBRequest.User.Ext -> extUserDataHuaweiAds. Error: " + err.Error())
		}

		var deviceId = extUserDataHuaweiAds.Data
		isValidDeviceId := false

		if len(deviceId.Oaid) > 0 {
			device.Oaid = deviceId.Oaid[0]
			isValidDeviceId = true
		}
		if len(deviceId.Gaid) > 0 {
			device.Gaid = deviceId.Gaid[0]
			isValidDeviceId = true
		}
		if len(deviceId.Imei) > 0 {
			device.Imei = deviceId.Imei[0]
			isValidDeviceId = true
		}

		if !isValidDeviceId {
			return errors.New("getDeviceID: Imei ,Oaid, Gaid are all empty.")
		}
		if len(deviceId.ClientTime) > 0 {
			device.ClientTime = getClientTime(deviceId.ClientTime[0])
		}
	} else {
		if len(device.Gaid) == 0 {
			return errors.New("getDeviceID: openRTBRequest.User.Ext is nil and device.Gaid is not specified.")
		}
	}
	return nil
}

func getClientTime(clientTime string) (newClientTime string) {
	var zone = defaultTimeZone
	t := time.Now().Local().Format(time.RFC822Z)
	index := strings.IndexAny(t, "-+")
	if index > 0 && len(t)-index == 5 {
		zone = t[index:]
	}
	if clientTime == "" {
		return time.Now().Format(timeFormat) + zone
	}
	if isMatched, _ := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}[+-]{1}\\d{4}$", clientTime); isMatched {
		return clientTime
	}
	if isMatched, _ := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}$", clientTime); isMatched {
		return clientTime + zone
	}
	return time.Now().Format(timeFormat) + zone
}

// getReqNetWorkInfo: for HuaweiAds request, include Carrier, Mcc, Mnc
func getReqNetWorkInfo(request *HuaweiAdsRequest, openRTBRequest *openrtb2.BidRequest) {
	if openRTBRequest.Device != nil {
		var network network
		if openRTBRequest.Device.ConnectionType != nil {
			network.Type = int32(*openRTBRequest.Device.ConnectionType)
		} else {
			network.Type = defaultUnknownNetworkType
		}

		var cellInfos []cellInfo
		if openRTBRequest.Device.MCCMNC != "" {
			var arr = strings.Split(openRTBRequest.Device.MCCMNC, "-")
			network.Carrier = 0
			if len(arr) >= 2 {
				cellInfos = append(cellInfos, cellInfo{
					Mcc: arr[0],
					Mnc: arr[1],
				})
				var str = arr[0] + arr[1]
				if str == "46000" || str == "46002" || str == "46007" {
					network.Carrier = 2
				} else if str == "46001" || str == "46006" {
					network.Carrier = 1
				} else if str == "46003" || str == "46005" || str == "46011" {
					network.Carrier = 3
				} else {
					network.Carrier = 99
				}
			}
		}
		network.CellInfo = cellInfos
		request.Network = network
	}
}

// getReqRegsInfo: get regs information for HuaweiAds request, include Coppa
func getReqRegsInfo(request *HuaweiAdsRequest, openRTBRequest *openrtb2.BidRequest) {
	if openRTBRequest.Regs != nil && openRTBRequest.Regs.COPPA >= 0 {
		var regs regs
		regs.Coppa = int32(openRTBRequest.Regs.COPPA)
		request.Regs = regs
	}
}

// getReqGeoInfo: get geo information for HuaweiAds request, include Lon, Lat, Accuracy, Lastfix
func getReqGeoInfo(request *HuaweiAdsRequest, openRTBRequest *openrtb2.BidRequest) {
	if openRTBRequest.Device != nil && openRTBRequest.Device.Geo != nil {
		var geo geo
		geo.Lon = float32(openRTBRequest.Device.Geo.Lon)
		geo.Lat = float32(openRTBRequest.Device.Geo.Lat)
		geo.Accuracy = int32(openRTBRequest.Device.Geo.Accuracy)
		geo.Lastfix = int32(openRTBRequest.Device.Geo.LastFix)
		request.Geo = geo
	}
}

// getReqGeoInfo: get GDPR consent
func getReqConsentInfo(request *HuaweiAdsRequest, openRTBRequest *openrtb2.BidRequest) {
	if openRTBRequest.User != nil && openRTBRequest.User.Ext != nil {
		var extUser ExtUser
		if err := json.Unmarshal(openRTBRequest.User.Ext, &extUser); err != nil {
			// fmt.Errorf("failed to parse ExtUser in HuaweiAds GDPR check: %v", err)
			return
		}
		request.Consent = extUser.Consent
	}
}

func getBannerFormat(adslot30 *adslot30, openRTBImp *openrtb2.Imp) {
	if openRTBImp.Banner.W != nil && openRTBImp.Banner.H != nil {
		adslot30.W = *openRTBImp.Banner.W
		adslot30.H = *openRTBImp.Banner.H
	}
	if len(openRTBImp.Banner.Format) != 0 {
		var formats = make([]format, 0, len(openRTBImp.Banner.Format))
		for _, f := range openRTBImp.Banner.Format {
			if f.H != 0 && f.W != 0 {
				formats = append(formats, format{f.W, f.H})
			}
		}
		adslot30.Format = formats
	}
}

func getNativeFormat(adslot30 *adslot30, openRTBImp *openrtb2.Imp) error {
	if openRTBImp.Native.Request == "" {
		return errors.New("extract openrtb native failed: imp.Native.Request is empty")
	}

	var nativePayload nativeRequests.Request
	if err := json.Unmarshal(json.RawMessage(openRTBImp.Native.Request), &nativePayload); err != nil {
		return err
	}

	// only compute the main image number, type = native1.ImageAssetTypeMain
	var numMainImage = 0
	var numVideo = 0
	var width int64
	var height int64
	for _, asset := range nativePayload.Assets {
		// Only one of the {title,img,video,data} objects should be present in each object.
		if asset.Video != nil {
			numVideo++
			continue
		}
		// every image has the same W, H.
		if asset.Img != nil {
			if asset.Img.Type == native1.ImageAssetTypeMain {
				numMainImage++
				if asset.Img.H != 0 && asset.Img.W != 0 {
					width = asset.Img.W
					height = asset.Img.H
				} else if asset.Img.WMin != 0 && asset.Img.HMin != 0 {
					width = asset.Img.WMin
					height = asset.Img.HMin
				}
			}
			continue
		}
	}
	adslot30.W = width
	adslot30.H = height

	var detailedCreativeTypeList = make([]string, 0, 2)
	if numVideo >= 1 {
		detailedCreativeTypeList = append(detailedCreativeTypeList, "903")
	} else if numMainImage > 1 {
		detailedCreativeTypeList = append(detailedCreativeTypeList, "904")
	} else if numMainImage == 1 {
		detailedCreativeTypeList = append(detailedCreativeTypeList, "901")
	} else {
		detailedCreativeTypeList = append(detailedCreativeTypeList, "913", "914")
	}
	adslot30.DetailedCreativeTypeList = detailedCreativeTypeList
	return nil
}

// roll ad need TotalDuration
func getVideoFormat(adslot30 *adslot30, adtype int32, openRTBImp *openrtb2.Imp) error {
	adslot30.W = openRTBImp.Video.W
	adslot30.H = openRTBImp.Video.H

	if adtype == roll {
		if openRTBImp.Video.MaxDuration == 0 {
			return errors.New("extract openrtb video failed: MaxDuration is empty when huaweiads adtype is roll.")
		}
		adslot30.TotalDuration = int32(openRTBImp.Video.MaxDuration)
	}
	return nil
}

func getFinalEndPoint(countryCode string) string {
	if countryCode == "" || len(countryCode) > 2 {
		return defaultEndpoint
	}
	var europeanSiteCountryCodeGroup = map[string]empty{"AX": {}, "AL": {}, "AD": {}, "AU": {}, "AT": {}, "BE": {},
		"BA": {}, "BG": {}, "CA": {}, "HR": {}, "CY": {}, "CZ": {}, "DK": {}, "EE": {}, "FO": {}, "FI": {},
		"FR": {}, "DE": {}, "GI": {}, "GR": {}, "GL": {}, "GG": {}, "VA": {}, "HU": {}, "IS": {}, "IE": {},
		"IM": {}, "IL": {}, "IT": {}, "JE": {}, "YK": {}, "LV": {}, "LI": {}, "LT": {}, "LU": {}, "MT": {},
		"MD": {}, "MC": {}, "ME": {}, "NL": {}, "AN": {}, "NZ": {}, "NO": {}, "PL": {}, "PT": {}, "RO": {},
		"MF": {}, "VC": {}, "SM": {}, "RS": {}, "SX": {}, "SK": {}, "SI": {}, "ES": {}, "SE": {}, "CH": {},
		"TR": {}, "UA": {}, "GB": {}, "US": {}, "MK": {}, "SJ": {}, "BQ": {}, "PM": {}, "CW": {}}
	var russianSiteCountryCodeGroup = map[string]empty{"RU": {}}
	var chineseSiteCountryCodeGroup = map[string]empty{"CN": {}}
	// choose site
	if _, exists := chineseSiteCountryCodeGroup[countryCode]; exists {
		return chineseSiteEndPoint
	} else if _, exists := russianSiteCountryCodeGroup[countryCode]; exists {
		return russianSiteEndPoint
	} else if _, exists := europeanSiteCountryCodeGroup[countryCode]; exists {
		return europeanSiteEndPoint
	} else {
		return asianSiteEndPoint
	}
}

func getHeaders(huaweiAdsImpExt *PublishersCredential, request *openrtb2.BidRequest, isTestAuthorization bool) http.Header {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")
	if huaweiAdsImpExt == nil {
		return headers
	}
	headers.Add("Authorization", getDigestAuthorization(huaweiAdsImpExt, isTestAuthorization))

	if request.Device != nil && len(request.Device.UA) > 0 {
		headers.Add("User-Agent", request.Device.UA)
	}
	return headers
}

func getDigestAuthorization(huaweiAdsImpExt *PublishersCredential, isTestAuthorization bool) string {
	var nonce = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	// this is for test case, time 2021/8/20 19:30
	if isTestAuthorization {
		nonce = "1629473330823"
	}
	var apiKey = huaweiAdsImpExt.PublisherId + ":ppsadx/getResult:" + huaweiAdsImpExt.SignKey
	return "Digest username=" + huaweiAdsImpExt.PublisherId + "," +
		"realm=ppsadx/getResult," +
		"nonce=" + nonce + "," +
		"response=" + computeHmacSha256(nonce+":POST:/ppsadx/getResult", apiKey) + "," +
		"algorithm=HmacSHA256,usertype=1,keyid=" + huaweiAdsImpExt.KeyId
}

func computeHmacSha256(message string, signKey string) string {
	h := hmac.New(sha256.New, []byte(signKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
