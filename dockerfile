# Use a lightweight base image
FROM golang:1.22.4-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    MY_ENV_VAR=12345 \
    GIN_MODE=release \
    ENGINE_PORT=":8443" \
    DATABASE_PROFILING=false \
    DATABASE_SERVER=172.17.0.3 \
    DATABASE_USERNAME=postgres \
    DATABASE_PASSWORD=mysecret \
    DATABASE_NAME=arctfrex \
    DATABASE_PORT=5432 \
    DATABASE_URL=postgres://user:password@localhost:5432/dbname \
    GOOGLE_SMTP_SERVER=smtp.gmail.com \
    GOOGLE_SMTP_PORT=587 \
    GOOGLE_SMTP_USERNAME=andreanadinata68@gmail.com \
    GOOGLE_SMTP_PASSWORD=jwuuxxeadrfpmsdb \
    EMAIL_FROM=arctfrex@gmail.com \
    OTP_GENERATOR_SECRET=JBSWY3DPEHPK3PXP \
    OTP_SEND_WITH_EMAIL=true \
    OTP_EMAIL_SUBJECT="YOUR OTP" \
    TWILIO_WHATSAPP_USERNAME=AC673bc23395eec3f38f7608578c3adfdc \
    TWILIO_WHATSAPP_PASSWORD=045f32e4555da7d444b7c3703194e742 \
    JWT_SECRET_KEY=m3X8Ib42ea06RSjIL1FAw8 \
    JWT_USERNAME=ARCVIS \
    JWT_PASSWORD=12345 \
    APPLICATION_NAME=ARCTFREX \
    RUN_MARKET_WORKER_PRICE_UPDATES=false \
    RUN_MARKET_WORKER_LIVE_MARKET_UPDATES=false \
    RUN_NEWS_WORKER_LATEST_NEWS_UPDATES=false \
    RUN_NEWS_WORKER_LATEST_NEWS_BULLETIN_UPDATES=false \
    RUN_ORDER_WORKER_CLOSE_ALL_EXPIRED_ORDER=false \
    MINIO_ENDPOINT=minio.arctfrex.com \
    MINIO_BASEURL=storage.arctfrex.com \
    MINIO_ENDPOINT_SECURED=false \
    MINIO_ACCESS_KEY=riFjlIyURmjokSIgNWej \ 
    MINIO_SECRET_KEY=7aba7M3oLAwkZaQaMz3V5MN0CkYQtwdtIauQzxMQ \
    BUCKET_NAME=arctfrex \
    ARC_META_INTEGRATOR_BASEURL=https://arcmetaintegrator.arctfrex.com/api

# Set the working directory
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Change to the directory containing the main.go file
WORKDIR /app/cmd

# Build the binary
RUN go build -o /app/main

# Set the working directory back to /app
WORKDIR /app

# Expose the application port
EXPOSE 8443

# Start the application
CMD ["./main"]
