from pydantic import BaseModel,Field

class resumeData(BaseModel):
    resumeString : str
    applicantName : str
    resumeURL: str
    

class JobInput(BaseModel):
    jobDescription: str
    topNumber: int

class responseFormatter(BaseModel):
    experience: list[str] = Field(
        description="Experience-related requirements for the role (years of experience, job titles like 'Accounting Supervisor', preferred industries, certifications, or work environment conditions)"
    )
    keywords: list[str] = Field(
        description="Important technical or operational keywords from the job description (e.g., 'cost accounting', 'IFRS', 'profitability')"
    )

class aiEvaluation(BaseModel):
    response: str = Field(description="A structured explanation of how well the candidate's resume matches the job description. "
                    "This should include specific references to matched experience, missing qualifications, "
                    "and an overall suitability assessment based on the job requirements.")
    
class resumeMatching(BaseModel):
    jobDescription: str
    resumeName: str
