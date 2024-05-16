# Use the Ubuntu 22.04 base image
FROM ubuntu:22.04

# Set the working directory
WORKDIR /app

# Install necessary packages
RUN apt-get update && apt-get install -y \
    software-properties-common \
    ca-certificates \
    wget \
    tar \
    && add-apt-repository ppa:ubuntu-toolchain-r/test \
    && apt-get install -y \
    gcc-11 \
    g++-11 \
    && update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-11 100 --slave /usr/bin/g++ g++ /usr/bin/g++-11 \
    && update-alternatives --config gcc \
    && rm -rf /var/lib/apt/lists/*

# Install Go
ENV GO_VERSION=1.22.1
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Set environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Install TensorFlow and Tesseract
RUN apt-get update && apt-get install -y \
    wget \
    tar \
    && wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.15.0.tar.gz \
    && tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-2.15.0.tar.gz \
    && ldconfig \
    && echo "/usr/local/lib" | tee -a /etc/ld.so.conf.d/libtensorflow.conf \
    && ldconfig \
    && apt-get install -y \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr-eng \
    tesseract-ocr-deu \
    tesseract-ocr-jpn \
    tesseract-ocr-ara \
    && rm -rf /var/lib/apt/lists/*

# Build the Go app
RUN go build ./cmd/server

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./server"]
