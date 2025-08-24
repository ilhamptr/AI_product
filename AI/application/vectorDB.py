from langchain_chroma import Chroma
from datetime import datetime,UTC
from uuid import uuid4
from base import responseFormatter,aiEvaluation
from langchain_core.documents import Document
from langchain_groq import ChatGroq
from langchain_cohere import CohereEmbeddings
from dotenv import load_dotenv
from fastapi import HTTPException
import json
import asyncio
import os

load_dotenv()

embedding = CohereEmbeddings(cohere_api_key=os.getenv("COHERE_API_KEY"),model='embed-english-v3.0')

vectorStore = Chroma(collection_name="resumes",
                    persist_directory="./resumes_collection",
                    embedding_function=embedding)
    
model = ChatGroq(  model="llama-3.3-70b-versatile",
    temperature=0,
    max_tokens=1000,
    # reasoning_format="parsed",
    # reasoning_effort="medium",
    timeout=None,)

async def addingResume(resumeStr:str,jobID:str,applicantName:str,resumeFileName:str):

    now = datetime.now(UTC)
    created_at = now.isoformat().replace('+00:00', 'Z')
    recordID = str(uuid4())
    document = Document(page_content=resumeStr,metadata={"applicant_name":applicantName,"job_id":jobID,"created_at":created_at,"resume_file":resumeFileName})
    vectorStore.add_documents(documents=[document],ids=[recordID])

    return {"message":"record has been added","data":{"id":recordID,"applicant_name":applicantName,"job_id":jobID,"created_at":created_at,"resume_file":resumeFileName}}


async def matchingAlgorithm(jobDesc:str,jobID:str,topNumber:int):
   

    messages = [
    (
        "system",
        "You are an expert assistant that extracts structured keywords and experience requirements from job descriptions. \
        'Experience requirements' should include years of experience, job titles (roles), industry preferences, qualifications, certifications, \
        and work conditions related to experience.",
    ),
    (
        "human",
        f"""Extract experience requirements and technical keywords from the following job description: {jobDesc}"""
    ),
]
    structuredModel = model.with_structured_output(responseFormatter)
    result = await structuredModel.ainvoke(messages)
    result
    query = f"""
    Find candidates whose experiences align with:
    Experience Requirements: {result.experience}
    Key Skills and Keywords: {result.keywords}

    """
    results = vectorStore.similarity_search_with_score(query=query, k=topNumber,filter={"job_id":str(jobID)})

    finalResult = []
    for doc, score in results:
        doc.metadata["id"] = doc.id
        doc.metadata["distance"] = score
        finalResult.append(doc.metadata)

    return {"data":finalResult}       


async def candidateDetails(jobDescription:str,resumeName:str):
    collection =  vectorStore._collection
    result = collection.get(
        where={"resume_file":resumeName}
    )
    if not result["ids"]:
        raise HTTPException(status_code=404,detail="can't find candidate data")

    messages = [
        (
            "You are an expert assistant that compares a job description with a candidate's resume. "
            "Your task is to determine how well the candidate's experience matches the job requirements. "
            "Focus on matching specific elements such as:\n"
            "- Required years of experience\n"
            "- Previous job titles or roles\n"
            "- Relevant industries or sectors\n"
            "- Educational background and certifications\n"
            "- Familiarity with tools, frameworks, or standards\n"
            "- Work conditions or expectations (e.g., WFO, remote, location)\n"
            "Be objective, structured, and explain where the resume aligns or falls short of the job description."
        )
        ,
        (
        "human",
            f"""
       Compare the following job description and resume.

        ### Job Description:
        {jobDescription}


        ### Resume:
        {result['documents'][0]}

        """
        )

    ]

    structured_model = model.with_structured_output(aiEvaluation)

    result = await structured_model.ainvoke(messages)
    return result


# asyncio.run(candidateDetails(jobDescription=""" Responsibilities:\n\n- Contribute to the development of company products and projects based on requirements.\n- Program the computer by entering coded information.\n- Confirm program operation by conducting tests; modifying program sequence and/or codes.\n- Contribute to team effort by accomplishing related results as needed.\n- Prepare reference materials for users by writing operating instructions.\n- Maintain historical records by documenting program development and revisions.\n\nGeneral Requirements:\n\n- Candidate must possess at least 
#                              a SMU, Diploma, or Bachelor's Degree in any field.\n- Applicants must be willing to work in Bandung.\n- Willing to travel and/or be placed out of town.\n- No work experience required.\n- 2 full-time positions available.\n- Fast learner, detail-oriented, and able to work under tight timelines.\n- Positive working attitude, self-initiated, and disciplined.\n- Good communication skills and self-managed.\n- Motivated to self-study and research.\n\nTechnical Requirements:\n\n- Strong knowledge of Object-Oriented Programming (OOP).\n- Good understanding of web programming concepts (client-server, HTTP methods, sessions, cookies).\n-
#                               Good knowledge of HTML & CSS.\n- Experience with web MVC frameworks (e.g., Laravel, CakePHP, Java iBatis, Java Spring, CodeIgniter) is a plus.
#                              \n- Comfortable with databases and SQL.\n- Knowledge and experience with Python or Node.js is a plus.\n- Freelance experience is a plus.\n- Familiar with Linux is a plus.\n- Familiar with Subversion or Git is a plus.""",resumeID="1bdeb99a-2c87-48f2-bfeb-f880a8e07bf9"))




