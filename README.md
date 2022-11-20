# Perché
## Lecteur CLI pour la station FIP

`perché` est un lecteur CLI de flux radio conçu spécifiquement pour la station
[FIP](https://www.radiofrance.fr/fip). Il se base sur le flux HLS officiel,
ainsi que sur l'API proposée par Radio France pour avoir les informations du
titre en cours.

`perché` permet surtout d'enregistrer dans un fichier texte le titre en cours
de diffusion, par simple pression de la touche `s`.


## Fonctionnalités à venir (j'espère)

 * Intégration à Gnome/KDE et macOS : les informations rendues par l'API
   permettent d'afficher une image de couverture et la progression actuelle du
   titre, en plus des informations basiques ;
 * choix du lecteur sous-jacent : actuellement `mpv` est requis et est la seule
   option ;
 * choix du fichier pour les titres sauvegardés : actuellement c'est
   `~/.fip-tracks.txt` seulement.

## Utilisation

```console
% ./perché
---
Title:  Juancito Trucupey
Artist: Celia Cruz
Year:   2011
```

Les titres suivants s'afficheront à la suite. À tout moment, on peut :

 * sauvegarder le titre courant : `s` (insensible à la casse) ;
 * quitter proprement : `q` ou Ctrl+C (insensible ç la casse).

## Compilation

```console
make build
./bin/perché
```
