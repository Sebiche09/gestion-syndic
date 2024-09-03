import fitz  # PyMuPDF
from fastapi import FastAPI, File, UploadFile, HTTPException
from paddleocr import PaddleOCR
from PIL import Image
import io
import numpy as np
from pdf2image import convert_from_bytes
from typing import Dict
import gc

app = FastAPI()

# Initialiser PaddleOCR
ocr = PaddleOCR(use_angle_cls=True, lang='fr')

def is_text_pdf(doc) -> bool:
    """Vérifie si le PDF contient principalement du texte."""
    text = ""
    for page_num in range(len(doc)):
        page = doc.load_page(page_num)
        text += page.get_text("text")
    return len(text.strip()) > 100  # Adjust threshold as needed

@app.post("/ocr")
async def ocr_endpoint(file: UploadFile = File(...)) -> Dict[str, str]:
    if file.filename == '':
        raise HTTPException(status_code=400, detail="Aucun fichier sélectionné")

    try:
        contents = await file.read()

        # Détection du type de fichier
        file_extension = file.filename.lower().split(".")[-1]

        if file_extension in ["pdf"]:
            # Ouvrir le fichier PDF avec PyMuPDF
            doc = fitz.open("pdf", contents)

            if is_text_pdf(doc):
                print("Text PDF detected")
                # Extraction directe du texte avec PyMuPDF
                text = ""
                for page_num in range(len(doc)):
                    page = doc.load_page(page_num)
                    text += page.get_text("text")
                return {"text": text.strip()}
            else:
                # Conversion en images et OCR avec PaddleOCR
                images = convert_from_bytes(contents)
                text = ""
                for page_number, image in enumerate(images, start=1):
                    if image.mode in ("RGBA", "LA") or (image.mode == "P" and "transparency" in image.info):
                        image = image.convert("RGB")
                    image_np = np.array(image)
                    result = ocr.ocr(image_np, cls=True)
                    page_text = ""
                    for line in result:
                        page_text += ('\n' + line[1][0])
                    text += f"\n\n--- Texte de la page {page_number} ---\n" + page_text.strip()
                    image.close()
                    gc.collect()
                return {"text": text.strip()}

        elif file_extension in ["png", "jpeg", "jpg"]:
            # Traiter les images (PNG, JPEG) directement avec PaddleOCR
            image = Image.open(io.BytesIO(contents))

            # Convertir l'image en mode RGB si nécessaire
            if image.mode in ("RGBA", "LA") or (image.mode == "P" and "transparency" in image.info):
                image = image.convert("RGB")

            # Convertir l'image PIL en tableau NumPy
            image_np = np.array(image)

            # Effectuer l'OCR avec PaddleOCR
            result = ocr.ocr(image_np, cls=True)

            text = ""
            # Extraire uniquement le texte des résultats
            for line in result:
                text += ('\n' + line[1][0])

            image.close()

            return {"text": text.strip()}

        else:
            raise HTTPException(status_code=400, detail="Type de fichier non supporté")

    except Exception as e:
        print(f"Error: {e}")
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host='0.0.0.0', port=5000)
