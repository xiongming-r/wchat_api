package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const HEADER_X_CONTENT_MD5 = "X-Content-MD5"
const HEADER_X_VERSION = "X-Sign-Version"
const HEADER_X_TIMESTAMP = "X-Timestamp"
const HEADER_X_NONCE = "X-Nonce"
// const accessKey = "0VHgS2PKU3bB1OMOH9LvA1oaDE"
const secretKey = "c5FrmwuwNMPOAMMT82Hie4tkYdk"

//const body = `{ "segmentDTOList": [{ "polygons": [{ "filterType": [ 2 ], "startDate": 2, "value": "[ [106.769048, 29.649794], [106.554945, 29.581484], [113.784273, 23.226937], [113.336932, 23.122341], [113.541928, 23.122284], [114.216637, 30.504789], [114.185324, 30.481072], [121.639907, 31.28643], [121.627924, 31.281475], [109.581154, 24.399206], [125.265613, 43.877379], [125.222859, 43.857183], [121.194974, 31.283154], [121.194309, 31.293923], [116.441518, 39.977133], [121.504601, 31.255256], [116.548518, 40.094699], [123.163908, 41.683483], [116.465288, 39.962046], [121.061136, 31.733576], [121.488855, 31.161518], [113.338961, 23.140184], [104.065535, 30.587856], [116.498827, 39.912432], [121.252796, 30.34069], [120.225391, 30.215155], [121.482679, 31.184247], [120.227439, 30.212953], [120.230987, 30.214679], [121.459118, 31.171379], [121.35092, 31.142298], [121.53083, 38.861143], [116.495642, 40.00542], [123.439774, 41.707437], [116.248512, 40.086023], [121.404486, 31.174327], [121.476839, 31.238141], [121.615111, 31.250975] ]" }] }] }`
// const estimateTempalte = `[{ "polygons": [%s] } ,{ "tagCodes":["102001002"] } ,{ "tagCodes":["102013001","102013003","102013004","102013005"] } ]` //高消费
const estimateTempalte = `[{ "polygons": [%s] } ,{ "tagCodes":[%s] } ,{ "tagCodes":[%s] } ]` //高消费
//const estimateTempalte = `[{ "polygons": [%s] } ,{ "tagCodes":["102001002"] } ,{ "tagCodes":["103008001","103008002","103008004"] } ]` //高档住宅
//const estimateTempalte = `[{ "polygons": [%s] } ,{ "tagCodes":["102001002"] } ]` //普通人群
//const estimateTempalte = `[{ "polygons": [%s] }]`
//const createTempalte = `{ "name": "ECARX投放项目人群包筛选", "desc": "ECARX投放项目，需要约100w设备号的人群包，用于在巨量引擎进行投放测试，验证平台投放流程的可用性", "segmentDTOList": [{ "polygons": [%s] },{ "tagCodes": [%s] }] }`
const createTempalte = `{ "name": "梵誓投放项目", "desc": "梵誓投放项目", "segmentDTOList": [{ "polygons": [%s] } ,{ "tagCodes": ["102001002"] } ,{ "tagCodes": ["102013001","102013003","102013004","102013005"] }] }` //高消费
//const createTempalte = `{ "name": "梵誓投放项目", "desc": "梵誓投放项目", "segmentDTOList": [{ "polygons": [%s] } ,{ "tagCodes": ["102001002"] } ,{ "tagCodes": ["103008001","103008002","103008004"] }] }` //高档住宅
//const createTempalte = `{ "name": "梵誓投放项目", "desc": "梵誓投放项目", "segmentDTOList": [{ "polygons": [%s] } ,{ "tagCodes": ["102001002"] }] }`  //普通人群
// const estimateURL = "https://lbi-api.newayz.com/openapi/v1/cloud/segment/estimate"
// const createURL = "https://lbi-api.newayz.com/openapi/v1/cloud/segment/create"
// const segmentURL = "https://lbi-api.newayz.com/openapi/v1/segment/tags?segmentId=%d"
// const segmentIdQueryURL = "https://lbi-api.newayz.com/openapi/v1/cloud/segment/ids?segmentIds=%d"

// const authURL = "https://lbi-api.newayz.com/openapi/precisionMarketing/threeParty/toAuthorizeUrl?advertiserId=1711765206912014&scope=%5B3%2C4%2C14%5D"
// const authURL = "https://lbi-api.newayz.com/openapi/precisionMarketing/threeParty/toAuthorizeUrl?advertiserId=1708510001353741&scope=3%2C4%2C14"
// const uploadURL = "https://lbi-api.newayz.com/openapi/precisionMarketing/threeParty/getJuLiangIdForClient"
// const pushURL = "https://lbi-api.newayz.com/openapi/precisionMarketing/threeParty/push"
// const publishURL = "https://lbi-api.newayz.com/openapi/precisionMarketing/threeParty/publish"

// const queryURL = "https://lbi-api.newayz.com/openapi/precisionMarketing/threeParty/query/juLiang?advertiserId=1708510001353741&wayzCrowdId=%d"
// const queryURL = "https://lbi-api.newayz.com/openapi/precisionMarketing/threeParty/query/juLiang?advertiserId=%d&wayzCrowdId=%d"

//var factory = flag.Int("factory", 1500, "the species we are studying")
//var office = flag.Int("office", 50, "the species we are studying")
// 这个下面的变量都是命令行的参数 | 参数名称 | 默认值 | 描述
var action = flag.String("action", "", "the species we are studying")
var crowd = flag.Int("crowd", -1, "the species we are studying")
var advertiserId = flag.Int("advertiser-id", -1, "advertiser-id")
var segmentId = flag.Int("segment-id", -1, "the species we are studying")
var accessKey = flag.String("access-key","0VHgS2PKU3bB1OMOH9LvA1oaDE", "access key")
var advertiserId = flag.String("advertiser-id","1708510001353741", "advertiser-id")
var estimateURL = flag.String("estimate-URL","","estimate-URL")
var createURL = flag.String("create-URL","","create-URL")
var segmentURL = flag.String("segment-URL","","segment-URL")
var segmentIdQueryURL = flag.String("segmentIdQuery-URL","","segmentIdQuery-URL")
var authURL = flag.String("auth-URL","","auth-URL")
var uploadURL = flag.String("upload-URL","","upload-URL")
var pushURL = flag.String("push-URL","","push-URL")
var publishURL = flag.String("publish-URL","","publish-URL")
var queryURL = flag.String("query-URL","","query-URL")
// var estimateTempalte = flag.String("estimate-Tempalte","","estimate-Tempalte")
// var createTempalte = flag.String("create-Tempalte","","create-Tempalte")

var tagCode = flag.String("tag-Code","","tag-Code")
var tagCodes = flag.String("tag-Codes","","tag-Codes")


func estimate() (string, *http.Request, error) {
	p, _ := autoGen()
	body := fmt.Sprintf(estimateTempalte, p,*tagCode,*tagCodes)
	// body := *estimateTempalte
	fmt.Printf("%s\n\n", body)
	req, err := http.NewRequest("POST", *estimateURL, bytes.NewBuffer([]byte(body)))
	return md5Body(body), req, err
}

func create() (string, *http.Request, error) {
	//p, t := autoGen() 为了让create可以直接加tagcode
	// p, _ := autoGen()
	//body := fmt.Sprintf(createTempalte, p, t)   为了让create可以直接加tagcode
	// body := fmt.Sprintf(createTempalte, p)
	body := *createTempalte
	fmt.Printf("%s\n\n", body)
	req, err := http.NewRequest("POST", *createURL, bytes.NewBuffer([]byte(body)))
	return md5Body(body), req, err
}
func getSign(req *http.Request) () {
	// flag.Parse()
	// var (
	// 	// md5 string
	// 	req *http.Request
	// 	// err error
	// 	// code = string
	// )
	//p, t := autoGen() 为了让create可以直接加tagcode
	// sigs := NewSignature(req)
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add(HEADER_X_CONTENT_MD5, md5)
	// req.Header.Add(HEADER_X_VERSION, "1.0")
	// req.Header.Add(HEADER_X_NONCE, genUUID())
	// req.Header.Add(HEADER_X_TIMESTAMP, strconv.FormatInt(time.Now().Unix(), 10))
	// //req.Header.Add(HEADER_X_TIMESTAMP, strconv.FormatInt(time.Date(2021, time.July, 17, 0, 0, 0, 0, time.UTC).Unix(), 10))
	// req.Header.Add("Authorization", sig.Sign())
	// code = sig.Sign()
	//body := fmt.Sprintf(createTempalte, p, t)   为了让create可以直接加tagcode
	// body := fmt.Sprintf(createTempalte, p)
	// fmt.Printf("%s\n\n", body)
	// req, err := http.NewRequest("POST", createURL, bytes.NewBuffer([]byte(body)))
	// return sig.Sign()
	fmt.Printf("%d\n", req.Header)
	return
}

func segment() (string, *http.Request, error) {
	url := fmt.Sprintf(*segmentURL, *segmentId)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	return md5Body(""), req, err
}

func segmentIdQuery() (string, *http.Request, error) {
	url := fmt.Sprintf(*segmentIdQueryURL, *segmentId)
	println(url)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	return md5Body(""), req, err
}

func auth() (string, *http.Request, error) {
	req, err := http.NewRequest("GET", *authURL, bytes.NewBuffer([]byte("")))
	return md5Body(""), req, err
}

func query() (string, *http.Request, error) {
	qURL := fmt.Sprintf(*queryURL, *advertiserId, *crowd)
	req, err := http.NewRequest("GET", qURL, bytes.NewBuffer([]byte("")))
	return md5Body(""), req, err
}

func upload() (string, *http.Request, error) {
	body := fmt.Sprintf(`{"wayzCrowdId":%d, "advertiserId":%d}`, *crowd,*advertiserId)
	fmt.Printf("%s\n\n", body)
	req, err := http.NewRequest("POST", *uploadURL, bytes.NewBuffer([]byte(body)))
	return md5Body(body), req, err
}

func push() (string, *http.Request, error) {
	body := `{"customAudienceId":334843577, "advertiserId":1708510001353741, "targetAdvertiserIds":[1708510001353741]}`
	req, err := http.NewRequest("POST", *pushURL, bytes.NewBuffer([]byte(body)))
	return md5Body(body), req, err
}

func publish() (string, *http.Request, error) {
	body := `{"customAudienceId":334843577, "advertiserId":1708510001353741}`
	req, err := http.NewRequest("POST", *publishURL, bytes.NewBuffer([]byte(body)))
	return md5Body(body), req, err
}

func main() {
	flag.Parse()
	fmt.Printf("action=%s\n", *action)

	var (
		md5 string
		req *http.Request
		err error
		// res string
	)
	switch *action {
	case "create":
		md5, req, err = create()
	case "estimate":
		md5, req, err = estimate()
	case "segment":
		md5, req, err = segment()
	case "auth":
		md5, req, err = auth()
	case "upload":
		md5, req, err = upload()
	case "push":
		md5, req, err = push()
	case "publish":
		md5, req, err = publish()
	case "query":
		md5, req, err = query()
	case "segmentIdQuery":
		md5, req, err = segmentIdQuery()
	case "getSign":
		getSign(req)
		
		// fmt.Printf(res)
	default:
		fmt.Printf("unknow action: '%s'\n", *action)
		return
	}

	if err != nil {
		fmt.Println(err)
	} else {
		sig := NewSignature(req)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add(HEADER_X_CONTENT_MD5, md5)
		req.Header.Add(HEADER_X_VERSION, "1.0")
		req.Header.Add(HEADER_X_NONCE, genUUID())
		req.Header.Add(HEADER_X_TIMESTAMP, strconv.FormatInt(time.Now().Unix(), 10))
		//req.Header.Add(HEADER_X_TIMESTAMP, strconv.FormatInt(time.Date(2021, time.July, 17, 0, 0, 0, 0, time.UTC).Unix(), 10))
		req.Header.Add("Authorization", sig.Sign())
		fmt.Printf("%v\n\n", req.Header)
		fmt.Printf(sig.Sign())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("failed, err=%v\n", err)
		} else {
			defer resp.Body.Close()

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("failed, err=%v\n", err)
			}
			bodyString := string(bodyBytes)
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("failed, resp=%v\n", bodyString)
			} else {
				fmt.Printf("ok, resp=%v\n", bodyString)
			}
		}
	}
}

type Signature struct {
	components []string
	request    *http.Request
}

func NewSignature(req *http.Request) *Signature {
	return &Signature{
		components: make([]string, 0),
		request:    req,
	}
}

func (s *Signature) Sign() string {
	//fmt.Printf("%v\n", s.request)
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(s.buildStringToSign()))
	str := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("wayz %s:%s", *accessKey, str)
}

func (s *Signature) buildStringToSign() string {
	s.add(*accessKey)
	s.add(s.request.Method)
	s.add(s.request.URL.Path)
	s.add(s.sortedParamStr())
	s.add(s.request.Header.Get("Accept"))
	s.add(s.request.Header.Get(HEADER_X_CONTENT_MD5))
	s.add(s.request.Header.Get("Content-Type"))
	s.add(s.request.Header.Get(HEADER_X_TIMESTAMP))
	s.add(s.request.Header.Get(HEADER_X_VERSION))
	s.add(s.request.Header.Get(HEADER_X_NONCE))
	return s.build()
}

func (s *Signature) sortedParamStr() string {
	keys := make([]string, 0)
	for k, _ := range s.request.URL.Query() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	kvs := make([]string, 0, len(keys))
	for _, k := range keys {
		v := s.request.URL.Query().Get(k)
		kvs = append(kvs, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(kvs, "&")
}

func (s *Signature) build() string {
	return strings.Join(s.components, "\r\n")
}

func (s *Signature) add(component string) {
	s.components = append(s.components, component)
}

func md5Body(body string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(body)))
}

func genUUID() string {
	uuidWithHyphen := uuid.New()
	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

func autoGen() (string, string) {
	points := []string{
		"121.475240&31.216590",
		"121.426939&31.210037",
		"116.433710&39.920980",
		"120.162410&30.260090",
		"118.787610&32.011701",
		"120.709864&31.321583",
		"114.264711&30.578856",
		"106.573522&29.576343",
		"121.478196&31.270738",
		"121.426939&31.210037",
		"121.412595&31.232962",
		"121.486047&31.238171",
		"121.462601&31.218525",
		"121.383530&31.245347",
		"121.414709&31.216709",
		"121.498930&31.236191",
		"121.513252&31.301154",
		"121.517833&31.299417",
		"121.437229&31.194816",
		"121.517570&31.227363",
		"121.545101&31.236498",
		"121.368951&31.170446",
		"121.400160&31.130959",
		"121.252425&31.331075",
		"121.569415&31.114466",
		"121.094983&31.144024",
		"116.372569&39.909710",
		"116.411731&39.915970",
		"116.352747&39.941756",
		"116.326592&39.787144",
		"116.414590&39.909421",
		"116.375358&39.908913",
		"116.314950&39.978627",
		"116.599626&39.925044",
		"116.460187&39.907508",
		"116.287337&39.958508",
		"116.468010&39.993210",
		"116.432774&39.920871",
		"116.475817&39.949311",
		"116.298894&39.824191",
		"116.141581&39.744316",
		"116.129218&39.911656",
		"116.325589&39.678294",
		"116.297023&39.908376",
		"116.320256&39.686835",
		"120.107566&30.299536",
		"120.141884&30.304046",
		"120.164514&30.252555",
		"120.206126&30.211463",
		"120.327336&30.309792",
		"120.185153&30.326174",
		"120.269216&30.181212",
		"120.213530&30.250566",
		"118.783441&32.039049",
		"118.820083&31.929672",
		"118.785953&32.038986",
		"118.722349&32.139833",
		"118.784567&32.043369",
		"118.729000&32.003875",
		"120.685691&31.322459",
		"120.929214&31.380983",
		"120.624743& 31.142381",
		"120.715308&31.319862",
		"120.972669&31.375556",
		"121.137505&31.460750",
		"120.602274&31.311330",
		"120.622030&31.310499",
		"120.546212&31.861169",
		"120.556166&31.294453",
		"114.274048&30.578618",
		"114.340712&30.586558",
		"114.352160&30.560179",
		"114.291393&30.579370",
		"114.275550&30.586082",
		"114.402341&30.505339",
		"114.309319&30.608918",
		"114.402127&30.505756",
		"114.167299&30.617150",
		"106.526732& 29.600746",
		"106.577367&29.558602",
		"106.518471&29.516477",
		"106.533302&29.577503",
		"106.511763&29.538737",
		"106.588509&29.563428",
		"106.571149&29.522626",
		"106.573522&29.576343",
		"106.529182&29.650231",
		"106.293288&29.605269",
		"121.324125&31.239902",
		"121.569415&31.114466",
		"121.448547&31.324193",
		"121.243485&31.057065",
		"121.483973&30.915380",
		"121.462982&30.929843",
		"121.219500&31.037832",
		"121.095923&31.144313",
		"121.252425&31.331075",
		"121.305986&31.304831",
		"121.446218&31.223982",
		"121.475117&31.235181",
		"121.416659&31.219129",
		"121.517409&31.229028",
		"121.378497&31.107354",
		"121.476271&31.232659",
		"121.335109&30.754436",
		"121.412630&31.233195",
		"121.462798&31.230082",
		"121.439651&31.192950",
		"121.370830&31.171530",
		"121.406454&31.071131",
		"121.447750&31.324000",
		"121.486047&31.238171",
		"121.414169&31.025234",
		"121.469236&31.221173",
		"121.416512&31.219580",
		"121.570851&31.114451",
		"121.437820&31.194890",
		"121.570851&31.114451",
		"121.472340&31.234598",
		"121.498930&31.236191",
		"121.764360&31.197446",
		"121.512889&31.301548",
		"121.517972&31.299904",
		"121.456988&31.228342",
		"116.375225&39.909381",
		"116.415547&40.059765",
		"116.414590&39.909421",
		"116.417894&39.897555",
		"116.373280&39.909911",
		"116.479027&39.910227",
		"116.410118&39.913767",
		"116.435347&39.940327",
		"116.519073&39.924529",
		"116.460849&39.895499",
		"116.468038&39.993481",
		"116.225949&39.905965",
		"116.332033&40.027918",
		"116.476123&39.949317",
		"116.440439&39.921741",
		"116.419844&39.898257",
		"116.501577&39.805441",
		"116.652469&40.127696",
		"116.597508&39.924930",
		"116.641114&39.905129",
		"116.288855&39.958990",
		"116.367876&39.813795",
		"116.835416&40.356078",
		"116.478548&39.893776",
		"116.171237&39.922886",
		"116.355851&39.860936",
		"116.629818&40.316679",
		"116.140617&39.74553",
		"116.548061&39.951677",
		"116.352747&39.941756",
		"116.420667&39.542956",
		"116.373544&39.910495",
		"116.315182&39.978316",
		"116.432458&39.940193",
		"116.356602&39.952927",
		"116.411660&39.915810",
		"116.447452&39.971439",
		"116.298284&39.824280",
		"116.460481&39.911856",
		"116.366225&39.854308",
		"116.595941&40.012906",
		"116.311594&40.028549",
		"116.341819&39.730177",
		"116.339339&39.991694",
		"116.659730&39.908917",
		"116.327024&39.787486",
		"116.239161&40.212115",
		"116.317465&40.067073",
		"116.123085&39.905659",
		"120.164394&30.252166",
		"120.107559&30.299526",
		"120.174114&30.327325",
		"120.292536&30.398450",
		"120.050504&30.247484",
		"120.163811&30.268079",
		"120.328500&30.309577",
		"120.162398&30.260845",
		"120.379250&30.337704",
		"120.204806&30.202733",
		"120.320294&30.308888",
		"120.212370&30.249109",
		"120.204625&30.259528",
		"120.269216&30.181212",
		"120.215364&30.251091",
		"120.119315&30.330456",
		"119.933344&30.089688",
		"120.164569&30.268779",
		"118.736256&32.032699",
		"118.843889&31.950780",
		"118.783109&32.038671",
		"118.724501&32.138236",
		"118.784700&32.044470",
		"118.785848&32.023570",
		"118.785381&32.038709",
		"118.779792&32.041479",
		"118.785199&32.040315",
		"118.820196&31.929918",
		"118.729000&32.003875",
		"118.698189&32.158664",
		"118.784766&32.091408",
		"118.657666&31.912402",
		"118.740245&32.023315",
		"118.722256&32.140479",
		"118.730849&32.125566",
		"118.783204&32.072722",
		"120.957475&31.375569",
		"120.622030&31.310499",
		"120.549210&31.283000",
		"120.545258&31.860961",
		"120.973477&31.375907",
		"120.955859&31.404539",
		"120.539700&31.864500",
		"120.539700&31.864500",
		"120.554422&31.293774",
		"120.608272&31.253170",
		"120.631921&31.414017",
		"120.760995&31.381710",
		"120.712498&31.323269",
		"120.758161&31.649164",
		"114.405277&30.506243",
		"114.405277&30.506243",
		"114.269796&30.580618",
		"114.309877&30.607751",
		"114.290644&30.577557",
		"114.342290&30.587631",
		"114.405208&30.505792",
		"114.290114&30.577408",
		"114.208968&30.560355",
		"114.377176&30.626728",
		"114.344730&30.554410",
		"114.268639&30.610261",
		"114.167652&30.616631",
		"114.411885&30.493288",
		"106.516999&29.510445",
		"106.568975&29.525250",
		"106.570625&29.523063",
		"107.398498&29.699892",
		"108.377972&30.809361",
		"106.532647&29.575069",
		"106.512359&29.535564",
		"108.383568&30.812136",
		"106.517658&29.516111",
		"106.253610&29.286682",
		"108.404499&31.170139",
		"106.544298&29.398654",
		"106.579528&29.558033",
		"105.925755&29.337105",
		"106.293500&29.605725",
		"106.568975&29.525250",
		"106.517385&29.510155",
		"106.567278&29.479740",
		"106.532647&29.575069",
		"106.532647&29.575069",
		"106.548486&29.641906",
		"106.532710&29.573794",
		"106.544347&29.398331",
		"106.462497&29.556759",
		"106.460578&29.553905",
		"106.475230&29.544385",
		"106.577349&29.558656",
		"106.532647&29.575069",
		"106.530680&29.652760",
		"106.533174&29.575669",
		"121.446350&31.225599",
		"116.475459&39.907243",
		"120.161499&30.273866",
		"118.785127&32.040776",
		"120.618866&31.316188",
		"106.533775&29.573132",
		"121.471180&31.265024",
		"121.337040&30.756073",
		"121.084380&31.135226",
		"121.511317&31.02753",
		"121.242947&31.057553",
		"121.601669&31.203772",
		"121.371776&31.023849",
		"121.252425&31.331075",
		"121.569415&31.114466",
		"121.324824&31.240346",
		"116.412554&39.915043",
		"116.482066&39.893687",
		"116.472155&39.909427",
		"116.415538&40.059723",
		"115.984869&40.473899",
		"120.161499&30.273866",
		"120.204798&30.201000",
		"120.1658&30.249794",
		"118.780708&32.041094",
		"118.734633&32.032971",
		"120.555416&31.294308",
		"120.554742&31.866642",
		"106.587236&29.564809",
	}
	values1 := make([]string, 0, len(points))
	for _, p := range points {
		values1 = append(values1, fmt.Sprintf(`{ "filterType": [0], "startDate": 3, "value": "%d&%s" }`, 5000, p))
	}
	codes := []string{
		"102002001",
		"102002002",
	}

	values2 := make([]string, 0, len(codes))
	for _, c := range codes {
		values2 = append(values2, fmt.Sprintf(`"%s"`, c))
	}
	return strings.Join(values1, ","), strings.Join(values2, ",")
}
