package quota

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github/yasun1/myquota/pkg/connection"
	"github/yasun1/myquota/pkg/constants/http"
	AMS "github/yasun1/myquota/pkg/endpoints/ams"
	. "github/yasun1/myquota/pkg/helpers"
	"github/yasun1/myquota/pkg/logs/debug"

	client "github.com/openshift-online/ocm-sdk-go"
)

const (
	skuRBTemplate = `{
		"sku": "%s",
		"sku_count": %d,
		"type": "%s"
	  }	
	`
	manualQuotaType = "Manual"
)

var SuperAdminConnection = connection.SuperAdminConnection

type Sku struct {
	Name     string
	QuotaID  string
	Type     string
	Allowed  int
	Consumed int
}

// var SkuMap = allSkus()

func AllSkus() map[string]Sku {
	params := map[string]interface{}{
		"size": 10000,
	}
	resp, err := AMS.ListSkuRules(SuperAdminConnection, params)
	if err != nil || resp.Status() != http.HTTPOK {
		panic(fmt.Errorf("[E] Failed to List skus:\n%v\n%v", err, resp.String()))
	}

	skuMap := make(map[string]Sku)
	skuRuleItems := DigArray(Parse(resp.Bytes()), "items")
	for _, skuRule := range skuRuleItems {
		skuName := DigString(skuRule, "sku")
		quotaID := DigString(skuRule, "quota_id")

		sku := Sku{
			Name:    skuName,
			QuotaID: quotaID,
		}
		skuMap[skuName] = sku
	}

	if len(skuMap) == 0 {
		panic(fmt.Errorf("[E] No valid skus in OCM\n"))
	}

	if debug.DebugMode() {
		for _, sku := range skuMap {
			fmt.Printf("[Debug] %s->%s\n", sku.Name, sku.QuotaID)
		}
	}

	return skuMap
}

// GetOrgID retturns the ocm organization id of the user
func GetOrgID(username string) (string, error) {
	params := map[string]interface{}{
		"search": fmt.Sprintf("username is '%s'", username),
	}
	resp, err := AMS.ListAccounts(SuperAdminConnection, params)
	if err != nil || resp.Status() != http.HTTPOK {
		err = fmt.Errorf("[E] Failed to List account\n%v\n%v\n", err, resp.String())
		return "", err
	}

	accountItems := DigArray(Parse(resp.Bytes()), "items")
	if len(accountItems) != 1 {
		err = fmt.Errorf("[E] Expect 1 but find %d for the account '%s'\n", len(accountItems), username)
		return "", err
	}

	organizationID := DigString(accountItems[0], "organization", "id")
	if organizationID == "" {
		err = fmt.Errorf("[E] The orgnization id is empty for the account '%s'\n", username)
	}

	return organizationID, err
}

// IsAssigned will check whether the quota is assigned
func IsAssigned(connection *client.Connection, orgID string, sku Sku) (string, bool, error) {
	params := map[string]interface{}{
		"search": fmt.Sprintf("sku is '%s' and type is '%s'", sku.Name, sku.Type),
	}
	resp, err := AMS.ListOrgResourceQuotas(connection, orgID, params)
	if err != nil || resp.Status() != http.HTTPOK {
		err = fmt.Errorf("[E] Failed to List resource quota\n%v\n%v\n", err, resp.String())
		return "", false, err
	}

	quotaItems := DigArray(Parse(resp.Bytes()), "items")
	if len(quotaItems) == 0 {
		return "", false, nil
	}

	resourceQuotaID := DigString(quotaItems[0], "id")
	return resourceQuotaID, true, nil
}

// OrgQuotas get the assigned resource quota in the organization
func OrgQuotas(orgID string) (map[string]string, error) {
	skuMap := AllSkus()

	params := map[string]interface{}{
		"size": 10000,
	}
	resp, err := AMS.ListOrgResourceQuotas(SuperAdminConnection, orgID, params)
	if err != nil || resp.Status() != http.HTTPOK {
		err = fmt.Errorf("[E] Failed to List resource quota\n")
		return nil, err
	}

	quotaMap := make(map[string]string)
	quotaItems := DigArray(Parse(resp.Bytes()), "items")
	for _, quota := range quotaItems {
		skuName := DigString(quota, "sku")
		sku := skuMap[skuName]

		if _, existed := quotaMap[sku.QuotaID]; existed {
			skuName = fmt.Sprintf("%s,%s", quotaMap[sku.QuotaID], skuName)
			quotaMap[sku.QuotaID] = skuName
			continue
		}

		quotaMap[sku.QuotaID] = skuName
	}

	return quotaMap, err
}

// AssignQuota assigns the quota to the organization.
// If the resource quota exists, will update its allowed to the new value.
// If the resource quota does not exist, will create a new resource quota.
func AssignQuota(orgID string, sku Sku) (string, error) {
	resourceQuotaID, existed, err := IsAssigned(SuperAdminConnection, orgID, sku)
	if err != nil {
		return "", err
	}

	var resp *client.Response
	quotaRB := fmt.Sprintf(skuRBTemplate, sku.Name, sku.Allowed, sku.Type)
	if existed {
		resp, err = AMS.PatchOrgResourceQuotaByID(SuperAdminConnection, orgID, resourceQuotaID, quotaRB)
	} else {
		resp, err = AMS.CreateOrgResourceQuota(SuperAdminConnection, orgID, quotaRB)
	}

	if err != nil || (resp.Status() != http.HTTPOK && resp.Status() != http.HTTPCreated) {
		err = fmt.Errorf("[E] Failed to assign %d %s_%s resource quota to the organization %s:\n%v\n%v",
			sku.Allowed, sku.Name, sku.Type, orgID, err, resp)
		return "", err
	}

	fmt.Printf("Successfully assign %d %s_%s resource quota to the organization %s\n", sku.Allowed, sku.Name, sku.Type, orgID)
	resourceQuotaID = DigString(Parse(resp.Bytes()), "id")
	return resourceQuotaID, err
}

// FPrintQuotaCost prints all the resource quotas in the organization.
func FPrintQuotaCost(orgID string) {
	quotaMap, err := OrgQuotas(orgID)
	if err != nil {
		panic(err)
	}

	params := map[string]interface{}{
		"size": 10000,
	}
	resp, err := AMS.RetrieveQuotaCost(SuperAdminConnection, orgID, params)
	if err != nil || resp.Status() != http.HTTPOK {
		err = fmt.Errorf("[E] Failed to get the quota cost of the organization %s\n", orgID)
		panic(err)
	}

	fmt.Printf("\n>>> The quota under the organization %s: \n", orgID)
	quotaCostItems := DigArray(Parse(resp.Bytes()), "items")
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	fmt.Fprintf(writer, "Name\tQuotaID\tAllowed\tConsumed\t\n")
	for _, quotaCost := range quotaCostItems {
		quotaID := DigString(quotaCost, "quota_id")
		allowed := DigInt(quotaCost, "allowed")
		consumed := DigInt(quotaCost, "consumed")
		skuNames := quotaMap[quotaID]

		fmt.Fprintf(writer, "%s\t%s\t%d\t%d\n",
			skuNames,
			quotaID,
			allowed,
			consumed,
		)
	}
	writer.Flush()
}

// getUsageForQuota gets the usage of the specified quota.
func getUsageForQuota(orgID string, sku Sku) Sku {
	params := map[string]interface{}{
		"search": fmt.Sprintf("quota_id is '%s'", sku.QuotaID),
	}

	resp, err := AMS.RetrieveQuotaCost(SuperAdminConnection, orgID, params)
	if err != nil || resp.Status() != http.HTTPOK {
		err = fmt.Errorf("[E] Failed to get the quota cost of the organization %s\n", orgID)
		panic(err)
	}

	quotaCostItems := DigArray(Parse(resp.Bytes()), "items")
	if len(quotaCostItems) == 1 {
		sku.Allowed = DigInt(quotaCostItems[0], "allowed")
		sku.Consumed = DigInt(quotaCostItems[0], "consumed")
	}

	return sku
}

// FPrintUsageForSkus prints the usage of the specified resource quotas.
func FPrintUsageForSkus(orgID string, skus ...Sku) {
	fmt.Printf("\n>>> The quota under the organization %s: \n", orgID)
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	fmt.Fprintf(writer, "Name\tQuotaID\tAllowed\tConsumed\t\n")
	for _, sku := range skus {
		sku = getUsageForQuota(orgID, sku)
		fmt.Fprintf(writer, "%s\t%s\t%d\t%d\n",
			sku.Name,
			sku.QuotaID,
			sku.Allowed,
			sku.Consumed,
		)
	}
	writer.Flush()
}

// RemoveQuota removes the resource quota from the organization.
// If the resource quota is in used, option '--force' is required.
func RemoveQuota(orgID string, sku Sku, force bool) {
	resourceQuotaID, existed, err := IsAssigned(SuperAdminConnection, orgID, sku)
	if err != nil {
		panic(err)
	}
	if !existed {
		fmt.Printf("[W] The resource quota with the sku '%s_%s' is not assigned. Give up removing.\n", sku.Name, sku.Type)
		os.Exit(0)
	}

	sku = getUsageForQuota(orgID, sku)
	if sku.Consumed != 0 && !force {
		FPrintUsageForSkus(orgID, sku)
		err = fmt.Errorf("[W] The resource quota is in used. If you truly remove the quota, please use with the option '--force'\n")
		panic(err)
	}

	resp, err := AMS.DeleteOrgResourceQuotaByID(SuperAdminConnection, orgID, resourceQuotaID)
	if err != nil || resp.Status() != http.HTTPNoContent {
		err = fmt.Errorf("[E] Failed to remove the %s_%s resource quota(%s) from the organization %s:\n%v\n%v\n",
			sku.Name, sku.Type, resourceQuotaID, orgID, err, resp)
		panic(err)
	}

	fmt.Printf("Successfully remove the %s_%s resource quota from the organization %s\n", sku.Name, sku.Type, orgID)
}
