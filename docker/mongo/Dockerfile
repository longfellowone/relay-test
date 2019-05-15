FROM mongo:latest

COPY /docker/mongo/initdb.js /docker-entrypoint-initdb.d

EXPOSE 27017

#RUN apt-get update && \
#    apt-get install gettext -y
#
#COPY /docker/envoy/envoy.yaml /tmpl/envoy.yaml.tmpl
#COPY /docker/envoy/entrypoint.sh /
#
#RUN chmod 500 /entrypoint.sh

#EXPOSE 9090
#
#ENTRYPOINT ["/entrypoint.sh"]