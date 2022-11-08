package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	numberspb "github.com/woodman231/api_dialing_grpc/protos/numberspb"
	stringspb "github.com/woodman231/api_dialing_grpc/protos/stringspb"
)

var (
	strings_server_host = flag.String("strings_server_host", "localhost", "The host name for the strings_server service")
	strings_server_port = flag.Int("strings_server_port", 50051, "The port for the strings_server service")
	numbers_server_host = flag.String("numbers_server_host", "localhost", "The host name for the numbers_server service")
	numbers_server_port = flag.Int("numbers_server_port", 50052, "The port for the numbers_server port")
)

type APIService struct {
	StringsServiceClient stringspb.StringServiceClient
	NumbersServiceClient numberspb.NumberServiceClient
	MuxServer            *http.ServeMux
}

func main() {
	flag.Parse()

	var apiService APIService

	stringsServerConnectionString := fmt.Sprintf("%v:%v", *strings_server_host, *strings_server_port)
	numbersServerConnectionString := fmt.Sprintf("%v:%v", *numbers_server_host, *numbers_server_port)

	log.Printf("stringsServerConnectionString: %v", stringsServerConnectionString)
	log.Printf("numbersServerConnectionString: %v", numbersServerConnectionString)

	// Connect to the Strings Service
	stringsServiceConnection, err := grpc.Dial(stringsServerConnectionString, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to strings grpc server")
	}

	defer stringsServiceConnection.Close()

	stringServiceClient := stringspb.NewStringServiceClient(stringsServiceConnection)
	apiService.StringsServiceClient = stringServiceClient

	// Connect to the Numbers Service
	numbersServiceConnection, err := grpc.Dial(numbersServerConnectionString, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to numbers grpc server")
	}

	defer numbersServiceConnection.Close()

	numbersServiceClient := numberspb.NewNumberServiceClient(numbersServiceConnection)
	apiService.NumbersServiceClient = numbersServiceClient

	// Create Mux Server
	muxServer := http.NewServeMux()
	apiService.MuxServer = muxServer

	apiService.RegisterRoutes()

	fmt.Println("Listening on port 8080")

	http.ListenAndServe(":8080", muxServer)
}

func (a *APIService) RegisterRoutes() {
	// Register a hello world route to be sure that the service is running
	a.MuxServer.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// Register routes for each service
	a.registerRoutesForService(stringspb.StringService_ServiceDesc)
	a.registerRoutesForService(numberspb.NumberService_ServiceDesc)
}

func (a *APIService) registerRoutesForService(serviceDesc grpc.ServiceDesc) {
	// Shorten the service name
	shortServiceName := getShortServiceName(serviceDesc.ServiceName)

	// For each method in the service description create a route for the API
	for i := range serviceDesc.Methods {
		methodName := serviceDesc.Methods[i].MethodName

		pattern := fmt.Sprintf("/api/%v/%v", shortServiceName, methodName)

		fmt.Printf("Registering: %v.%v @ %v\n", shortServiceName, methodName, pattern)

		a.MuxServer.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			a.handleRequest(shortServiceName, methodName, w, r)
		})
	}
}

func getShortServiceName(currentServiceName string) string {
	words := strings.Split(currentServiceName, ".")
	wordCount := len(words)

	if wordCount > 0 {
		return words[wordCount-1]
	}

	return words[0]
}

func (a *APIService) handleRequest(serviceName string, methodName string, w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for %v", r.URL.String())

	// Set Up Some Variables to Hold our grpcResponse, and possible grpcError
	var grpcResponse interface{}
	var grpcErr error

	switch serviceName {
	case "StringService":
		var stringOperationRequest stringspb.OperationRequest

		// Convert the HTTP Request Body to a StringService OperationRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&stringOperationRequest)

		// Validate the request
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}

		if stringOperationRequest.InputString == "" {
			http.Error(w, "Missing input_string, or input_string is empty", http.StatusBadRequest)
			return
		}

		// Use the decoded request body with the appropriate method that was passed in
		switch methodName {
		case "MakeUpperCase":
			grpcResponse, grpcErr = a.StringsServiceClient.MakeUpperCase(context.Background(), &stringOperationRequest)
		case "MakeLowerCase":
			grpcResponse, grpcErr = a.StringsServiceClient.MakeLowerCase(context.Background(), &stringOperationRequest)
		default:
			http.Error(w, "Not Implemented", http.StatusInternalServerError)
			return
		}
	case "NumberService":
		var numberOperationRequest numberspb.OperationRequest

		// Convert the HTTP Request Body to a NumbersSErvice OperationRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&numberOperationRequest)

		// Validate the request
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}

		if numberOperationRequest.InputNumberOne == 0 && numberOperationRequest.InputNumberTwo == 0 {
			http.Error(w, "Either input_number_one or input_number_two are missing or they are both 0", http.StatusBadRequest)
			return
		}

		// Use the decoded request body with the appropraite method that was passed in
		switch methodName {
		case "AddTwoNumbers":
			grpcResponse, grpcErr = a.NumbersServiceClient.AddTwoNumbers(context.Background(), &numberOperationRequest)
		case "SubtractTwoNumbers":
			grpcResponse, grpcErr = a.NumbersServiceClient.SubtractTwoNumbers(context.Background(), &numberOperationRequest)
		default:
			http.Error(w, "Not Implemented", http.StatusInternalServerError)
			return
		}
	}

	// Write the error from the GRPC Client if applicable
	if grpcErr != nil {
		log.Print("Could not get response from GRPC")
		http.Error(w, grpcErr.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response from the GRPC Client if applicable
	marshaledResponse, marshalErr := json.Marshal(grpcResponse)

	if marshalErr != nil {
		log.Printf("Could not marshal response from GRPC")
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(marshaledResponse)
}
