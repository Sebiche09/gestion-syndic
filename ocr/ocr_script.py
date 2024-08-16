from flask import Flask, request, jsonify
from paddleocr import PaddleOCR
from PIL import Image
import io
import numpy as np
import gc

app = Flask(__name__)

# Initialisez PaddleOCR avec les paramètres appropriés
ocr = PaddleOCR(use_angle_cls=True, lang='fr')

@app.route('/ocr', methods=['POST'])
def ocr_endpoint():
    if 'file' not in request.files:
        return jsonify({'error': 'Aucun fichier fourni'}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({'error': 'Aucun fichier sélectionné'}), 400

    try:
        # Lire l'image depuis le fichier
        image = Image.open(io.BytesIO(file.read()))

        # Convertir l'image PIL en tableau NumPy
        image_np = np.array(image)

        # Effectuer l'OCR avec PaddleOCR
        result = ocr.ocr(image_np, cls=True)

        # Préparer le texte extrait
        text = '\n'.join([line[1][0] for line in result[0]])

        # Nettoyer les ressources temporaires
        del image_np
        del result
        image.close()

        # Forcer la collecte des déchets pour libérer la mémoire
        gc.collect()

        return jsonify({'text': text}), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
