variable "GITHUB_TOKEN" {
  type = string
}

terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "6.7.5"
    }
  }
}

provider "github" {
  token = var.GITHUB_TOKEN
  owner = "gidoichi"
}

resource "github_repository" "this" {
  name                        = "ical-converter"
  allow_auto_merge            = true
  allow_merge_commit          = false
  allow_rebase_merge          = false
  delete_branch_on_merge      = true
  description                 = "Convert iCalendar to register a calendar"
  has_issues                  = true
  squash_merge_commit_message = "BLANK"
  squash_merge_commit_title   = "PR_TITLE"
  pages {
    build_type = "workflow"
  }
}

resource "github_branch_protection" "default" {
  repository_id = github_repository.this.node_id
  pattern       = "main"
  required_status_checks {
    strict = true
    contexts = [
      "build-container",
      "go-test",
      "no-diff",
      "pull-request",
      "terraform-plan",
    ]
  }
}

resource "github_repository_environment" "github-pages" {
  environment = "github-pages"
  repository  = github_repository.this.name
}

resource "github_repository_environment" "dockerhub" {
  environment = "dockerhub"
  repository  = github_repository.this.name
}
