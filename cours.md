# Séance 1 
mail prof : pierre.francois@insa-lyon.fr
Projet :
    - groupe de 3 de preference mixtes
    - A la fin de chaque cours bilan, même minime, au prof par mail avec objet specifique : [GO] GR2-<numProjet> etc

Objectif de cette séance :
    - Installation GO
    - Formation du groupe
    - Programme qui ouvre un fichier, lire le fichier et afficher à l'écran les derniers mots de chaque ligne 
    - Créer un repo git
    - Préparer le mail 

## Notre groupe : GR2-8 
Piste de recherche :
    - Traitement d'image
    - Compression/amélioration qualité d'image

Premier travail : tranformer une image en noir et blanc 

A voir pour go routine -> systeme de canal
Sinon plusieurs go routine qui ecrive dans le meme fichier mais il faut verifier si il est thread safe -> mutex

## Bilan test image 6000x4000
Sans go routines : 2,75 secondes

