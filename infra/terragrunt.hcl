terraform {
  source = ".///"
}

include "backend" {
  path = "${get_repo_root()}/terragrunt_backend.hcl"
}


inputs =   {
    project="recommendation-engine",
    package_name="main"
  }