provider "ibm" {
  ibmcloud_api_key = var.ibmcloud_api_key
}

// Provision iam_trusted_profiles_claim_rule resource instance
resource "ibm_iam_trusted_profiles_claim_rule" "iam_trusted_profiles_claim_rule_instance" {
  profile_id = "profile_id"
  type = "type"
  name = "name"
  realm_name = "realm_name"
  expiration = 43200
  conditions {
				claim = "claim"
				operator = "operator"
				value = "value"
			}
}

// Create iam_trusted_profiles_claim_rule data source
data "ibm_iam_trusted_profiles_claim_rule" "iam_trusted_profiles_claim_rule_instance" {
  profile_id = var.iam_trusted_profiles_claim_rule_profile_id
  rule_id = var.iam_trusted_profiles_claim_rule_rule_id
}
