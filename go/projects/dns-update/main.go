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
	stdoutEncoder  *json.Encoder
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
	stdoutEncoder.Encode(recordResult)
	return recordResult
}

func GetDNSRecords(apiKey string, apiUrl string, zoneID string) *CloudFlareDnsRecordResponse {
	rawURL := fmt.Sprintf("%s/client/v4/zones/%s/dns_records", apiUrl, zoneID)
	fmt.Printf("Record Query URL: %s\n", rawURL)
	client := &http.Client{}
	body := &bytes.Buffer{}
	request, _ := http.NewRequest("GET", rawURL, body)
	request.Header.Add("Authorization", "Bearer "+apiKey)
	response, _ := client.Do(request)
	data, _ := io.ReadAll(response.Body)
	stdoutEncoder.Encode(data)
	recordResult := &CloudFlareDnsRecordResponse{}
	json.Unmarshal(data, recordResult)
	// stdoutEncoder.Encode(recordResult)
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

func ExtractSettingsFromRecords(settings []RecordSetting, existingRecords []CloudFlareDnsRecord, defaultIP string) (payload CloudFlareBatchRecordPayload) {
	recordsByName := make(map[string]map[string][]*CloudFlareDnsRecord)
	for _, record := range existingRecords {
		if recordsByName[record.Name] == nil {
			recordsByName[record.Name] = make(map[string][]*CloudFlareDnsRecord)
		}
		recordsByName[record.Name][record.Type] = append(recordsByName[record.Name][record.Type], &record)
	}
	stdoutEncoder.Encode(recordsByName)
	for _, setting := range settings {
		if recordForSetting := recordsByName[setting.Name][setting.Type]; len(recordForSetting) != 0 {
			switch setting.State {
			case "present":
				content := setting.Content
				if content == "" {
					content = recordsByName[setting.Name][setting.Type][0].Content
				}
				payload.Patches = append(payload.Patches,
					PatchRecord{
						RecordSetting: RecordSetting{
							Name:     setting.Name,
							Ttl:      setting.Ttl,
							Type:     setting.Type,
							Proxied:  setting.Proxied,
							Comment:  setting.Comment,
							Content:  content,
							Priority: setting.Priority,
						},
						RecordID: recordForSetting[0].ID,
					})
			case "absent":
				payload.Deletes = append(payload.Deletes,
					DeleteRecord{
						ID: recordForSetting[0].ID,
					})
			}
		} else {
			content := setting.Content
			if content == "" {
				content = defaultIP
			}
			payload.Posts = append(payload.Posts,
				RecordSetting{
					Name:     setting.Name,
					Ttl:      setting.Ttl,
					Type:     setting.Type,
					Proxied:  setting.Proxied,
					Comment:  setting.Comment,
					Content:  content,
					Priority: setting.Priority,
				})
		}
	}
	return payload
}

func SetDnsRecordFromResponse(existingRecords *CloudFlareDnsRecordResponse, recordSettings []RecordSetting, apiUrl string, apiKey string, zoneID string, defaultIP string) []*http.Response {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	recordUpdates := ExtractSettingsFromRecords(recordSettings, existingRecords.Result, defaultIP)
	body := &bytes.Buffer{}
	payload, _ := json.Marshal(recordUpdates)
	body.Write(payload)
	fmt.Println("Record Updates Payload: ")
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
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")
	bodyEncoder := json.NewEncoder(body)
	client := http.Client{}
	responses := []*http.Response{}
	for _, patch := range records.Patches {
		body.Reset()
		url := fmt.Sprintf("%s/client/v4/zones/%s/dns_records/%s", apiUrl, zoneID, patch.RecordID)
		bodyEncoder.Encode(patch)
		request, _ := http.NewRequest("PATCH", url, body)
		request.Header.Add("Authorization", "Bearer "+apiKey)
		response, _ := client.Do(request)
		responses = append(responses, response)
		fmt.Printf("PATCH: %s STATUS CODE: %d\n", patch.Name, response.StatusCode)
		// encoder.Encode(patch)
		// body, _ := io.ReadAll(response.Body)
		// fmt.Println(string(body))
	}
	for _, post := range records.Posts {
		body.Reset()
		url := fmt.Sprintf("%s/client/v4/zones/%s/dns_records", apiUrl, zoneID)
		bodyEncoder.Encode(post)
		request, _ := http.NewRequest("POST", url, body)
		request.Header.Add("Authorization", "Bearer "+apiKey)
		response, _ := client.Do(request)
		responses = append(responses, response)
		fmt.Printf("POST: %s STATUS CODE: %d\n", post.Name, response.StatusCode)
		// encoder.Encode(post)
		// body, _ := io.ReadAll(response.Body)
		// fmt.Println(string(body))
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
	endpointURL := fmt.Sprintf("%s/client/v4/zones/%s/dns_records/batch", apiUrl, zoneID)
	fmt.Printf("Batch record endpiont URL: %s\n", endpointURL)
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
	responseBody := &CloudFlareDnsBatchRecordResponse{}
	json.Unmarshal(responseContent, &responseBody)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(responseBody)
	// fmt.Println("returning")
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
		SetConfigDefaults(&config)
	}
	return config
}

func SetConfigDefaults(config *Config) {
	if config.ApiKey == "" {
		config.ApiKey = GetAPIKey()
	}
	if config.IpUrl == "" {
		config.IpUrl = "https://api.ipify.org"
	}
	if config.DefaultIP == "" {
		config.DefaultIP = GetIpAddress(config.IpUrl)
	}
}

func CreateDefaultConfig(apiKey string) Config {
	config := Config{
		ApiKey:   apiKey,
		ZoneName: "harmonlab.io",
		RecordSettings: []RecordSetting{
			{
				Name:    "harmonlab.io",
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

func AppendBaseDomain(rs *RecordSetting, baseDomain string) {
	if !strings.HasSuffix(rs.Name, baseDomain) {
		rs.Name = rs.Name + baseDomain
	}
}

// func GetConfigWithDefaults("config.yaml", api
func main() {
	stdoutEncoder = json.NewEncoder(os.Stdout)
	stdoutEncoder.SetIndent("", "  ")
	stdoutEncoder.SetEscapeHTML(false)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	config := GetConfig("config.yaml", apiKey)
	stdoutEncoder.Encode(config)
	// fmt.Println("api key: " + config.ApiKey)
	ipAddress := GetIpAddress(config.IpUrl)
	encoder.Encode(ipAddress)
	dnsZone := GetCloudFlareDnsZone(config.ApiKey, apiUrl, config.ZoneName)
	// encoder.Encode(dnsZone)
	dnsRecords := GetDNSRecords(config.ApiKey,
		apiUrl,
		dnsZone.Result[0].ID)
	fmt.Println("DNS record result:")
	stdoutEncoder.Encode(dnsRecords)
	for idx := range config.RecordSettings {
		(&config.RecordSettings[idx]).SetDefaultValues(config.ZoneName)
		// AppendBaseDomain(&rs, config.ZoneName)
	}
	fmt.Println("Record Settings After Defaults: ")
	stdoutEncoder.Encode(config.RecordSettings)
	SetDnsRecordFromResponse(dnsRecords, config.RecordSettings, apiUrl, config.ApiKey, dnsZone.Result[0].ID,
		config.DefaultIP)
	// fmt.Println("Set response: ")
	// bodyBuffer := &bytes.Buffer{}
	//
	//	for _, response := range setResponse {
	//		bodyBuffer.Reset()
	//		tmp, _ := io.ReadAll(response.Body)
	//		bodyBuffer.Write(tmp)
	//		encoder.Encode(bodyBuffer.Bytes())
	//	}
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
	} `json:"settings"`
	Tags []string `json:"tags"`
	ID   string   `json:"id"`
}

type Config struct {
	ApiKey         string          `yaml:"apiKey" json:"apiKey"`
	ZoneName       string          `yaml:"zoneName" json:"zoneName"`
	RecordSettings []RecordSetting `yaml:"recordSettings" json:"recordSettings"`
	RecordType     string          `yaml:"recordType" json:"recordType"`
	IpUrl          string          `yaml:"ipUrl" json:"ipUrl"`
	DefaultIP      string          `yaml:"defaultIP" json:"defaultIP"`
}

// Top-level response
type CloudFlareDnsRecordWriteResponse struct {
	Errors   []APIMessage `json:"errors"`
	Messages []APIMessage `json:"messages"`
	Success  bool         `json:"success"`
	Result   Zone         `json:"result"`
}

type CloudFlareDnsBatchRecordResponse struct {
	Errors   []APIMessage `json:"errors"`
	Messages []APIMessage `json:"messages"`
	Success  bool         `json:"success"`
	Result   struct {
		Deletes []CloudFlareDnsRecord `json:"deletes"`
		Posts   []CloudFlareDnsRecord `json:"posts"`
		Patches []CloudFlareDnsRecord `json:"patches"`
		Puts    []CloudFlareDnsRecord `json:"puts"`
	} `json:"result"`
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
	Name     string  `yaml:"name" json:"name"`
	Type     string  `yaml:"type" json:"type"`
	Ttl      int     `yaml:"ttl" json:"ttl"`
	Proxied  *bool   `yaml:"proxied" json:"proxied"`
	Comment  *string `yaml:"comment" json:"comment"`
	Content  string  `yaml:"content" json:"content"`
	Priority *int    `yaml:"priority" json:"priority,omitempty"`
	State    string  `yaml:"state"`
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
	Posts   []RecordSetting `yaml:"posts" json:"posts"`
	Patches []PatchRecord   `yaml:"patches" json:"patches"`
	Deletes []DeleteRecord  `yaml:"deletes" json:"deletes"`
}

type DeleteRecord struct {
	ID string `yaml:"id" json:"id"`
}

func (rs *RecordSetting) SetDefaultValues(baseDomain string) {
	AppendBaseDomain(rs, baseDomain)
	if rs.Type == "" {
		rs.Type = "A"
	}
	if rs.Proxied == nil {
		rs.Proxied = BoolPtr(true)
	}
	if rs.Content == "" {
		rs.Content = ipAddress
	}
	if rs.State == "" {
		fmt.Println("Setting record setting state")
		rs.State = "present"
	} else if rs.State != "present" && rs.State != "absent" {
		fmt.Println("Setting record setting state")
		rs.State = "absent"
	}
}
