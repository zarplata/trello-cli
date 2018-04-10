pkgname=trello-cli
pkgver=1
pkgrel=1
pkgdesc="Trello console client"
arch=('any')
license=('GPL')
makedepends=('go')
depends=()
install=()
source=("https://github.com/zarplata/trello-cli.git")
md5sums=('SKIP')

pkgver() {
    cd "$srcdir/$pkgname"

    make ver
}
    
build() {
    cd "$srcdir/$pkgname"

    make
}

package() {
    cd "$srcdir/$pkgname"

    install -Dm 0755 .out/"${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
}
