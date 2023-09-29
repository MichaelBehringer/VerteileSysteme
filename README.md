# VerteileSysteme
## Wichtige Vorraussetzung
Der Ordner, im welem sich dieses Projekt befindet muss VerteilteSysteme heißen
Docker muss installiert sein
Port 80 auf PC ist offen
## Starten
sudo docker-compose --compatibility up --build
(compatibility stellt sicher, dass die Container auf Windows und Linux gleich benannt werden)
## Anwendung
Wenn die Anwendung gestartet ist, kann unter http://localhost/ die Website erreicht werden. Bzw die IP-Adresse des Servers auf dem der Container läuft falls dieser nicht lokal ist.
## Nutzer
Es werden automatisch Test benutzer angelegt:
Benutzername Passwort
Berry abc123
Dj smash
Marco kommentar
## Sonstiges
Viel Spaß ^^