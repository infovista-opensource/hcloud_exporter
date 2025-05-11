ARG VARIANT="1.23"
FROM golang:${VARIANT}

# Set environment variables
ENV CGO_ENABLED=0 \
    DEBIAN_FRONTEND=noninteractive

# Install needed packages
RUN apt-get update \
    && apt-get -y install --no-install-recommends \
       apt-utils \
       git \
       curl \
       wget \
       unzip \
       zip \
       make \
       gcc \
       libc6-dev \
       sudo \
       vim \
       nano \
       jq \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN groupadd --gid $USER_GID vscode \
    && useradd --uid $USER_UID --gid $USER_GID -m vscode \
    && echo vscode ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/vscode \
    && chmod 0440 /etc/sudoers.d/vscode

# Install Go tools
RUN go install -v golang.org/x/tools/gopls@latest \
    && go install -v honnef.co/go/tools/cmd/staticcheck@latest \
    && go install -v github.com/go-delve/delve/cmd/dlv@latest \
    && go install -v github.com/mgechev/revive@latest \
    && go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set up workspace directory
WORKDIR /workspace

# Switch to vscode user
USER vscode
