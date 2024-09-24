# Use Amazon Linux 2023 as the base image
FROM amazonlinux:2023

# Set Go version
ENV GO_VERSION=1.23.1
ENV PACKAGE_VERSION=1.0.0

# Install necessary build tools and dependencies
RUN yum install -y \
        rpm-build \
        rpmdevtools \
        yum-utils \
        gcc \
        git \
        wget \
        tar \
        make \
    && yum clean all

# Install Go from the official tarball
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"

# Create a user for building
RUN useradd -ms /bin/bash builder
USER builder
WORKDIR /home/builder

# Set up RPM build environment
RUN rpmdev-setuptree

# Copy your Go application source code
COPY --chown=builder:builder . /home/builder/build

RUN pwd && cd rpmbuild && ls -lha && uname -m

# Change to the build directory
WORKDIR /home/builder/build

# RUN tar czf postfix-aws-cassandra-${PACKAGE_VERSION}.tar.gz src/ postfix-aws-cassandra.service

RUN mkdir postfix-aws-cassandra-${PACKAGE_VERSION} && \
    cp -r src/ postfix-aws-cassandra.service postfix-aws-cassandra-${PACKAGE_VERSION}/ && \
    tar czf postfix-aws-cassandra-${PACKAGE_VERSION}.tar.gz postfix-aws-cassandra-${PACKAGE_VERSION}

# Build the Go application
# RUN pwd && make build

# Copy the spec file
RUN cp /home/builder/build/postfix-aws-cassandra.spec /home/builder/rpmbuild/SPECS/

# Copy the binary to the RPM build directory
RUN cp /home/builder/build/postfix-aws-cassandra-${PACKAGE_VERSION}.tar.gz /home/builder/rpmbuild/SOURCES/

# Build the RPM package
RUN rpmbuild -bb /home/builder/rpmbuild/SPECS/postfix-aws-cassandra.spec

# The built RPM is now in /home/builder/rpmbuild/RPMS/x86_64/