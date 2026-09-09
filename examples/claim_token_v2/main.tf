terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

# defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Create a Foundation Central claim token used to register nodes.
resource "nutanix_claim_token_v2" "example" {
  name            = var.claim_token_name
  expiry_time     = var.claim_token_expiry_time
  max_usage_count = var.claim_token_max_usage_count
}

# Read a single claim token by its external id.
data "nutanix_claim_token_v2" "get-claim-token" {
  ext_id = nutanix_claim_token_v2.example.id
}

# Read the secret value of the claim token.
data "nutanix_secret_v2" "get-claim-token-secret" {
  ext_id = nutanix_claim_token_v2.example.id
}

# List all claim tokens, filtered by name.
data "nutanix_claim_tokens_v2" "list-claim-tokens" {
  filter = "name eq '${nutanix_claim_token_v2.example.name}'"
}

# Read a single node registered with the claim token (by node ext_id).
data "nutanix_node_v2" "get-node" {
  claim_token_ext_id = nutanix_claim_token_v2.example.id
  ext_id             = var.node_ext_id
}

# List all nodes registered with the claim token.
data "nutanix_nodes_v2" "list-nodes" {
  claim_token_ext_id = nutanix_claim_token_v2.example.id
}
