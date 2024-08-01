import os
import time
from PIL import Image
import pytesseract

# Chemin vers le répertoire FTP monté
ftp_dir = '/home/ftpusers/user'
log_file = '/app/ocr_log.txt'  # Chemin vers le fichier de log

def log_message(message):
    with open(log_file, 'a') as f:
        f.write(message + '\n')

def process_image(file_path):
    image = Image.open(file_path)
    text = pytesseract.image_to_string(image)
    log_message(f"Texte extrait de {file_path} :")
    log_message(text)

def main():
    processed_files = set()
    while True:
        files = os.listdir(ftp_dir)
        for file in files:
            file_path = os.path.join(ftp_dir, file)
            if file_path not in processed_files:
                process_image(file_path)
                processed_files.add(file_path)
        time.sleep(10)  # Attendre avant de vérifier à nouveau

if __name__ == "__main__":
    main()
