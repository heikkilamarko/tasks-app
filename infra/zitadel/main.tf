terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "~> 2.0"
    }
  }
}

provider "zitadel" {
  domain           = var.zitadel_domain
  insecure         = false
  port             = 443
  jwt_profile_file = "zitadel-admin-sa.json"
}

# Organization

resource "zitadel_org" "tasks_app" {
  name       = "tasks-app"
  is_default = true
}

# Machine Users

resource "zitadel_machine_user" "email_notifier" {
  org_id            = zitadel_org.tasks_app.id
  user_name         = "email_notifier"
  name              = "email_notifier"
  access_token_type = "ACCESS_TOKEN_TYPE_BEARER"
}

resource "zitadel_personal_access_token" "email_notifier" {
  org_id          = zitadel_org.tasks_app.id
  user_id         = zitadel_machine_user.email_notifier.id
  expiration_date = "9999-12-31T23:59:59Z"
}

resource "zitadel_org_member" "email_notifier" {
  org_id  = zitadel_org.tasks_app.id
  user_id = zitadel_machine_user.email_notifier.id
  roles   = ["ORG_USER_MANAGER"]
}

# Human Users

resource "zitadel_human_user" "app_user" {
  for_each = tomap({ for user in var.app_users : user.user_name => user })

  org_id             = zitadel_org.tasks_app.id
  user_name          = each.value.user_name
  first_name         = each.value.first_name
  last_name          = each.value.last_name
  email              = each.value.email
  initial_password   = var.initial_password
  is_email_verified  = true
  preferred_language = "en"
}

# Projects

resource "zitadel_project" "tasks_app" {
  name                   = "tasks-app"
  org_id                 = zitadel_org.tasks_app.id
  project_role_assertion = true
  project_role_check     = false
  has_project_check      = true
}

# Roles

resource "zitadel_project_role" "user" {
  org_id       = zitadel_org.tasks_app.id
  project_id   = zitadel_project.tasks_app.id
  role_key     = "user"
  display_name = "User"
}

# User Grants

resource "zitadel_user_grant" "app_user" {
  for_each = zitadel_human_user.app_user

  org_id     = zitadel_org.tasks_app.id
  project_id = zitadel_project.tasks_app.id
  user_id    = each.value.id
  role_keys  = ["user"]
}

# Applications

resource "zitadel_application_oidc" "tasks_app" {
  org_id                      = zitadel_org.tasks_app.id
  project_id                  = zitadel_project.tasks_app.id
  name                        = "tasks-app"
  app_type                    = "OIDC_APP_TYPE_WEB"
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  auth_method_type            = "OIDC_AUTH_METHOD_TYPE_NONE"
  access_token_type           = "OIDC_TOKEN_TYPE_JWT"
  redirect_uris               = ["https://${var.app_domain}/ui", "https://${var.app_domain}/ui/auth/callback"]
  post_logout_redirect_uris   = ["https://${var.app_domain}/"]
  access_token_role_assertion = true
  id_token_role_assertion     = true
  id_token_userinfo_assertion = true
  dev_mode                    = false

  depends_on = [
    zitadel_user_grant.app_user
  ]
}
