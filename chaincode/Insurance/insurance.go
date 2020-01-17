package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"strings"
	"time"
)

//SimpleAssets implements a simple chaincode to manage an asset
type SimpleAssets struct {

}

type insurance struct {
	ObjectType string `json:"docType"`
	name string `json:"name"`
	email string `json:"email"`
	phno int `json:"phno"`
	Owner string `json:"Owner"`
}

//Init initializes chaincode
// ============================================
func (t *SimpleAssets) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

//Init is called during chaincde instantiation to initialize any data.
//Chaincode upgrade also call this function to reset or to migrate data
//func (t *SimpleAssets) Init(stub shim.ChaincodeStubInterface) peer.Response {
	//Getting the args from the trasaction proposal
//	args := stub.GetStringArgs()
//	if len(args) != 2 {
//		return shim.Error("Incorrect arguments. Expecting a key and a value")
//	}
	//Setting up  any variables or assets here by calling stub.PutState()
	//We store the key and value in the ledger
//	err := stub.PutState(args[0], []byte(args[1]))
//	if err != nil {
//		return shim.Error(fmt.Sprintf("Failed to create asset: %s",args[0]))
//	}
//	return shim.Success(nil)
//}

//Invoke - Our entry point for Invocations
//========================================
func (t *SimpleAssets) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function , args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke is running"+ function)

	//Handle different functions
	if function == "initInsurance" {
		return t.initInsurance(stub, args)
	}else if function == "transferInsurance" {
		return t.transferInsurance(stub, args)
	}else if function == "transferPolicyBasedOnType"{
		return t.transferPolicyBasedOnemail(stub,args)
	}else if function == "delete" {
		return t.delete(stub,args)
	}else if function == "readPolicy" {
		return t.readPolicy(stub, args)
	}else if function == "queryPolicyByOwner" {
		return t.queryPolicyByOwner(stub, args)
	}else if function == "getHistoryForPolicy" {
		return t.getHistoryForPolicy(stub, args)
	}else if function == "queryPolicy"  {
		return t.queryPolicy(stub, args)
	}
	fmt.Printf("invoke did not find func:"+ function)
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleAssets) initInsurance(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	var err error
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	//=======================
	//Input SanitationgetHistoryForPolicy
	fmt.Println("- start init policy")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	policyName := args[0]
	email := strings.ToLower(args[1])
	OwnerName := strings.ToLower(args[3])
	phno, err := strconv.Atoi(args[2])
	if err !=nil  {
		return shim.Error("2nd argument must be a numeric string ")
	}
	//======Checking if insurance already exits.
	policyAsByte, err :=stub.GetState(policyName)
	if err != nil {
		return shim.Error("Failed to get policy:" + err.Error())
	}else if policyAsByte != nil{
		fmt.Println("This policy already exists:" + policyName)
		return shim.Error("This policy already exists:" + policyName)
	}

	//==============Creating policy object and marshal to JSON
	objectType := "Insurance"
	insurance  := &insurance{objectType,policyName,email,phno,OwnerName}
	insuranceJSONasBytes, err :=json.Marshal(insurance)
	if err != nil{
		return shim.Error(err.Error())
	}
	err= stub.PutState(policyName,insuranceJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Indexing the Insurance to enable email-based range queries.eg:All of ram Insurance
	//An index is the normal key/value entry in the state
	//The key is composite key, with the elements that we want to range query on the listed first
	//In our case, the composite key is based on the index name
	//This will enable very efficient state range queries based on composite keys matching index-name
	indexName := "email"
	emailIndexKey , err := stub.CreateCompositeKey(indexName,[]string{insurance.email,insurance.name})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(emailIndexKey,value)

	//===Insurance saved and indexed ====//
	fmt.Println("-end init Insurance")
	return shim.Success(nil)
}

//Invoke is called per transaction on the chaincode.
//Each transaction is either a 'get' or a 'set' on the asset created by Init function.
//The 'set' method may create a new asset by specifying a new key-value pair.
//func (t *SimpleAssets) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//Extract the function and args from the transaction proposal
//	fn, args :=stub.GetFunctionAndParameters()
//	var result string
//	var err error
//	if fn == "set" {
//		result, err = set(stub, args)
//	} else {
//		result, err = get(stub, args)
//	}
//	if err != nil {
//		return shim.Error(err.Error())
//	}
	//Return the result as success payload
//	return shim.Success([]byte(result))
//}

//Set store the asset(both key and value) on the ledger.
//If the key exists,it will override the value with new one.




//readPolicy - read the policy from the chaincode state
func (t *SimpleAssets) readPolicy(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	var name, jsonResp string
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.Expecting name of the policy to query")
	}
	name = args[0]
	valAsbytes, err := stub.GetState(name)
	if err!=nil {
		jsonResp = "{\"Error\":\"Failed to get state for" + name + "\"}"
		return shim.Error(jsonResp)
	}else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Policy doesn't exist:" + name + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

func (t *SimpleAssets) delete(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	var jsonResp string
	var policyJSON insurance
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.Expecting only one")
	}
	policyName := args[0]
	valAsbytes, err := stub.GetState(policyName)
	if err != nil {
		jsonResp = "{\"Error\":\"failed to get state for" + policyName + "\"}"
		return shim.Error(jsonResp)
	}else if valAsbytes == nil{
		jsonResp = "{\"Error\":\"Policy doesn't exist:"+policyName+"\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal([]byte(valAsbytes), &policyJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"failed to decode JSON of: " + policyName + "\"}"
		return shim.Error(jsonResp)
	}
	err = stub.DelState(policyName)
	if err != nil {
		return shim.Error("Failed to delete state:"+err.Error())
	}
	//maintaining the index
	indexName := "email"
	emailIndexkey, err := stub.CreateCompositeKey(indexName , []string{policyJSON.email,policyJSON.name})
	if err != nil{
		return shim.Error(err.Error())
	}

	//delete index entry to state.
	err = stub.DelState(emailIndexkey)
	if err != nil {
		return shim.Error("Failed to delete state:"+ err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleAssets) transferInsurance(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) !=2 {
		shim.Error("Incorrect number of arguments.Expecting 2")
		policyName := args[0]
		newOwner := args[1]
		fmt.Println("Start Transferable",policyName,newOwner)
		policyAsByte, err :=stub.GetState(policyName)
		if err!= nil {
			shim.Error("Failed to get the policy"+err.Error())
		} else if policyAsByte == nil{
			shim.Error("Policy doesn't exist")
		}

		policyTotransfer := insurance{}
		err =json.Unmarshal(policyAsByte, &policyTotransfer)
		if err != nil {
			shim.Error(err.Error())
		}
		//Changing the owner
		policyTotransfer.Owner = newOwner
		policyJSONasByte, _ := json.Marshal(policyTotransfer)
		err = stub.PutState(policyName,policyJSONasByte)
		if err != nil {
			shim.Error(err.Error())			
		}
		fmt.Println("- end transfer policy(success)")
	}
	return shim.Success(nil)
}

func (t *SimpleAssets) queryPolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		shim.Error("Incorrect number of arguments.Expecting 1")
	}
	queryString := args[0]
	queryResults, err := getQueryResultforQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())		
	}
	return shim.Success(queryResults)

}

func (t *SimpleAssets) transferPolicyBasedOnemail(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments .Expecting 2")
	}
	email := args[0]
	newOwner := strings.ToLower(args[1])
	fmt.Println("- start transferPolicyBasedOnEmail", email, newOwner)

	//Query the email~name index by email
	//This will execute the key range query on all keys starting with email
	emailPolicyResultsIterator, err :=stub.GetStateByPartialCompositeKey("email~name", []string{email})
	if err != nil {
		return shim.Error(err.Error())
	}
	//Iterate through result set and for each email found,transfer to new owner
	var i int
	for i = 0; emailPolicyResultsIterator.HasNext(); i++ {
		// we don't get the value,we will just get the policy name from the composite key
		responseRange, err :=  emailPolicyResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		objectType, CompositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedemail := CompositeKeyParts[0]
		returnedPolicyName := CompositeKeyParts[1]
		fmt.Println("- found a policy from index:%s type:%s name:%s\n", objectType, returnedemail, returnedPolicyName)
		//Calling the transfer function for the found policy
		//Reusing the same function that is used to transfer individual policy
		response := t.transferInsurance(stub, []string{returnedPolicyName, newOwner})
		if response.Status != shim.OK {
			return shim.Error("Transfer failed" + response.Message)
		}
	}
	responsePayLoad := fmt.Sprintf("Transferred #{i} #{email} insurance to #{newOwner}")
	fmt.Println(" - end transferpolicyBasedonEmail" + responsePayLoad)
	return shim.Success([]byte (responsePayLoad))
}

func (t *SimpleAssets) queryPolicyByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.Expecting 1")
	}
	owner := strings.ToLower(args[0])
	queryString:= fmt.Sprintf("{\"selector\":{\"docType\":\"insurance\",\"owner\":\"%s\"}}", owner)
	queryResults, err := getQueryResultforQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (t *SimpleAssets) getHistoryForPolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		shim.Error("Incorrect number of arguments.Expecting 1")
	}
	policyName := args[0]

	fmt.Printf("- start getHistoryforPolicy: #{policyName}\n")
	resultsIterator, err := stub.GetHistoryForKey(policyName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	//buffer is a JSON array containing historic value for the policy
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext(){
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Value\":")
		if response.IsDelete{
			buffer.WriteString("null")
		}else {
			buffer.WriteString(string(response.Value))
		}
		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- getHistoryForPolicy returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func getQueryResultforQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Println("- getQueryResultforQueryString queryString:\n#{queryString}\n")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err!=nil{
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil,err
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes() ,nil

}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	//buffer is the JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse,err := resultsIterator.Next()
		if err != nil {
			return nil,err
		}
		//Adding a comma before array members
		//Suppressing it for first array members
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(",\"Record\":")
		//Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}






//func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
//	if len(args) != 2 {
//		return "",fmt.Errorf("Incorrect arguments. Expecting a key and a value")
//	}
//	err := stub.PutState(args[0], []byte(args[0]))
//	if err != nil {
//		return "", fmt.Errorf("Failed to set assests:%s",args[0])
//	}
//	return args[1], nil
//}

//Get returns the value of the specified asset key
//func get(stub shim.ChaincodeStubInterface, args []string) (string,error){
//	if len(args) != 1 {
//		return "",fmt.Errorf("Incorrect arguments.Expecting a key")
//	}
//	value, err :=stub.GetState(args[0])
//	if err != nil {
//		return "",fmt.Errorf("Failed to get asset: %s and error :%s",args[0],err)
//	}
//	if value == nil {
//		return "",fmt.Errorf("Asset not found: %s",args[0])
//	}
//	return string(value), nil
//}

//main function starts up the chaincode in the container during instantiate
func main()  {
	if err := shim.Start(new(SimpleAssets)); err != nil {
		fmt.Printf("Error starting simpleAsset chaincode : %s", err)
	}
}


