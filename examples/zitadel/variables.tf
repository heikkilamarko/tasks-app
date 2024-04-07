variable "auth_domain" {
  description = "Auth domain"
  type        = string
}

variable "iam_admin_users" {
  description = "The list of iam admin users to create"
  type = list(object({
    user_name  = string
    first_name = string
    last_name  = string
    email      = string
  }))
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

variable "initial_password" {
  description = "The initial password set for the user"
  type        = string
  sensitive   = true
}
