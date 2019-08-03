ARG ALPINE_VERSION=3.10

FROM alpine:${ALPINE_VERSION}
ARG BUILD_DATE
ARG VCS_REF
ARG NAME=Automated
ARG EMAIL=quentin.mcgaw@gmail.com
LABEL org.label-schema.schema-version="1.0.0-rc1" \
      maintainer="quentin.mcgaw@gmail.com" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/qdm12/updated" \
      org.label-schema.url="https://github.com/qdm12/updated" \
      org.label-schema.vcs-description="Docker container to update and push files to Github" \
      org.label-schema.vcs-usage="https://github.com/qdm12/updated/blob/master/README.md#setup" \
      org.label-schema.docker.cmd="docker run -d qmcgaw/updated" \
      org.label-schema.docker.cmd.devel="docker run -it --rm qmcgaw/updated" \
      org.label-schema.docker.params="" \
      org.label-schema.version="" \
      image-size="64.4MB" \
      ram-usage="1MB" \
      cpu-usage=""
RUN apk --update --no-cache --progress -q add openssh-client git ca-certificates wget sed grep bind-tools perl-xml-xpath && \
    rm -rf /var/cache/apk/*
COPY key /home/user/.ssh/id_rsa
RUN adduser -D user --uid 1000 && \
    printf "[user]\nname = ${NAME}\nemail = ${EMAIL}" > /home/user/.gitconfig && \
    touch /home/user/.ssh/known_hosts && \
    chown -R user /home/user/.ssh && \
    chmod 700 /home/user/.ssh && \
    chmod 600 /home/user/.ssh/known_hosts && \
    chmod 400 /home/user/.ssh/id_rsa
ENV VERBOSE=1
COPY entrypoint.sh /
RUN chown user /entrypoint.sh && \
    chmod 500 /entrypoint.sh && \
    mkdir /updated && \
    chown -R user /updated && \
    chmod -R 700 /updated
USER user
ENTRYPOINT ["/entrypoint.sh"]
