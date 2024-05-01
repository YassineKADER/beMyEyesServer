#Todo: Test the dockerfile

# Start from the latest golang base image
FROM golang:1.22.1

# Add Maintainer Info
LABEL maintainer="EROS <yassinekader.contact@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

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
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr-eng \
    tesseract-ocr-deu \
    tesseract-ocr-jpn \
 && rm -rf /var/lib/apt/lists/* \
 && wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.15.0.tar.gz \
 && tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-2.15.0.tar.gz \
 && ldconfig \
 && echo "/usr/local/lib" | tee -a /etc/ld.so.conf.d/libtensorflow.conf \
 && ldconfig

# Build the Go app
RUN go build ./cmd/server

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./server"]
