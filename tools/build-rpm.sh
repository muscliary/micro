
#Builds two .rpm packages, for x86 (i386) and amd64 (amd64)
#Based on the build-deb.sh but for rpm instead

function BuildRPMSpec() {
dategenerated=$(date +%F_%H:%M)
echo Name: micro
echo Version: $(echo $2 | tr "-" "." )
echo Release: 1
echo "Summary: A feature-rich terminal text editor"
echo URL: https://micro-editor.github.io
echo "Packager: Zachary Yedidia \<zyedidia@gmail.com\>"
echo License: MIT
if [ $1 == "amd64" ]
then
echo BuildArch: "x86_64"
fi
if [ $1 == "i386" ]
then
echo BuildArch: "x86"
fi
echo "Group: Applications/Editors"
echo "%description"
echo "A modern and intuitive terminal-based text editor."
echo " This package contains a modern alternative to other terminal-based"
echo " Editors. It is easy to use, supports mouse input, and is customizable"
echo " via themes and plugins."
echo "%install"
echo "mkdir -p /usr/share/doc/micro"
echo "install -m 755 micro /usr/local/bin/micro"
echo "install -m 744 AUTHORS /usr/share/doc/micro/AUTHORS"
echo "install -m 744 LICENSE /usr/share/doc/micro/LICENSE"
echo "install -m 744 LICENSE-THIRD-PARTY /usr/share/doc/micro/LICENSE-THIRD-PARTY"
echo "install -m 744 README.md /usr/share/doc/micro/README.md"
echo "install -m 744 micro.1.gz /usr/share/man/man1/micro.1.gz"
echo "install -m 744 micro.desktop /usr/share/applications/micro.desktop"
echo "install -m 744 micro.svg /usr/share/icons/micro.svg"
echo "%files"
echo "/usr/local/bin/micro"
echo "/usr/share/doc/micro"
echo "/usr/share/doc/micro/AUTHORS"
echo "/usr/share/doc/micro/LICENSE"
echo "/usr/share/doc/micro/LICENSE-THIRD-PARTY"
echo "/usr/share/doc/micro/README.md"
echo "/usr/share/man/man1/micro.1.gz"
echo "/usr/share/applications/micro.desktop"
echo "/usr/share/icons/micro.svg"
echo "%changelog"
echo "*Version: $1-$2"
echo "*Auto generated on $dategenerated by $USER@$HOSTNAME"
}

function installFiles() {
TO="$1/$2/usr/share/doc/micro/"
mkdir -p $TO
mkdir -p "$1/$2/usr/share/man/man1/"
mkdir -p "$1/$2/usr/share/applications/"
mkdir -p "$1/$2/usr/share/icons/"
cp ../AUTHORS $TO
cp ../LICENSE $TO
cp ../LICENSE-THIRD-PARTY $TO
cp ../README.md $TO
gzip -c ../assets/packaging/micro.1 > $1/$2/usr/share/man/man1/micro.1.gz
cp ../assets/packaging/micro.desktop $1/$2/usr/share/applications/
cp ../assets/logo.svg $1/$2/usr/share/icons/micro.svg
}

version=$1
if [ "$1" == "" ]
then
  version=$(go run build-version.go)
fi
echo "Building packages for Version '$version'"
echo "Running Cross-Compile"
./cross-compile.sh $version

echo "Beginning package build process"

PKGPATH="../packages/rpm"

rm -fr $PKGPATH
mkdir -p $PKGPATH/amd64/
mkdir -p $PKGPATH/i386/
mkdir -p $PKGPATH/arm/

BuildRPMSpec "amd64" "$version" > "$PKGPATH/amd64/micro-$version-amd64.spec"
#BuildRPMSpec "amd64" "$version"
tar -xzf "../binaries/micro-$version-linux64.tar.gz" "micro-$version/micro"
mkdir -p $PKGPATH/amd64/usr/local/bin/
mv "micro-$version/micro" "$PKGPATH/amd64/usr/local/bin"

BuildRPMSpec "i386" "$version" > "$PKGPATH/i386/micro-$version-i386.spec"
#BuildRPMSpec "i386" "$version"
tar -xzf "../binaries/micro-$version-linux32.tar.gz" "micro-$version/micro"
mkdir -p $PKGPATH/i386/usr/local/bin/
mv "micro-$version/micro" "$PKGPATH/i386/usr/local/bin/"

BuildRPMSpec "arm" "$version" > "$PKGPATH/arm/micro-$version-arm.spec"
tar -xzf "../binaries/micro-$version-linux-arm.tar.gz" "micro-$version/micro"
mkdir -p $PKGPATH/arm/usr/local/bin
mv "micro-$version/micro" "$PKGPATH/arm/usr/local/bin"

rm -rf "micro-$version"

installFiles $PKGPATH "amd64"
installFiles $PKGPATH "i386"
installFiles $PKGPATH "arm"

rpmbuild -bb --buildroot $PKGPATH/amd64 $PKGPATH/amd64/micro-$version-amd64.spec
rpmbuild -bb --buildroot $PKGPATH/i386 $PKGPATH/i386/micro-$version-i386.spec
rpmbuild -bb --buildroot $PKGPATH/arm $PKGPATH/arm/micro/$version-arm.spec
