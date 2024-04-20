FROM postgres:16
RUN apt-get update && apt-get -y install postgresql-16-cron
COPY ./initdb.d/ /docker-entrypoint-initdb.d/
