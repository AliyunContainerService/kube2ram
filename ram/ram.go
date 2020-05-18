package ram

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	"github.com/AliyunContainerService/kube2ram/metrics"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/denverdino/aliyungo/metadata"
	"github.com/karlseguin/ccache"
	log "github.com/sirupsen/logrus"
)

var cache = ccache.New(ccache.Configure())

const (
	maxSessNameLength = 32
	//STSEndpoint is default sts server address
	STSEndpoint = "https://sts.aliyuncs.com"
)

// Client represents an RAM client.
type Client struct {
	BaseARN         string
	Endpoint        string
	Region          string
	AccessKey       string
	AccessSecret    string
	STSAccessKey    string
	STSAccessSecret string
	STSToken        string
	BaseRoleName    string
}

// Credentials represent the security credentials response
type Credentials struct {
	Code            string
	LastUpdated     string
	AccessKeyID     string `json:"AccessKeyId"`
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
}

func getHash(text string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(text))
	if err != nil {
		return text
	}
	return fmt.Sprintf("%x", h.Sum32())
}

// GetInstanceRAMRole get instance RAM role from metadata service.
func GetInstanceRAMRole() (string, error) {
	metadata := metadata.NewMetaData(nil)

	ramRole, err := metadata.RoleName()
	if err != nil {
		return "", err
	}
	if ramRole == "" || err != nil {
		return "", errors.New("ECS Metadata didn't returned any RAM Role")
	}
	return ramRole, nil
}

func sessionName(roleARN, remoteIP string) string {
	idx := strings.LastIndex(roleARN, "/")
	name := fmt.Sprintf("%s-%s", getHash(remoteIP), roleARN[idx+1:])
	return fmt.Sprintf("%.[2]*[1]s", name, maxSessNameLength)
}

// Helper to format RAM return codes for metric labeling
func getRAMCode(err error, statusCode int) string {
	if err != nil {
		if statusCode != 0 {
			return strconv.Itoa(statusCode)
		}
		return metrics.RAMUnknownFailCode
	}
	return metrics.RAMSuccessCode
}

//AlibabaCloudConfig define metadata struct return from metadata API
type AlibabaCloudConfig struct {
	RegionID        string    `json:"regionId" yaml:"regionId"`
	AccessKeyID     string    `json:"accessKeyId" yaml:"accessKeyId"`
	AccessKeySecret string    `json:"accessKeySecret" yaml:"accessKeySecret"`
	VPCID           string    `json:"vpcId" yaml:"vpcId"`
	RoleName        string    `json:"-" yaml:"-"` // For ECS RAM role only
	StsToken        string    `json:"-" yaml:"-"`
	ExpireTime      time.Time `json:"-" yaml:"-"`
}

// GetCloudConfigFromMetaServer return the credentials and instance config from ECS metadata server
func GetCloudConfigFromMetaServer() (*AlibabaCloudConfig, error) {
	cfg := &AlibabaCloudConfig{}
	// Load config from Metadata Service
	m := metadata.NewMetaData(nil)
	roleName := ""
	var err error
	if roleName, err = m.RoleName(); err != nil {
		return cfg, fmt.Errorf("failed to get role name from Metadata Service: %v", err)
	}
	vpcID, err := m.VpcID()
	if err != nil {
		return cfg, fmt.Errorf("failed to get VPC ID from Metadata Service: %v", err)
	}
	regionID, err := m.Region()
	if err != nil {
		return cfg, fmt.Errorf("failed to get Region ID from Metadata Service: %v", err)
	}
	role, err := m.RamRoleToken(roleName)
	if err != nil {
		return cfg, fmt.Errorf("failed to get STS Token from Metadata Service: %v", err)
	}
	cfg.RegionID = regionID
	cfg.RoleName = roleName
	cfg.VPCID = vpcID
	cfg.AccessKeyID = role.AccessKeyId
	cfg.AccessKeySecret = role.AccessKeySecret
	cfg.StsToken = role.SecurityToken
	cfg.ExpireTime = role.Expiration
	return cfg, nil
}

// AssumeRole returns an RAM role Credentials using Alibaba Cloud STS.
func (ram *Client) AssumeRole(roleARN, remoteIP string, sessionTTL time.Duration) (*Credentials, error) {
	hitCache := true
	item, err := cache.Fetch(roleARN, sessionTTL, func() (interface{}, error) {
		hitCache = false

		// Set up a prometheus timer to track the Alibaba Cloud request duration. It stores the timer value when
		// observed. A function gets err at observation time to report the status of the request after the function returns.
		var err error
		var status int
		lvsProducer := func() []string {
			return []string{getRAMCode(err, status), roleARN}
		}
		timer := metrics.NewFunctionTimer(metrics.RAMRequestSec, lvsProducer, nil)
		defer timer.ObserveDuration()

		var svc *sts.Client
		if ram.AccessKey == "" || ram.AccessSecret == "" {
			log.Printf("init sts client with ram role %s", ram.BaseRoleName)

			svc, err = sts.NewClientWithEcsRamRole(ram.Region, ram.BaseRoleName)
			if err != nil {
				return nil, err
			}
		} else {
			svc, err = sts.NewClientWithAccessKey(ram.Region, ram.AccessKey, ram.AccessSecret)
			if err != nil {
				return nil, err
			}
		}

		log.Infof("assume role arn is %s", roleARN)
		request := sts.CreateAssumeRoleRequest()
		request.Scheme = "https"
		request.RoleArn = roleARN
		request.RoleSessionName = sessionName(roleARN, remoteIP)
		request.DurationSeconds = requests.NewInteger(int(sessionTTL.Seconds() * 2))

		resp, err := svc.AssumeRole(request)
		if resp != nil {
			status = resp.GetHttpStatus()
		}
		if err != nil {
			return nil, err
		}

		return &Credentials{
			AccessKeyID:     resp.Credentials.AccessKeyId,
			Code:            "Success",
			Expiration:      resp.Credentials.Expiration,
			LastUpdated:     time.Now().Format("2006-01-02T15:04:05Z"),
			AccessKeySecret: resp.Credentials.AccessKeySecret,
			SecurityToken:   resp.Credentials.SecurityToken,
		}, nil
	})
	if hitCache {
		metrics.RAMCacheHitCount.WithLabelValues(roleARN).Inc()
	}
	if err != nil {
		return nil, err
	}
	return item.Value().(*Credentials), nil
}

// NewClient returns a new STS client.
func NewClient(acssessKey, accessSecret, baseARN string, config *AlibabaCloudConfig) *Client {

	return &Client{
		Region:          config.RegionID,
		BaseARN:         baseARN,
		Endpoint:        STSEndpoint,
		AccessKey:       acssessKey,
		AccessSecret:    accessSecret,
		STSAccessKey:    config.AccessKeyID,
		STSAccessSecret: config.AccessKeySecret,
		STSToken:        config.StsToken,
		BaseRoleName:    config.RoleName,
	}
}
