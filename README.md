# INSATutorat  
[![wakatime](https://wakatime.com/badge/user/395b07f0-60f3-4a92-a7fd-d4c38328182c/project/51582e5a-5406-4da2-9824-d5b913e2954e.svg)](https://wakatime.com/badge/user/395b07f0-60f3-4a92-a7fd-d4c38328182c/project/51582e5a-5406-4da2-9824-d5b913e2954e)

**INSATutorat** est une application développée dans le cadre d'un Travail d'Initiative Personnelle (TIP) pour le département STPI de l'INSA Rouen Normandie.  
Elle a pour objectif de faciliter la gestion du tutorat entre étudiants grâce à une plateforme web intuitive et complète.

---

## 📄 Documents du projet

- [Cahier des charges](./TIP_Cahier_des_charges.pdf)
- [Présentation diapo](./TIP_Présentation_diapo.pdf)
- [Rapport final](./TIP_Rapport_final.pdf)

---

## 🚀 Fonctionnalités principales

- Gestion des **tuteurs** et **tutorés**  
- Attribution automatique ou manuelle des appariements  
- Visualisation des disponibilités des étudiants sous forme de calendrier  
- Interface d'administration pour gérer les inscriptions, affectations et quotas  
- Suivi des sessions de tutorat et des heures réalisées

---

## 🛠️ Architecture du projet

Le projet est divisé en deux parties principales :  

- **Backend (Go + Gin + GORM)** :  
  Fournit une API REST pour gérer les inscriptions, les appariements, les comptes, et l'administration.  

- **Frontend (Vue/Nuxt + Tailwind)** :  
  Interface utilisateur permettant aux étudiants et administrateurs d'interagir avec la plateforme.  

Base de données : **MariaDB** (via GORM).

---

## 📋 Prérequis

Avant d'installer et de lancer le projet, assurez-vous d'avoir :  

- **Golang** (testé avec 1.24.0) 
- **Node.js** (testé avec 20.11.0) et **npm** ≥ (testé avec 10.8.1)
- **MariaDB** (ou MySQL compatible)  
- (Optionnel) un serveur web pour la version statique (NGINX, Caddy, ...)

---

## ⚙️ Installation & Lancement

### 1. Cloner le projet
```bash
git clone https://github.com/Romitou/INSATutorat
cd INSATutorat
````

### 2. Configurer la base de données

Créez une base MariaDB et configurez vos identifiants dans le fichier de configuration du backend (par ex. `.env` ou variables d'environnements).

### 3. Lancer le backend

```bash
go build
./insatutorat
```

Le serveur API est maintenant actif.

### 4. Installer et lancer le frontend

```bash
cd client
npm install
npm run dev
```

Le frontend est accessible sur [http://localhost:3000](http://localhost:3000).

---

## 📦 Build pour la production

Pour générer une version statique du frontend :

```bash
BASE_URL=https://tutorat-stpi.foo.bar npm run generate
```

Les fichiers générés se trouvent dans `client/dist`.
Ils peuvent être servis avec **NGINX**, **Caddy**, ou tout autre serveur web.

Le backend (Go) peut être déployé indépendamment comme binaire.
