from fitz import open
from fastapi import UploadFile,File,HTTPException

async def validateFile(file:UploadFile):
    allowed_extension = {"pdf"}
    max_file_size = 5 * 1024 * 1024 
    extension = file.filename.rsplit(".",1)[-1].lower()
    if extension not in allowed_extension:
        raise HTTPException(status_code=400,detail="unsupported file extension")
    content = await file.read()
    if len(content) > max_file_size:
        raise HTTPException(status_code=400,detail="file too large (maximum 5MB)")
    return content,file.filename

async def extractResume(file:bytes):
    pdfData = open(stream=file,filetype="pdf")
    extracted_text = ""
    for page in pdfData:
        extracted_text += page.get_text()

    return extracted_text