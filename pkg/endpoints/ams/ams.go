package ams

import (
	"fmt"

	client "github.com/openshift-online/ocm-sdk-go"
)

//*************************** AMS *************************
const allowedNum = 1

func parameterError(inputNum int) error {
	if inputNum != allowedNum {
		return fmt.Errorf("Too many parameters, only %d allowed.", allowedNum)
	}
	return nil
}
func parameters(request *client.Request, parameters ...map[string]interface{}) *client.Request {
	for _, params := range parameters {
		for key, value := range params {
			request = request.Parameter(key, value)
		}
	}

	return request
}

// Account
func ListAccounts(connection *client.Connection, params ...map[string]interface{}) (resp *client.Response, err error) {
	if len(params) > 1 {
		err = parameterError(len(params))
		return
	}

	request := connection.Get().Path(accountURL)
	request = parameters(request, params...)
	return request.Send()
}

// Quota

func ListSkuRules(connection *client.Connection, params ...map[string]interface{}) (resp *client.Response, err error) {
	if len(params) > 1 {
		err = parameterError(len(params))
		return nil, err
	}
	request := connection.Get().Path(skuRuleURL)
	if len(params) != 0 {
		for key, value := range params[0] {
			request = request.Parameter(key, value)
		}
	}
	return request.Send()
}

func RetrieveSkuRuleByID(connection *client.Connection, skuRuleID string) (resp *client.Response, err error) {
	resp, err = connection.Get().Path(fmt.Sprintf(skuRuleIDURL, skuRuleID)).Send()
	return
}

func RetrieveQuotaCost(connection *client.Connection, organizationID string, params ...map[string]interface{}) (resp *client.Response, err error) {
	if len(params) > 1 {
		err = parameterError(len(params))
		return nil, err
	}

	request := connection.Get().Path(fmt.Sprintf(quotaCostURL, organizationID))
	request = parameters(request, params...)
	return request.Send()
}

func ListOrgResourceQuotas(connection *client.Connection, organizationID string, params ...map[string]interface{}) (resp *client.Response, err error) {
	if len(params) > 1 {
		err = parameterError(len(params))
		return
	}

	request := connection.Get().Path(fmt.Sprintf(resourceQuotaURL, organizationID))
	request = parameters(request, params...)
	return request.Send()
}

func CreateOrgResourceQuota(connection *client.Connection, organizationID string, body string) (resp *client.Response, err error) {
	resp, err = connection.Post().Path(fmt.Sprintf(resourceQuotaURL, organizationID)).String(body).Send()
	return
}

func RetrieveOrgResourceQuotaByID(connection *client.Connection, organizationID string, quotaID string) (resp *client.Response, err error) {
	resp, err = connection.Get().Path(fmt.Sprintf(resourceQuotaIDURL, organizationID, quotaID)).Send()
	return
}

func PatchOrgResourceQuotaByID(connection *client.Connection, organizationID string, quotaID string, body string) (resp *client.Response, err error) {
	resp, err = connection.Patch().Path(fmt.Sprintf(resourceQuotaIDURL, organizationID, quotaID)).String(body).Send()
	return
}

func DeleteOrgResourceQuotaByID(connection *client.Connection, organizationID string, quotaID string) (resp *client.Response, err error) {
	resp, err = connection.Delete().Path(fmt.Sprintf(resourceQuotaIDURL, organizationID, quotaID)).Send()
	return
}
