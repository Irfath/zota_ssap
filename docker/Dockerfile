#FROM public.ecr.aws/lambda/go:1
#COPY main ${LAMBDA_TASK_ROOT}
#CMD ["main"]

#FROM public.ecr.aws/lambda/provided:al2
#COPY main /var/task/main
#CMD ["main"] 


#FROM public.ecr.aws/lambda/go:1
#COPY main ${LAMBDA_TASK_ROOT}
#CMD ["main"]

# Stage 1: Builder - for compiling Go binary
FROM public.ecr.aws/lambda/go:1.21 AS builder

# Set working directory
WORKDIR /build

# Copy go module files for dependency caching
COPY app/go.mod app/go.sum ./
RUN go mod download

# Copy source code
COPY app/main.go .

# Build optimized binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -tags netgo \
    -o main main.go

# Stage 2: Final Lambda image
FROM public.ecr.aws/lambda/go:1.21

# Set timezone for consistent logs
ENV TZ=UTC

# Set working directory
WORKDIR ${LAMBDA_TASK_ROOT}

# Copy binary from builder stage
COPY --from=builder /build/main .

# Ensure binary is executable
RUN chmod 755 main

# Metadata labels for image scanning
LABEL maintainer="irfathr@gmail.com" \
      version="1.0" \
      description="Zota Secure Serverless Asset Proxy"

# Command to run the Lambda
CMD ["main"]