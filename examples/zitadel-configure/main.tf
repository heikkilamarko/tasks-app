terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.0.5"
    }
  }
}

provider "zitadel" {
  domain           = "auth.tasks-app.com"
  insecure         = false
  port             = 443
  jwt_profile_file = "./machinekey/zitadel-admin-sa.json"
}

# Organizations

data "zitadel_orgs" "zitadel" {
  name        = "zitadel"
  name_method = "TEXT_QUERY_METHOD_EQUALS"
}

resource "zitadel_org" "tasks_app" {
  name = "tasks-app"
}

# Users

resource "zitadel_human_user" "zitadel_admin" {
  org_id             = data.zitadel_orgs.zitadel.ids[0]
  user_name          = "zitadel-admin"
  first_name         = "Zitadel"
  last_name          = "Admin"
  preferred_language = "en"
  email              = "zitadel-admin@tasks-app.com"
  is_email_verified  = true
  initial_password   = "S3c_r3t!"
}

resource "zitadel_human_user" "editor" {
  org_id             = zitadel_org.tasks_app.id
  user_name          = "editor"
  first_name         = "Editor"
  last_name          = "User"
  preferred_language = "en"
  email              = "editor@tasks-app.com"
  is_email_verified  = true
  initial_password   = "S3c_r3t!"
}

resource "zitadel_human_user" "viewer" {
  org_id             = zitadel_org.tasks_app.id
  user_name          = "viewer"
  first_name         = "Viewer"
  last_name          = "User"
  preferred_language = "en"
  email              = "viewer@tasks-app.com"
  is_email_verified  = true
  initial_password   = "S3c_r3t!"
}

# Instance Members

resource "zitadel_instance_member" "default" {
  user_id = zitadel_human_user.zitadel_admin.id
  roles   = ["IAM_OWNER"]
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

resource "zitadel_project_role" "editor" {
  org_id       = zitadel_org.tasks_app.id
  project_id   = zitadel_project.tasks_app.id
  role_key     = "editor"
  display_name = "Editor"
}

resource "zitadel_project_role" "viewer" {
  org_id       = zitadel_org.tasks_app.id
  project_id   = zitadel_project.tasks_app.id
  role_key     = "viewer"
  display_name = "Viewer"
}

# User Grants

resource "zitadel_user_grant" "editor_editor" {
  org_id     = zitadel_org.tasks_app.id
  project_id = zitadel_project.tasks_app.id
  user_id    = zitadel_human_user.editor.id
  role_keys  = ["editor"]
}

resource "zitadel_user_grant" "viewer_viewer" {
  org_id     = zitadel_org.tasks_app.id
  project_id = zitadel_project.tasks_app.id
  user_id    = zitadel_human_user.viewer.id
  role_keys  = ["viewer"]
}

# Applications

resource "zitadel_application_oidc" "tasks_app" {
  org_id                      = zitadel_org.tasks_app.id
  project_id                  = zitadel_project.tasks_app.id
  name                        = "tasks-app"
  app_type                    = "OIDC_APP_TYPE_USER_AGENT"
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  auth_method_type            = "OIDC_AUTH_METHOD_TYPE_NONE"
  access_token_type           = "OIDC_TOKEN_TYPE_JWT"
  redirect_uris               = ["https://www.tasks-app.com/ui", "https://www.tasks-app.com/ui/auth/callback"]
  post_logout_redirect_uris   = ["https://www.tasks-app.com/", "http://www.tasks-app.com/"]
  access_token_role_assertion = true
  id_token_role_assertion     = true
  id_token_userinfo_assertion = true
  dev_mode                    = false

  depends_on = [
    zitadel_user_grant.editor_editor,
    zitadel_user_grant.viewer_viewer
  ]
}

# Outputs

output "tasks_app_client_id" {
  value     = zitadel_application_oidc.tasks_app.client_id
  sensitive = true
}
