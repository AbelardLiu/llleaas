FROM alpine:3.8

USER root

ENV USER_UID=1001 \
    USER_NAME=llleaas \
    GROUP_ID=1001 \
    GROUP_NAME=llleaas \
    HOME=/home/llleaas \
    LLLEAAS=/home/llleaas/bin/llleaas

RUN mkdir -p ${HOME}/log; \
    mkdir -p ${HOME}/bin

COPY _output/bin/llleaas ${LLLEAAS}

COPY bin ${HOME}/bin
COPY config ${HOME}/config

RUN ${HOME}/bin/user_setup

RUN chown -R ${USER_UID}:${GROUP_ID} ${HOME}

USER ${USER_UID}

ENTRYPOINT ["/home/llleaas/bin/entrypoint"]