zitadel_domain   = "zitadel.test"
app_domain       = "tasks-app.test"
initial_password = "S3c_r3t!"
# zitadel_domain   = "zitadel.${DOMAIN}"
# app_domain       = "tasks-app.${DOMAIN}"

app_users = [
  {
    user_name  = "johndoe"
    first_name = "John"
    last_name  = "Doe"
    email      = "john.doe@tasks-app.test"
    # email      = "john.doe@tasks-app.${DOMAIN}"
  },
  {
    user_name  = "janedoe"
    first_name = "Jane"
    last_name  = "Doe"
    email      = "jane.doe@tasks-app.test"
    # email      = "jane.doe@tasks-app.${DOMAIN}"
  }
]
