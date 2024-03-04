package ams

// Account Management Service
const (
	// account
	accountURL = "/api/accounts_mgmt/v1/accounts"

	// quota
	skuRuleURL         = "/api/accounts_mgmt/v1/sku_rules"
	skuRuleIDURL       = "/api/accounts_mgmt/v1/sku_rules/%s"
	quotaCostURL       = "/api/accounts_mgmt/v1/organizations/%s/quota_cost"
	resourceQuotaURL   = "/api/accounts_mgmt/v1/organizations/%s/resource_quota"
	resourceQuotaIDURL = "/api/accounts_mgmt/v1/organizations/%s/resource_quota/%s"
)
