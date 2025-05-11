variable "zitadel_domain" {
  description = "Zitadel domain"
  type        = string
}

variable "app_domain" {
  description = "App domain"
  type        = string
}

variable "initial_password" {
  description = "The initial password set for the user"
  type        = string
  sensitive   = true
}

variable "app_users" {
  description = "The list of app users to create"
  type = list(object({
    user_name  = string
    first_name = string
    last_name  = string
    email      = string
  }))
}
