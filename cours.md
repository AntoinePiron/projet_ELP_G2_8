# Suivi du projet
*Le but de ce fichier est de suivre l'évolution de notre projet est d'en garder une trace écrite notamment pour les bilans*  

## Séance 1 
mail prof : pierre.francois@insa-lyon.fr
Projet :
    - groupe de 3 de preference mixtes
    - A la fin de chaque cours bilan, même minime, au prof par mail avec objet specifique : [GO] GR2-<numProjet> etc

Objectif de cette séance :   
    1. Installation GO  
    2. Formation du groupe  
    3. Programme qui ouvre un fichier, lire le fichier et afficher à l'écran les derniers mots de chaque ligne   
    4. Créer un repo git   
    5. Préparer le mail    

## Notre groupe : GR2-8 
Piste de recherche :   
    - Traitement d'image   
    - Compression/amélioration qualité d'image    

Premier travail : tranformer une image en noir et blanc    
Réussite ! --> code dispo dans *image_bw/main.go*    

A voir pour go routine -> systeme de canal    
Sinon plusieurs go routine qui ecrive dans le meme fichier mais il faut verifier si il est thread safe -> mutex     

## Séance 2 (autonomie)
On test notre programe avec/sans go routines pour une image en 6000x4000 donc plutot lourde.   
Sans go routines : 2,75 secondes   
Avec go routine : envrion 10 secondes    
Conclusion : il y a un problème dans le code avec des go routines il va falloir trouver un moyen d'optimiser la chose   

## Séance 3 du 08/10/2021
Début de reflexion :  
On est passé dans un premier temps sur un buffered channel -> gain de 2 secondes  
On a ensuite essayé de créer un autre programme fonctionnant avec un système de mutex qui va directement écrire chaque pixel dans le fichier final on se retrouve avec un temps d'execution avoisinant les 4 secondes pour 4 subdivisions. Ce n'est toujours pas mieux que le programe sans go routines mais on s'améliore.  
On a enfin réussi à rendre les go routines plus efficace !  
Solutions ? Dans la go routine on utilisait un mutex qui bloquait à chaque écriture de pixel or cela n'est pas efficaace car il faudrait faire un blocage si jamais deux go routines tentent d'écrire au meme endroit ce qui n'est pas le cas pour nous : chaque go routines écrit dans des pixels différents donc pas de risque de collision.   
