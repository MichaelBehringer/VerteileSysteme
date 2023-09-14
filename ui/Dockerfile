# Verwenden Sie ein offizielles Node.js-Image als Basis
FROM node:lts-alpine

# Erstellen Sie das Arbeitsverzeichnis im Container
WORKDIR /appUI

# Kopieren Sie die Abhängigkeiten und den Code in den Container
COPY package*.json ./
RUN npm install
COPY . .

# Bauen Sie die React-Anwendung
RUN npm run build
RUN npm install -g serve

# Starten Sie die Anwendung
CMD ["serve", "-s", "build"]