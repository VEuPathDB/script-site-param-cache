#!/usr/bin/env sh

#
# Constants
#
readonly GITHUB_API_PATH="https://api.github.com"
readonly REPO_API_PATH="${GITHUB_API_PATH}/repos/VEuPathDB/script-site-param-cache"
readonly REPO_TARGET="${REPO_API_PATH}/releases/latest"
readonly BINARY_NAME="param-cache"

#
# Execution Time
#
readonly FILE_URL="$(curl -s "${REPO_TARGET}" \
  | jq -r '.assets[].browser_download_url | select(. | match("linux"))')"
readonly FILE_NAME="$(basename "${FILE_URL}")"


#
# Do the thing
#
wget -q "${FILE_URL}" \
  && tar -xzf "${FILE_NAME}" \
  && rm "${FILE_NAME}" \
  && ./param-cache $@
rm -rf ./param-cache
