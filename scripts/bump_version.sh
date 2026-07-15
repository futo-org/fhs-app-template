#!/usr/bin/env sh
# Bump the semver in server/VERSION and print the new version.
#   usage: scripts/bump_version.sh major|minor|patch
set -eu

usage() {
	echo "usage: $0 major|minor|patch" >&2
	exit 1
}
[ $# -eq 1 ] || usage

file="$(dirname "$0")/../server/VERSION"
version="$(tr -d '[:space:]' <"$file")"

IFS=. read -r major minor patch <<EOF
$version
EOF

for part in "$major" "$minor" "$patch"; do
	case "$part" in
	"" | *[!0-9]*)
		echo "error: '$version' in server/VERSION is not a valid x.y.z version" >&2
		exit 1
		;;
	esac
done

case "$1" in
major)
	major=$((major + 1))
	minor=0
	patch=0
	;;
minor)
	minor=$((minor + 1))
	patch=0
	;;
patch) patch=$((patch + 1)) ;;
*) usage ;;
esac

new="${major}.${minor}.${patch}"
printf '%s\n' "$new" >"$file"
echo "$new"
