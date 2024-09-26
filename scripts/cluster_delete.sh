#!/bin/bash

# Check if the kind cluster already exists
if kind get clusters | grep -q "^ip2country$"; then
  # Delete a new kind cluster
    kind delete cluster --name ip2country
else
  echo "Cluster 'ip2country' doesnt exists."
fi

