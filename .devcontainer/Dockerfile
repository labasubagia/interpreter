FROM golang:1.22-bookworm

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y git zsh sudo pre-commit

ARG USER=nonroot
ARG UID=1000
ARG GID=1000
RUN groupadd --gid $GID $USER \
    && useradd --uid $UID --shell $(which zsh) --gid $GID -m $USER \
    && echo $USER ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USER \
    && chmod 0440 /etc/sudoers.d/$USER

RUN apt-get install -y python3 python3-pip python3-venv

RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install github.com/cweill/gotests/gotests@latest \
    && go install github.com/fatih/gomodifytags@latest \
    && go install github.com/haya14busa/goplay/cmd/goplay@latest \
    && go install golang.org/x/tools/cmd/goimports@latest

USER $USER

RUN touch /home/$USER/.zshrc

USER root
