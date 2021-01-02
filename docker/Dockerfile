FROM couchdb
WORKDIR $PWD
COPY ./docker/etc/local.ini /opt/couchdb/etc/local.d
ADD ./docker/cert/*.crt /opt/couchdb/etc/cert/
ADD ./docker/cert/*.key /opt/couchdb/etc/cert/
ADD ./docker/cert/*.csr /opt/couchdb/etc/cert/
ADD ./docker/cert/*.req /opt/couchdb/etc/cert/