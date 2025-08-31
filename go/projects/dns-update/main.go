package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	GlobaLDefaults = Config{}
	apiKey         string
	zoneName       string
	recordName     []string
	recordType     int
	apiUrl         = "https://api.cloudflare.com"
	ipAddress      string
)

func GetAPIKey() string {
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if kv[0] == "CLOUDFLARE_API_KEY" {
			// fmt.Println(kv[0] + ": " + kv[1])
			return kv[1]
		}
	}
	return ""
}

func GetZoneID(apiKey string, apiUrl string, zoneName string) string {
	requestURL := fmt.Sprintf("%s/client/v4/zones?name=%s", apiUrl, zoneName)
	client := &http.Client{}
	body := &bytes.Buffer{}
	req, _ := http.NewRequest("GET", requestURL, body)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	client.Do(req)
	bodyResponse, _ := io.ReadAll(req.Body)
	result := &CloudFlareDnsZoneResponse{}
	json.Unmarshal(bodyResponse, result)
	return result.Result[0].ID
}

func GetCloudFlareDnsZone(apiKey string, apiUrl string, zoneName string) *CloudFlareDnsZoneResponse {
	requestURL := fmt.Sprintf("%s/client/v4/zones?name=%s", apiUrl, zoneName)
	fmt.Println(requestURL)
	client := &http.Client{}
	body := &bytes.Buffer{}
	req, _ := http.NewRequest("GET", requestURL, body)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	response, _ := client.Do(req)
	bodyResponse, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading body: %s\n", err)
	}
	rawBody := map[string]any{}
	json.Unmarshal(bodyResponse, &rawBody)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(rawBody)
	result := &CloudFlareDnsZoneResponse{}
	json.Unmarshal(bodyResponse, result)
	return result
}

func GetDNSRecord(apiKey string, apiUrl string, zoneID string, recordName string) *CloudFlareDnsRecordResponse {
	rawURL := fmt.Sprintf("%s/client/v4/zones/%s/dns_records?name=%s", apiUrl, zoneID, recordName)
	client := &http.Client{}
	body := &bytes.Buffer{}
	request, _ := http.NewRequest("GET", rawURL, body)
	request.Header.Add("Authorization", "Bearer "+apiKey)
	client.Do(request)
	data, _ := io.ReadAll(request.Body)
	recordResult := &CloudFlareDnsRecordResponse{}
	json.Unmarshal(data, recordResult)
	return recordResult
}

func BoolPtr(b bool) *bool {
	return &b
}

func ReconcileRecordSettingndExistingEntry(setting *RecordSetting, existingRecord CloudFlareDnsRecord) {
	switch {
	case setting.Ttl == 0:
		setting.Ttl = existingRecord.Ttl
	case setting.Type == "":
		setting.Type = existingRecord.Type
	case setting.Proxied == nil:
		setting.Proxied = BoolPtr(true)
	}
}

func ExtractSettingsFromRecords(settings []RecordSetting, existingRecords []CloudFlareDnsRecord) (payload CloudFlareBatchRecordPayload) {
	recordsByName := make(map[string]map[string][]*CloudFlareDnsRecord)
	for _, record := range existingRecords {
		recordsByName[record.Name][record.Type] = append(recordsByName[record.Name][record.Type], &record)
	}
	for _, setting := range settings {
		setting.SetZeroValues()
		if recordForSetting := recordsByName[setting.Name][setting.Type][0]; recordForSetting != nil {
			payload.Patches = append(payload.Patches,
				PatchRecord{
					RecordSetting: RecordSetting{
						Name:    setting.Name,
						Ttl:     setting.Ttl,
						Type:    setting.Type,
						Proxied: setting.Proxied,
						Comment: setting.Comment,
						Content: setting.Content,
					},
					RecordID: recordForSetting.ID,
				})
		} else {
			payload.Posts = append(payload.Posts,
				RecordSetting{
					Name:    setting.Name,
					Ttl:     setting.Ttl,
					Type:    setting.Type,
					Proxied: setting.Proxied,
					Comment: setting.Comment,
					Content: setting.Content,
				})
		}
	}
	return payload
}

func SetDnsRecordFromResponse(existingRecords *CloudFlareDnsRecordResponse, recordSettings []RecordSetting, ipAddress string, apiUrl string, apiKey string, zoneID string) []*http.Response {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	recordUpdates := ExtractSettingsFromRecords(recordSettings, existingRecords.Result)
	body := &bytes.Buffer{}
	payload, _ := json.Marshal(recordUpdates)
	body.Write(payload)
	encoder.Encode(recordUpdates)
	batchResponse := SetDnsRecordBatch(recordUpdates, apiUrl, apiKey, zoneID)
	if batchResponse.StatusCode != 200 {
		responses := SetDnsRecordsIndividually(recordUpdates, apiUrl, apiKey, zoneID)
		return responses
	}
	return []*http.Response{batchResponse}
}

func SetDnsRecordsIndividually(records CloudFlareBatchRecordPayload, apiUrl string, apiKey string, zoneID string) []*http.Response {
	body := &bytes.Buffer{}
	bodyEncoder := json.NewEncoder(body)
	client := http.Client{}
	responses := []*http.Response{}
	for _, patch := range records.Patches {
		body.Reset()
		url := fmt.Sprintf("%s/zones/%s/records/%s", apiUrl, zoneID, patch.RecordID)
		bodyEncoder.Encode(patch)
		request, _ := http.NewRequest("PATCH", url, body)
		response, _ := client.Do(request)
		responses = append(responses, response)
	}
	for _, post := range records.Posts {
		body.Reset()
		url := fmt.Sprintf("%s/zones/%s/records", apiUrl, zoneID)
		bodyEncoder.Encode(post)
		request, _ := http.NewRequest("POST", url, body)
		response, _ := client.Do(request)
		responses = append(responses, response)
	}
	return responses
}

func CreateCloudFlareDnsRecordPtr(name string, ttl int, r_type string, comment string, content string, proxied bool) *CloudFlareDnsRecord {
	record := &CloudFlareDnsRecord{
		Name:    name,
		Ttl:     ttl,
		Type:    r_type,
		Comment: comment,
		Content: content,
		Proxied: proxied,
	}
	return record
}

func SetDnsRecordBatch(batch CloudFlareBatchRecordPayload, apiUrl string, apiKey string, zoneID string) *http.Response {
	endpointURL := fmt.Sprintf("%s/zones/%s/dns_records/batch", apiUrl, zoneID)
	client := &http.Client{}
	fmt.Println(endpointURL)
	body := &bytes.Buffer{}
	bodyContent, _ := json.Marshal(batch)
	body.Write(bodyContent)
	request, _ := http.NewRequest("POST", endpointURL, body)
	request.Header.Add("Authorization", "Bearer "+apiKey)
	response, _ := client.Do(request)
	fmt.Printf("Batch Record Response: %d", response.StatusCode)
	responseContent, _ := io.ReadAll(response.Body)
	responseBody := &CloudFlareDnsRecordWriteResponse{}
	json.Unmarshal(responseContent, &responseBody)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(responseBody)
	return response
}

func SetDnsRecord(record *CloudFlareDnsRecord, method string, apiUrl string, apiKey string, zoneID string) *http.Response {
	var endpointUrl string
	if method == "PATCH" {
		endpointUrl = fmt.Sprintf("%s/client/v4/zones/%s/dns_records/%s", apiUrl, zoneID, record.ID)
	} else {
		endpointUrl = fmt.Sprintf("%s/client/v4/zones/%s/dns_records", apiUrl, zoneID)
	}
	client := &http.Client{}
	fmt.Println(endpointUrl)
	body := &bytes.Buffer{}
	bodyContent, _ := json.Marshal(record)
	body.Write(bodyContent)
	request, _ := http.NewRequest(method, endpointUrl, body)
	request.Header.Add("Authorization", "Bearer "+apiKey)
	response, _ := client.Do(request)
	fmt.Printf("Set Record Response: %d", response.StatusCode)
	responseContent, _ := io.ReadAll(response.Body)
	responseBody := &CloudFlareDnsRecordWriteResponse{}
	json.Unmarshal(responseContent, &responseBody)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(responseBody)
	return response
}

func GetConfig(fileName string, apiKey string) Config {
	var config Config
	file, err := os.Open(fileName)
	if err != nil {
		config = CreateDefaultConfig(apiKey)
	} else {
		config = LoadConfig(file)
	}
	return config
}

func CreateDefaultConfig(apiKey string) Config {
	config := Config{
		ApiKey:   apiKey,
		ZoneName: "harmonlab.io",
		RecordSettings: []RecordSetting{
			{
				Name:    "@",
				Type:    "A",
				Ttl:     3600,
				Proxied: BoolPtr(false),
			},
		},
		RecordType: "A",
		IpUrl:      "https://api.ipify.org",
	}
	return config
}

func LoadConfig(file *os.File) (config Config) {
	rawConfig, _ := io.ReadAll(file)
	yaml.Unmarshal(rawConfig, &config)
	return config
}

func GetIpAddress(ipUrl string) string {
	response, _ := http.Get(ipUrl)
	ipAddress, _ := io.ReadAll(response.Body)
	return string(ipAddress)
}

func main() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	config := GetConfig("config.yaml", apiKey)
	if config.ApiKey == "" {
		config.ApiKey = GetAPIKey()
	}
	// fmt.Println("api key: " + config.ApiKey)
	ipAddress := GetIpAddress(config.IpUrl)
	encoder.Encode(ipAddress)
	dnsZone := GetCloudFlareDnsZone(config.ApiKey, apiUrl, config.ZoneName)
	// encoder.Encode(dnsZone)
	dnsRecords := GetDNSRecord(config.ApiKey,
		apiUrl,
		dnsZone.Result[0].ID,
		config.RecordSettings[0].Name)
	setResponse := SetDnsRecordFromResponse(dnsRecords, config.RecordSettings, ipAddress, apiUrl, config.ApiKey, dnsZone.Result[0].ID)
	fmt.Println("Set response: ")
	encoder.Encode(setResponse)
}

type CloudFlareAPIError struct {
	Code              int    `json:"code"`
	Message           string `json:"message"`
	Documentation_url string `json:"documentation_url"`
	Source            struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
}

type CloudFlareAPIMessage struct {
	Code              int    `json:"code"`
	Message           string `json:"message"`
	Documentation_url string `json:"documentation_url"`
	Source            struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
}

type CloudFlareDnsZoneResult struct {
	ID      string `json:"id"`
	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Meta struct {
		Cdn_only                 bool `json:"cdn_only"`
		Custom_certificate_quota int  `json:"custom_certificate_quota"`
		Dns_noly                 bool `json:"dns_noly"`
		Foundation_dns           bool `json:"foundation_dns"`
		Page_rule_quota          int  `json:"page_rule_quota"`
		Phishing_detected        bool `json:"phishing_detected"`
		Step                     int  `json:"step"`
	} `json:"meta"`
	Name  string `json:"name"`
	Owner struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"owner"`
	Plan struct {
		ID                 string  `json:"id"`
		Can_subscribe      bool    `json:"can_subscribe"`
		Currency           string  `json:"currency"`
		Externally_managed bool    `json:"externally_managed"`
		Frequency          string  `json:"frequency"`
		Is_subscribed      bool    `json:"is_subscribed"`
		Legacy_discount    bool    `json:"legacy_discount"`
		Legacy_id          string  `json:"legacy_id"`
		Name               string  `json:"name"`
		Price              float64 `json:"price"`
	} `json:"plan"`
	Cname_suffix string   `json:"cname_suffix"`
	Paused       bool     `json:"paused"`
	Permissions  []string `json:"permissions"`
	Tenant       struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tenant"`
	Tenant_unit struct {
		ID string `json:"id"`
	} `json:"tenant_unit"`
	Type                string   `json:"type"`
	Vanity_name_servers []string `json:"vanity_name_servers"`
}

type CloudFlareDnsZoneResultInfo struct {
	Count       int `json:"count"`
	Page        int `json:"page"`
	Per_page    int `json:"per_page"`
	Total_count int `json:"total_count"`
	Total_pages int `json:"total_pages"`
}

type CloudFlareDnsZoneMessage struct {
	Code              int    `json:"code"`
	Message           string `json:"message"`
	Dovumentation_url string `json:"dovumentation_url"`
	Source            struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
}

type CloudFlareDnsZoneResponse struct {
	Errors      []CloudFlareAPIError        `json:"errors"`
	Messages    []CloudFlareAPIMessage      `json:"messages"`
	Success     bool                        `json:"success"`
	Result      []CloudFlareDnsZoneResult   `json:"result"`
	Result_info CloudFlareDnsZoneResultInfo `json:"result_info"`
}

type CloudFlareDnsRecordResponse struct {
	Errors   []CloudFlareAPIError   `json:"errors"`
	Messages []CloudFlareAPIMessage `json:"messages"`
	Success  bool                   `json:"success"`
	Result   []CloudFlareDnsRecord  `json:"result"`
}

type CloudFlareDnsRecordResult struct {
	Name    string `json:"name"`
	Ttl     int    `json:"ttl"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	Setting struct {
		Ipv4_only bool `json:"ipv4_only"`
		Ipv6_only bool `json:"ipv6_only"`
	} `json:"setting"`
	Tags []string `json:"tags"`
	ID   string   `json:"id"`
}

type CloudFlareDnsRecord struct {
	Name    string `json:"name"`
	Ttl     int    `json:"ttl"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	Setting struct {
		Ipv4_only bool `json:"ipv4_only"`
		Ipv6_only bool `json:"ipv6_only"`
	} `json:"setting"`
	Tags []string `json:"tags"`
	ID   string   `json:"id"`
}

type Config struct {
	ApiKey         string          `yaml:"apiKey"`
	ZoneName       string          `yaml:"zoneName"`
	RecordSettings []RecordSetting `yaml:"recordNames"`
	RecordType     string          `yaml:"recordType"`
	IpUrl          string          `yaml:"ipUrl"`
}

// Top-level response
type CloudFlareDnsRecordWriteResponse struct {
	Errors   []APIMessage `json:"errors"`
	Messages []APIMessage `json:"messages"`
	Success  bool         `json:"success"`
	Result   Zone         `json:"result"`
}

// Reused error/message shape
type APIMessage struct {
	Code             int            `json:"code"`
	Message          string         `json:"message"`
	DocumentationURL string         `json:"documentation_url"`
	Source           *MessageSource `json:"source,omitempty"`
}

type MessageSource struct {
	Pointer string `json:"pointer"`
}

// "result" object
type Zone struct {
	ID                string     `json:"id"`
	Account           Account    `json:"account"`
	Meta              Meta       `json:"meta"`
	Name              string     `json:"name"`
	Owner             Owner      `json:"owner"`
	Plan              Plan       `json:"plan"`
	CNameSuffix       string     `json:"cname_suffix"`
	Paused            bool       `json:"paused"`
	Permissions       []string   `json:"permissions"`
	Tenant            Tenant     `json:"tenant"`
	TenantUnit        TenantUnit `json:"tenant_unit"`
	Type              string     `json:"type"` // JSON key "type"
	VanityNameServers []string   `json:"vanity_name_servers"`
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Meta struct {
	CDNOnly                bool `json:"cdn_only"`
	CustomCertificateQuota int  `json:"custom_certificate_quota"`
	DNSOnly                bool `json:"dns_only"`
	FoundationDNS          bool `json:"foundation_dns"`
	PageRuleQuota          int  `json:"page_rule_quota"`
	PhishingDetected       bool `json:"phishing_detected"`
	Step                   int  `json:"step"`
}

type Owner struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Plan struct {
	ID                string  `json:"id"`
	CanSubscribe      bool    `json:"can_subscribe"`
	Currency          string  `json:"currency"`
	ExternallyManaged bool    `json:"externally_managed"`
	Frequency         string  `json:"frequency"`
	IsSubscribed      bool    `json:"is_subscribed"`
	LegacyDiscount    bool    `json:"legacy_discount"`
	LegacyID          string  `json:"legacy_id"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
}

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TenantUnit struct {
	ID string `json:"id"`
}

type RecordSetting struct {
	Name    string  `yaml:"name" json:"name"`
	Type    string  `yaml:"type" json:"type"`
	Ttl     int     `yaml:"ttl" json:"ttl"`
	Proxied *bool   `yaml:"proxied" json:"proxied"`
	Comment *string `yaml:"comment" json:"comment"`
	Content string  `yaml:"content" json:"content"`
}

type RecordsAndSettings struct {
	Records  map[string]map[string][]CloudFlareDnsRecord
	Settings map[string]map[string][]RecordSetting
}

type PatchRecord struct {
	RecordSetting
	RecordID string `json:"id"`
}

type CloudFlareBatchRecordPayload struct {
	Posts   []RecordSetting `yaml:"posts"`
	Patches []PatchRecord   `yaml:"patches"`
}

func (rs *RecordSetting) SetZeroValues() {
	if rs.Type == "" {
		rs.Type = "A"
	}
	if rs.Proxied == nil {
		rs.Proxied = BoolPtr(true)
	}
	if rs.Content == "" {
		rs.Content = ipAddress
	}
}
