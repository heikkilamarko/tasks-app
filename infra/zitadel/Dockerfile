FROM ghcr.io/zitadel/zitadel
COPY ./zitadel.yaml /zitadel.yaml
COPY ./steps.yaml /steps.yaml
CMD ["start-from-init", "--masterkeyFromEnv", "--config", "/zitadel.yaml", "--steps", "/steps.yaml"]
