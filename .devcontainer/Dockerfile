ARG VARIANT=1-bullseye
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT} as base
ARG HOME=/home/vscode
RUN chown -R vscode /usr/local/bin
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ${HOME}/.fzf && \
    ${HOME}/.fzf/install