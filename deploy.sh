#!/bin/bash

# Sailtix Documentation Server Deployment Script
# This script builds and deploys the secure documentation server

set -e

echo "🚀 Starting Sailtix Documentation Server Deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if Go is installed (for local development)
if command -v go &> /dev/null; then
    print_status "Go is installed. You can also run the server locally."
else
    print_warning "Go is not installed. You can only run the server with Docker."
fi

# Check if required files exist
required_files=("index.html" "openapi.yaml" "openapi-v2.yaml" "agents.json")
for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        print_error "Required file $file is missing!"
        exit 1
    fi
done

print_success "All required files found."

# Create config.json if it doesn't exist
if [ ! -f "config.json" ]; then
    print_status "Creating default config.json..."
    cat > config.json << EOF
{
  "port": ":8080",
  "session_secret": "$(openssl rand -hex 32)",
  "session_duration_hours": 24,
  "max_login_attempts": 5,
  "lockout_duration_minutes": 15
}
EOF
    print_success "Default config.json created with secure session secret."
fi

# Function to run with Docker
run_with_docker() {
    print_status "Building and running with Docker Compose..."
    
    # Stop existing containers
    docker-compose down 2>/dev/null || true
    
    # Build and start
    docker-compose up --build -d
    
    print_success "Documentation server is running!"
    print_status "Access URL: http://localhost:8080"
    print_status "To view logs: docker-compose logs -f"
    print_status "To stop: docker-compose down"
}

# Function to run locally with Go
run_locally() {
    print_status "Running locally with Go..."
    
    # Download dependencies
    go mod download
    
    # Run the server
    go run server.go
}

# Function to build binary
build_binary() {
    print_status "Building binary..."
    
    # Download dependencies
    go mod download
    
    # Build
    go build -o sailtix-docs server.go
    
    print_success "Binary built successfully: sailtix-docs"
    print_status "Run with: ./sailtix-docs"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  docker     Build and run with Docker (default)"
    echo "  local      Run locally with Go"
    echo "  build      Build binary only"
    echo "  stop       Stop Docker containers"
    echo "  logs       Show Docker logs"
    echo "  help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 docker    # Run with Docker"
    echo "  $0 local     # Run locally with Go"
    echo "  $0 build     # Build binary only"
}

# Main script logic
case "${1:-docker}" in
    "docker")
        run_with_docker
        ;;
    "local")
        if ! command -v go &> /dev/null; then
            print_error "Go is not installed. Use 'docker' option instead."
            exit 1
        fi
        run_locally
        ;;
    "build")
        if ! command -v go &> /dev/null; then
            print_error "Go is not installed."
            exit 1
        fi
        build_binary
        ;;
    "stop")
        print_status "Stopping Docker containers..."
        docker-compose down
        print_success "Containers stopped."
        ;;
    "logs")
        print_status "Showing Docker logs..."
        docker-compose logs -f
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    *)
        print_error "Unknown option: $1"
        show_usage
        exit 1
        ;;
esac

echo ""
print_success "Deployment completed successfully! 🎉" 