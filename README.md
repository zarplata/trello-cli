# trello-cli

trello-cli - Easy trello client in order to add a card on the Board.

### Supported features

Add a card on the Board.

### Installation

#### Manual build

```sh
# Building
git clone https://github.com/zarplata/trello-cli.git
cd trello-cli
make

#Installing
make install

# By default, binary installs into /usr/bin/ and zabbix config in /etc/zabbix/zabbix_agentd.conf.d/ but,
# you may manually copy binary to your executable path and zabbix config to specific include directory
```

#### Arch Linux package
```sh
# Building
git clone https://github.com/zarplata/trello-cli.git
cd trello-cli
git checkout pkgbuild

makepkg

#Installing
pacman -U *.tar.xz
```

### Configuration

#### Login in trello.com

Login in https://trello.com or create account.

Take AppKey from https://trello.com/app-key.

Next generate Token and take him.

#### Run setup mode

Run setup mode and enter AppKey, Token and select Board and List:
```sh
trello-cli -S
```

### Add a card on the Board

```sh
trello-cli CardName CardDescription
```
