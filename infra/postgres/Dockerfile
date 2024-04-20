FROM postgres:15
RUN apt-get update && apt-get -y install postgresql-15-cron
COPY ./initdb.d/ /docker-entrypoint-initdb.d/
