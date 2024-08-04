#!/bin/sh

echo "PassivePorts 30000 30009" >> /etc/pure-ftpd/pure-ftpd.conf

# Créer les répertoires nécessaires
mkdir -p /home/ftpusers/user/temporary/pending
mkdir -p /home/ftpusers/user/temporary/processed
mkdir -p /home/ftpusers/user/final

# Créer un utilisateur virtuel et définir le mot de passe
(echo "$FTP_USER_PASS"; echo "$FTP_USER_PASS") | pure-pw useradd $FTP_USER_NAME -u ftpuser -d $FTP_USER_HOME

# Mettre à jour la base de données des utilisateurs
pure-pw mkdb

# Changer les permissions du répertoire utilisateur
chown -R ftpuser:ftpgroup /home/ftpusers/user
chmod -R 755 /home/ftpusers/user

# Démarrer Pure-FTPd
exec "$@"