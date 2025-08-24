from fastapi import FastAPI,UploadFile,File,Form,Depends
from fileProcessing import validateFile,extractResume
from vectorDB import addingResume,matchingAlgorithm,candidateDetails
from base import JobInput,resumeMatching
from auth import verify_jwt

app = FastAPI()

@app.post("/resume-processing/{jobID}/",status_code=201)
async def uploadResume(jobID:str,file: UploadFile = File(...),
                       applicantName: str = Form(...),
                       resumeFile:str = Form(...),
                    ):

    fileData,_ = await validateFile(file=file)
    extractedString = await extractResume(file=fileData)
    response = await addingResume(resumeStr=extractedString,jobID=jobID,
                                  applicantName=applicantName,
                                  resumeFileName=resumeFile)

    return response

@app.post("/get-applicants/{jobID}/",status_code=200)
async def getCandidates(jobID,
                        job:JobInput,
                        _: dict = Depends(verify_jwt)):
    candidates = await matchingAlgorithm(
        jobDesc=job.jobDescription,
        jobID=jobID,
        topNumber=job.topNumber
    )
    return candidates

@app.post("/matching-process/",status_code=200)
async def getMatchingResume(resume:resumeMatching,
                            _: dict = Depends(verify_jwt)):
    resumeProcessing = await candidateDetails(jobDescription=resume.jobDescription,
                                              resumeName=resume.resumeName)
    return resumeProcessing

# @app.post("/matching-process/",status_code=200)
# async def getMatchingResume(resume:resumeMatching,
#                             _: dict = Depends(verify_jwt)):
#     print(resume.jobDescription,resume.resumeName)
#     return {"message":"okay"}