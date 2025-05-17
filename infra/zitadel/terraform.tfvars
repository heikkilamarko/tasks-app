zitadel_domain   = "zitadel.local"
app_domain       = "tasks-app.local"
initial_password = "S3c_r3t!"
# zitadel_domain   = "zitadel.${DOMAIN}"
# app_domain       = "tasks-app.${DOMAIN}"

app_users = [
  {
    user_name  = "johndoe"
    first_name = "John"
    last_name  = "Doe"
    email      = "john.doe@tasks-app.local"
    # email      = "john.doe@tasks-app.${DOMAIN}"
  },
  {
    user_name  = "janedoe"
    first_name = "Jane"
    last_name  = "Doe"
    email      = "jane.doe@tasks-app.local"
    # email      = "jane.doe@tasks-app.${DOMAIN}"
  }
]
