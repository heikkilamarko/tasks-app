FROM ghcr.io/zitadel/zitadel:v2.49.0
COPY ./zitadel.yaml /zitadel.yaml
COPY ./steps.yaml /steps.yaml
CMD ["start-from-init", "--masterkeyFromEnv", "--config", "/zitadel.yaml", "--steps", "/steps.yaml"]
