output "tasks_app_client_id" {
  value     = zitadel_application_oidc.tasks_app.client_id
  sensitive = true
}

output "email_notifier_token" {
  value     = zitadel_personal_access_token.email_notifier.token
  sensitive = true
}
