#!/usr/bin/env bash
#
# Script to mirror a tag from Prysm into EthereumAPIs protocol buffers
#
# Example:
#
# mirror-ethereumapis.sh
#
set -e

# Validate settings.
[ "$TRACE" ] && set -x

# Define variables.
GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/prysmaticlabs/ethereumapis"

AUTH="Authorization: token $GITHUB_SECRET_ACCESS_TOKEN"
# skipcq: SH-2034
export WGET_ARGS="--content-disposition --auth-no-challenge --no-cookie"
# skipcq: SH-2034
export CURL_ARGS="-LJO#"

# Validate token.
curl -o /dev/null -sH "$AUTH" "$GH_REPO" || { echo "Error: Invalid repo, token or network issue!";  exit 1; }

git config --global user.email contact@prysmaticlabs.com
git config --global user.name prylabsbot
git config --global url."https://git:'$GITHUB_SECRET_ACCESS_TOKEN'@github.com/".insteadOf "git@github.com/"

# Clone ethereumapis and prysm
git clone https://github.com/prysmaticlabs/prysm /tmp/prysm/
git clone https://github.com/prysmaticlabs/ethereumapis /tmp/ethereumapis/

# Checkout the release tag in prysm and copy over protos
cd /tmp/prysm && git checkout "$BUILDKITE_BRANCH"
cp -Rf /tmp/prysm/proto/eth /tmp/ethereumapis
cd /tmp/ethereumapis || exit

# Replace imports in go files and proto files as needed
find ./eth -name '*.go' -print0 |
    while IFS= read -r -d '' line; do
        sed -i 's/prysm\/proto\/eth/ethereumapis\/eth/g' "$line"
    done

find ./eth -name '*.go' -print0 |
    while IFS= read -r -d '' line; do
        sed -i 's/proto\/eth/eth/g' "$line"
    done

find ./eth -name '*.go' -print0 |
    while IFS= read -r -d '' line; do
        sed -i 's/proto_eth/eth/g' "$line"
    done

find ./eth -name '*.proto' -print0 |
    while IFS= read -r -d '' line; do
        sed -i 's/"proto\/eth/"eth/g' "$line"
    done

find ./eth -name '*.proto' -print0 |
    while IFS= read -r -d '' line; do
        sed -i 's/prysmaticlabs\/prysm\/proto\/eth/prysmaticlabs\/ethereumapis\/eth/g' "$line"
    done

if git status | grep -q 'nothing to commit'; then
   echo "nothing to push, exiting early"
   exit
fi

# Push to the mirror repository
git add --all
git commit -am "$BUILDKITE_BRANCH"
git push origin master
