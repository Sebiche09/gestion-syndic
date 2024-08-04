import os
import time
from PIL import Image
import pytesseract

# Chemin vers le répertoire FTP monté
ftp_dir = '/home/ftpusers/user'
log_file = '/app/ocr_log.txt'  # Chemin vers le fichier de log

def log_message(message):
    """ Écrit un message dans le fichier de log """
    with open(log_file, 'a') as f:
        f.write(message + '\n')

def process_image(file_path):
    """ Extrait le texte de l'image et le log """
    try:
        image = Image.open(file_path)
        text = pytesseract.image_to_string(image)
        log_message(f"Texte extrait de {file_path} :")
        log_message(text)
    except Exception as e:
        log_message(f"Erreur lors du traitement de {file_path}: {e}")

def process_directory(directory):
    """ Traite les fichiers dans le répertoire """
    for root, dirs, files in os.walk(directory):
        for file in files:
            file_path = os.path.join(root, file)
            if file_path.endswith(('.png', '.jpg', '.jpeg', '.tiff', '.bmp')):  # Extension des images à traiter
                process_image(file_path)

def main():
    processed_files = set()
    while True:
        for root, dirs, files in os.walk(ftp_dir):
            for file in files:
                file_path = os.path.join(root, file)
                if file_path not in processed_files:
                    process_image(file_path)
                    processed_files.add(file_path)
        time.sleep(10)  # Attendre avant de vérifier à nouveau

if __name__ == "__main__":
    main()
