#!/usr/bin/env bash
set -euo pipefail

LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -z "$LAST_TAG" ]; then
  COMMITS=$(git log HEAD --pretty=format:"%s")
else
  COMMITS=$(git log "${LAST_TAG}..HEAD" --pretty=format:"%s")
fi

if [ -z "$COMMITS" ]; then
  echo "No commits since ${LAST_TAG}, skipping release."
  exit 0
fi

BUMP="patch"
while IFS= read -r msg; do
  if echo "$msg" | grep -qE "^(feat|fix|refactor|perf)(\(.+\))?!:|^BREAKING CHANGE"; then
    BUMP="major"
    break
  elif echo "$msg" | grep -qE "^feat(\(.+\))?:"; then
    if [ "$BUMP" != "major" ]; then
      BUMP="minor"
    fi
  fi
done <<< "$COMMITS"

VERSION=${LAST_TAG#v}
VERSION=${VERSION:-0.0.0}
MAJOR=$(echo "$VERSION" | cut -d. -f1)
MINOR=$(echo "$VERSION" | cut -d. -f2)
PATCH=$(echo "$VERSION" | cut -d. -f3)

case $BUMP in
  major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
  minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
  patch) PATCH=$((PATCH + 1)) ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

git tag "${NEW_VERSION}"
git push origin "${NEW_VERSION}"
