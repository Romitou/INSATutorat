
# INSATutorat
Cette application a été réalisé dans le cadre d'un Travail d'Initiative Personnelle (TIP) pour le département STPI de l'INSA Rouen Normandie. Il s'agit d'une plateforme permettant de gérer le tutorat au sein de ce département avec différentes fonctionnalités.


## Développement / tests

Pré-requis :
* Golang (testé avec go1.24.0)
* NodeJS & npm (testé avec v20.11.0 & 10.8.1)
* MariaDB (c.f. gorm)

Clonez le projet et accédez au répertoire du dépôt :

```bash
  git clone https://github.com/Romitou/INSATutorat
  cd INSATutorat
```

Buildez et démarrez le binaire du serveur avec :

```bash
  go build
  ./insatutorat
```

Accédez au client et installez les dépendances :

```bash
  cd client
  npm i
```

Démarrez le serveur de développement :

```bash
  npm run dev
```

