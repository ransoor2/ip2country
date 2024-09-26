#!/bin/bash

# Check if the kind cluster already exists
if kind get clusters | grep -q "^ip2country$"; then
  echo "Cluster 'ip2country' already exists."
else
  # Create a new kind cluster
  kind create cluster --name ip2country --config scripts/kind_config.yaml
fi

# Tag the Docker image
docker tag ip2country:latest kind.local/ip2country:latest

# Load the Docker image into the kind cluster
kind load docker-image kind.local/ip2country:latest --name ip2country

# Apply Kubernetes configurations
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/