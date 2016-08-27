# This script creates releases on Github for micro
# It assumes that the binaries are in the current directory
# You must have the correct Github access token to run this script

# $1 is the title, $2 is the description

commitID=$(git rev-parse HEAD)

echo "Creating tag"
git tag $1 $commitID
git push --tags

echo "Creating new release"
github-release release \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "$1" \
    --description "$2" \
    --pre-release

echo "Uploading OSX binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-osx.tar.gz" \
    --file binaries/micro-$1-osx.tar.gz

echo "Uploading Linux 64 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-linux64.tar.gz" \
    --file binaries/micro-$1-linux64.tar.gz

echo "Uploading Linux 32 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-linux32.tar.gz" \
    --file binaries/micro-$1-linux32.tar.gz

echo "Uploading Linux Arm binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-linux-arm.tar.gz" \
    --file binaries/micro-$1-linux-arm.tar.gz

echo "Uploading FreeBSD 64 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-freebsd64.tar.gz" \
    --file binaries/micro-$1-freebsd64.tar.gz

echo "Uploading FreeBSD 32 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-freebsd32.tar.gz" \
    --file binaries/micro-$1-freebsd32.tar.gz

echo "Uploading OpenBSD 64 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-openbsd64.tar.gz" \
    --file binaries/micro-$1-openbsd64.tar.gz

echo "Uploading OpenBSD 32 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-openbsd32.tar.gz" \
    --file binaries/micro-$1-openbsd32.tar.gz

echo "Uploading NetBSD 64 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-netbsd64.tar.gz" \
    --file binaries/micro-$1-netbsd64.tar.gz

echo "Uploading NetBSD 32 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-netbsd32.tar.gz" \
    --file binaries/micro-$1-netbsd32.tar.gz

echo "Uploading Windows 64 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-win64.zip" \
    --file binaries/micro-$1-win64.zip

echo "Uploading Windows 32 binary"
github-release upload \
    --user zyedidia \
    --repo micro \
    --tag $1 \
    --name "micro-$1-win32.zip" \
    --file binaries/micro-$1-win32.zip
