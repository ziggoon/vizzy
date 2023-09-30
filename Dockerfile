# Use an official Go runtime as a parent image
FROM golang:1.19
# Set the working directory inside the container
WORKDIR /app
# Copy the source code from your host to the working directory inside the container
COPY . .
# Build your Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o myapp
# Expose port 42029
EXPOSE 42069
# Command to run your application
CMD ["./myapp"]