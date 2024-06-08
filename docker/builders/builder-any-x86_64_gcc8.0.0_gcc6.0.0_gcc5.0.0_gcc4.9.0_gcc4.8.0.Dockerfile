FROM debian:buster

LABEL maintainer="cncf-diginfra-dev@lists.cncf.io"

ARG TARGETARCH

RUN cp /etc/skel/.bashrc /root && cp /etc/skel/.profile /root

RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
	bash-completion \
	bc \
	clang \
    	llvm \
	ca-certificates \
	curl \
	dkms \
	dwarves \
	gnupg2 \
	gcc \
	jq \
	libc6-dev \
	libelf-dev \
	netcat \
	xz-utils \
	rpm2cpio \
	cpio \
	flex \
	bison \
	openssl \
	libssl-dev \
	libncurses-dev \
	libudev-dev \
	libpci-dev \
	libiberty-dev \
	lsb-release \
	wget \
	software-properties-common \
	gpg \
	zstd \
	gawk \
	mawk \
	cmake \
	git \
	&& rm -rf /var/lib/apt/lists/*

RUN if [ "$TARGETARCH" = "amd64" ] ; then apt-get install -y --no-install-recommends libmpx2; fi

# gcc 4.9 is required on x86 to build some 3.10+ kernels
# note: on arm gcc 4.9 could not be found.
RUN echo 'deb http://dk.archive.ubuntu.com/ubuntu/ xenial main' >> /etc/apt/sources.list     && \
	echo 'deb http://dk.archive.ubuntu.com/ubuntu/ xenial universe' >> /etc/apt/sources.list
# repo is unsigned therefore the APT options
RUN if [ "$TARGETARCH" = "amd64" ] ; then apt-get -o Acquire::AllowInsecureRepositories=true -o Acquire::AllowDowngradeToInsecureRepositories=true update && apt-get -o APT::Get::AllowUnauthenticated=true install -y --no-install-recommends gcc-4.9; fi

# gcc 6 is no longer included in debian stable, but we need it to
# build kernel modules on the default debian-based ami used by
# kops. So grab copies we've saved from debian snapshots with the
# prefix https://snapshot.debian.org/archive/debian/20170517T033514Z
# or so.

RUN if [ "$TARGETARCH" = "amd64" ]; then curl -L -o libcilkrts5_6.3.0-18_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libcilkrts5_6.3.0-18_${TARGETARCH}.deb; fi; \
    curl -L -o cpp-6_6.3.0-18_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/cpp-6_6.3.0-18_${TARGETARCH}.deb \
	&& curl -L -o gcc-6-base_6.3.0-18_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/gcc-6-base_6.3.0-18_${TARGETARCH}.deb \
	&& curl -L -o gcc-6_6.3.0-18_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/gcc-6_6.3.0-18_${TARGETARCH}.deb \
	&& curl -L -o libasan3_6.3.0-18_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libasan3_6.3.0-18_${TARGETARCH}.deb \
	&& curl -L -o libgcc-6-dev_6.3.0-18_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libgcc-6-dev_6.3.0-18_${TARGETARCH}.deb \
	&& curl -L -o libubsan0_6.3.0-18_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libubsan0_6.3.0-18_${TARGETARCH}.deb \
	&& curl -L -o libmpfr4_3.1.3-2_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libmpfr4_3.1.3-2_${TARGETARCH}.deb \
	&& curl -L -o libisl15_0.18-1_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libisl15_0.18-1_${TARGETARCH}.deb \
	&& dpkg -i cpp-6_6.3.0-18_${TARGETARCH}.deb gcc-6-base_6.3.0-18_${TARGETARCH}.deb gcc-6_6.3.0-18_${TARGETARCH}.deb libasan3_6.3.0-18_${TARGETARCH}.deb; \
    if [ "$TARGETARCH" = "amd64" ]; then dpkg -i libcilkrts5_6.3.0-18_${TARGETARCH}.deb; fi; \
    dpkg -i libgcc-6-dev_6.3.0-18_${TARGETARCH}.deb libubsan0_6.3.0-18_${TARGETARCH}.deb libmpfr4_3.1.3-2_${TARGETARCH}.deb libisl15_0.18-1_${TARGETARCH}.deb \
	&& rm -f cpp-6_6.3.0-18_${TARGETARCH}.deb gcc-6-base_6.3.0-18_${TARGETARCH}.deb gcc-6_6.3.0-18_${TARGETARCH}.deb libasan3_6.3.0-18_${TARGETARCH}.deb libcilkrts5_6.3.0-18_${TARGETARCH}.deb libgcc-6-dev_6.3.0-18_${TARGETARCH}.deb libubsan0_6.3.0-18_${TARGETARCH}.deb libmpfr4_3.1.3-2_${TARGETARCH}.deb libisl15_0.18-1_${TARGETARCH}.deb

# gcc 5 is no longer included in debian stable, but we need it to
# build centos kernels, which are 3.x based and explicitly want a gcc
# version 3, 4, or 5 compiler. So grab copies we've saved from debian
# snapshots with the prefix https://snapshot.debian.org/archive/debian/20190122T000000Z.

RUN if [ "$TARGETARCH" = "amd64" ]; then curl -L -o libmpx0_5.5.0-12_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libmpx0_5.5.0-12_${TARGETARCH}.deb; fi; \
    curl -L -o cpp-5_5.5.0-12_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/cpp-5_5.5.0-12_${TARGETARCH}.deb \
	&& curl -L -o gcc-5-base_5.5.0-12_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/gcc-5-base_5.5.0-12_${TARGETARCH}.deb \
	&& curl -L -o gcc-5_5.5.0-12_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/gcc-5_5.5.0-12_${TARGETARCH}.deb \
	&& curl -L -o libasan2_5.5.0-12_${TARGETARCH}.deb	https://download.diginfra.khulnasoft.com/dependencies/libasan2_5.5.0-12_${TARGETARCH}.deb \
	&& curl -L -o libgcc-5-dev_5.5.0-12_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libgcc-5-dev_5.5.0-12_${TARGETARCH}.deb \
	&& curl -L -o libisl15_0.18-4_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libisl15_0.18-4_${TARGETARCH}.deb \
	&& dpkg -i cpp-5_5.5.0-12_${TARGETARCH}.deb gcc-5-base_5.5.0-12_${TARGETARCH}.deb gcc-5_5.5.0-12_${TARGETARCH}.deb libasan2_5.5.0-12_${TARGETARCH}.deb; \
    if [ "$TARGETARCH" = "amd64" ]; then dpkg -i libmpx0_5.5.0-12_${TARGETARCH}.deb; fi; \
    dpkg -i libgcc-5-dev_5.5.0-12_${TARGETARCH}.deb libisl15_0.18-4_${TARGETARCH}.deb \
	&& rm -f cpp-5_5.5.0-12_${TARGETARCH}.deb gcc-5-base_5.5.0-12_${TARGETARCH}.deb gcc-5_5.5.0-12_${TARGETARCH}.deb libasan2_5.5.0-12_${TARGETARCH}.deb libgcc-5-dev_5.5.0-12_${TARGETARCH}.deb libisl15_0.18-4_${TARGETARCH}.deb libmpx0_5.5.0-12_${TARGETARCH}.deb

# gcc 4 is no longer included in debian stable, but we need it to
# build centos kernels, which are 2.x based and explicitly want a gcc
# version 4 compiler. So grab copies we've saved from debian
# snapshots with the prefix http://ftp.debian.org/debian/pool/main/g/gcc-4.8/.

RUN if [ "$TARGETARCH" = "amd64" ] ; then curl -L -o libasan0_4.8.4-1_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libasan0_4.8.4-1_${TARGETARCH}.deb && dpkg -i libasan0_4.8.4-1_${TARGETARCH}.deb; fi; \
    curl -L -o cpp-4.8_4.8.4-1_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/cpp-4.8_4.8.4-1_${TARGETARCH}.deb \
	&& curl -L -o gcc-4.8-base_4.8.4-1_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/gcc-4.8-base_4.8.4-1_${TARGETARCH}.deb \
	&& curl -L -o gcc-4.8_4.8.4-1_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/gcc-4.8_4.8.4-1_${TARGETARCH}.deb \
	&& curl -L -o libgcc-4.8-dev_4.8.4-1_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libgcc-4.8-dev_4.8.4-1_${TARGETARCH}.deb \
	&& curl -L -o libisl10_0.12.2-2_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libisl10_0.12.2-2_${TARGETARCH}.deb \
	&& curl -L -o multiarch-support_2.19-18+deb8u10_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/multiarch-support_2.19-18%2Bdeb8u10_${TARGETARCH}.deb \
	&& curl -L -o libcloog-isl4_0.18.4-1+b1_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libcloog-isl4_0.18.4-1%2Bb1_${TARGETARCH}.deb \
	&& dpkg -i multiarch-support_2.19-18+deb8u10_${TARGETARCH}.deb \
	&& dpkg -i libisl10_0.12.2-2_${TARGETARCH}.deb gcc-4.8-base_4.8.4-1_${TARGETARCH}.deb; \
    if [ "$TARGETARCH" = "amd64" ] ; then dpkg -i libasan0_4.8.4-1_${TARGETARCH}.deb; fi; \
    dpkg -i libgcc-4.8-dev_4.8.4-1_${TARGETARCH}.deb libcloog-isl4_0.18.4-1+b1_${TARGETARCH}.deb cpp-4.8_4.8.4-1_${TARGETARCH}.deb gcc-4.8_4.8.4-1_${TARGETARCH}.deb \
	&& rm -f multiarch-support_2.19-18+deb8u10_${TARGETARCH}.deb libisl10_0.12.2-2_${TARGETARCH}.deb gcc-4.8-base_4.8.4-1_${TARGETARCH}.deb libasan0_4.8.4-1_${TARGETARCH}.deb libgcc-4.8-dev_4.8.4-1_${TARGETARCH}.deb libcloog-isl4_0.18.4-1+b1_${TARGETARCH}.deb cpp-4.8_4.8.4-1_${TARGETARCH}.deb gcc-4.8_4.8.4-1_${TARGETARCH}.deb

# debian:stable head contains binutils 2.31, which generates
# binaries that are incompatible with kernels < 4.16. So manually
# forcibly install binutils 2.30-22 instead.

RUN if [ "$TARGETARCH" = "amd64" ] ; then \
    curl -L -o binutils-x86-64-linux-gnu_2.30-22_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/binutils-x86-64-linux-gnu_2.30-22_${TARGETARCH}.deb; \
    else  \
    curl -L -o  binutils-aarch64-linux-gnu_2.30-22_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/binutils-aarch64-linux-gnu_2.30-22_${TARGETARCH}.deb; \
    fi

RUN curl -L -o binutils_2.30-22_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/binutils_2.30-22_${TARGETARCH}.deb \
	&& curl -L -o libbinutils_2.30-22_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/libbinutils_2.30-22_${TARGETARCH}.deb \
	&& curl -L -o binutils-common_2.30-22_${TARGETARCH}.deb https://download.diginfra.khulnasoft.com/dependencies/binutils-common_2.30-22_${TARGETARCH}.deb \
	&& dpkg -i *binutils*.deb \
	&& rm -f *binutils*.deb

# Properly create soft link
RUN ln -s /usr/bin/gcc-4.8 /usr/bin/gcc-4.8.0
RUN if [ "$TARGETARCH" = "amd64" ] ; then ln -s /usr/bin/gcc-4.9 /usr/bin/gcc-4.9.0; fi;
RUN ln -s /usr/bin/gcc-5 /usr/bin/gcc-5.0.0
RUN ln -s /usr/bin/gcc-6 /usr/bin/gcc-6.0.0
RUN ln -s /usr/bin/gcc-8 /usr/bin/gcc-8.0.0
