#!/bin/bash

# This is tmp solution for upload release because GitHub action cannot build 
# for ARM arch.

RELEASE_NAME=""
RELEASE_DESC=""
REPO_USER=""
TOKEN=""

REPO_NAME="rpi4_network_controller"
FILE_PATH="./network"

RELEASE_RESPONSE=$(curl --data "{\"tag_name\": \"${RELEASE_NAME}\", \"name\": \"${RELEASE_NAME}\", \"body\": \"${RELEASE_DESC}\"}" \
  -H "Authorization: token ${TOKEN}" \
  -H "Content-Type: application/json" \
  -X POST https://api.github.com/repos/"${REPO_USER}"/"${REPO_NAME}"/releases)

UPLOAD_URL=$(echo "${RELEASE_RESPONSE}" | jq -r .upload_url | sed -e "s/{?name,label}//")

# Check for errors in the response
if echo "${RELEASE_RESPONSE}" | grep -q "message"; then
  echo "Error creating the release:"
  echo "${RELEASE_RESPONSE}" | jq .message
  exit 1
fi

if [ -z "${UPLOAD_URL}" ] || [ "${UPLOAD_URL}" == "null" ]; then
  echo "Failed to get the upload URL. Check if the release was created successfully or if you've exceeded the rate limit."
  exit 1
fi

# Upload a file to the release
curl -H "Authorization: token $TOKEN" \
  -H "Content-Type: $(file -b --mime-type ${FILE_PATH})" \
  -X POST "${UPLOAD_URL}?name=$(basename ${FILE_PATH})" \
  --upload-file "${FILE_PATH}"
