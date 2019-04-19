package ram

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"

	"github.com/denverdino/aliyungo/metadata"
)

const fullArnPrefix = "acs:"

// ARNRegexp is the regex to check that the base ARN is valid,
// see https://help.aliyun.com/document_detail/28628.html.
var ARNRegexp = regexp.MustCompile(`^acs:ram::\d+:role\/?(\w+|-|\/|\.)*$`)

// IsValidBaseARN validates that the base ARN is valid.
func IsValidBaseARN(arn string) bool {
	return ARNRegexp.MatchString(arn)
}

// RoleARN returns the full ram role ARN.
func (ram *Client) RoleARN(role string) string {
	log.Infof("RoleARN, ram.BaseARN is %++v", ram.BaseARN)

	if strings.HasPrefix(strings.ToLower(role), fullArnPrefix) {
		return role
	}
	//return if the wantedRole is the ecs ram role
	if strings.HasSuffix(ram.BaseARN, strings.ToLower(role)) {
		return role
	}

	prefixRole := strings.Split(ram.BaseARN, "/")[0]
	return fmt.Sprintf("%s/%s", prefixRole, role)
}

// GetBaseArn get the base ARN from metadata service.
func GetBaseArn() (string, error) {
	metadata := metadata.NewMetaData(nil)
	ownerAccountID, err := metadata.OwnerAccountID()
	if err != nil {
		return "", err
	}

	roleName, err := metadata.RoleName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("acs:ram::%s:role/%s", ownerAccountID, strings.ToLower(roleName)), nil
}
