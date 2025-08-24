package main
import(
	"fmt"
	"net/http"
	"bytes"
	"io"
	"encoding/json"
	"os"

)

// func getResumeFile(c *gin)

func generateSignedURL(fileName string) (string, error) {
	projectID := os.Getenv("PROJECT_ID")
	bucket := os.Getenv("BUCKET")
	token := os.Getenv("SUPABASE_TOKEN")
    signURL := fmt.Sprintf(
        "https://%s.supabase.co/storage/v1/object/sign/%s/%s",
        projectID,
        bucket,
        fileName,
    )

    reqBody := `{"expiresIn": 3600}` 

    req, err := http.NewRequest("POST", signURL, bytes.NewBuffer([]byte(reqBody)))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var result map[string]interface{}
    json.Unmarshal(body, &result)

    if signedPath, ok := result["signedURL"].(string); ok {
        return fmt.Sprintf("https://%s.supabase.co/storage/v1/%s", projectID, signedPath), nil
    }

    return "", fmt.Errorf("failed to generate signed URL: %s", body)
}
