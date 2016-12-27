package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"
)

type SearchService interface{
    // method to get the rest url
    get_rest_url() string
    
    //method to get search result
    get_search_result(channel chan string)

}


type GoogleService struct{
    // query string to hold query
    query string
    
}

func (google_service GoogleService)get_rest_url() string{
    // method to get google api url

    var google_url *url.URL
    google_url_str := "https://www.googleapis.com/customsearch/v1"
    google_api_key := "AIzaSyA9ZagpXQIfiAIRovz_qAR7-PSxh7tjsNA"
    google_cx_param := "017576662512468239146:omuauf_lfve"
    
    google_url, err := url.Parse(google_url_str)
    
    if err != nil{
        panic("some error ocurred while parsing ur")    
    }
    parameters := url.Values{}
    parameters.Add("key", google_api_key)
    parameters.Add("cx", google_cx_param)
    parameters.Add("q", google_service.query)
    google_url.RawQuery = parameters.Encode()
    return google_url.String()

}

func (google_service GoogleService)get_search_result(channel chan string){
    // error message to handle errors
    var errorMsg string 

    // rest url for google service
    restUrl := google_service.get_rest_url()

    // 1 sec timeout for hhtp client
    timeout := time.Duration(1* time.Second)
    // client for http request
    client := &http.Client{Timeout: timeout}
    
    // get the response from google service
    resp, getErr := client.Get(restUrl)
    googleServiceResult := map[string]string{
            "google": "",
            }
    
    serviceResultMap := map[string]string{
            "error": errorMsg,
            "url": restUrl,
        }

    if getErr !=nil {
        errorMsg = "Some error occured while getting the data"
    }else{
        // parse response
        htmlData, err := ioutil.ReadAll(resp.Body)
        if err !=nil {
            errorMsg = "some error ocurred while reading data"
        }
        
        serviceResultMap["text"] = string(htmlData)
    }

    serviceResultMap["error"] = errorMsg
    //jsonify data
    jsonServiceResult, _ := json.Marshal(serviceResultMap)
    googleServiceResult["google"] = string(jsonServiceResult)
    // jsonify google data
    jsonGoogleServiceResult, _ := json.Marshal(googleServiceResult)

    //sernd serach result through channel
    channel <- string(jsonGoogleServiceResult)
}


type DuckDuckGoService struct{
    // query string to hold search query
    query string   
}

func (duckduckgo_service DuckDuckGoService)get_rest_url() string{
    // method to get duckduckgo url api
    var duckduckgoUrl *url.URL
    duckduckgoUrlStr := "https://api.duckduckgo.com/"

    duckduckgoUrl, err := url.Parse(duckduckgoUrlStr)
    
    if err != nil{
        panic("some error ocurred while parsing url")    
    }
    parameters := url.Values{}
    parameters.Add("q", duckduckgo_service.query)
    parameters.Add("format", "json")
    duckduckgoUrl.RawQuery = parameters.Encode()
    return duckduckgoUrl.String()
}

func (duckduckgo_service DuckDuckGoService)get_search_result(channel chan string){
    // error message string for handling error
    var errorMsg string 
    errorMsg = ""
    
    // rest url for calling service api
    restUrl := duckduckgo_service.get_rest_url()

    // 1 sec timeout for http client
    timeout := time.Duration(1* time.Second)
    
    // initialize client
    client := &http.Client{Timeout: timeout}
    
    // duckduckgo service result map 
    duckduckgoServiceResult := map[string]string{
            "duckduckgo": "",
            }
    // service result map
    serviceResultMap := map[string]string{
            "error": errorMsg,
            "url": restUrl,
        }
    
    // get the data from client
    resp, getErr := client.Get(restUrl)
    if getErr !=nil {
         errorMsg = "Some error occured while getting the data"

    }else{
        // parse html data for response
        htmlData, err := ioutil.ReadAll(resp.Body)
        if err !=nil {
            errorMsg = "some error ocurred while reading data"
        }
        // update service result map
        serviceResultMap["text"] = string(htmlData)
            
    }

    serviceResultMap["error"] = errorMsg
    // jsonify data
    jsonServiceResult, _ := json.Marshal(serviceResultMap)
    duckduckgoServiceResult["duckduckgo"] = string(jsonServiceResult)
    // jsonify data 
    jsonDuckduckgoServiceResult, _ := json.Marshal(duckduckgoServiceResult)

    // send the search result through channel
    channel <- string(jsonDuckduckgoServiceResult)

}

func handler(w http.ResponseWriter, r *http.Request){
    // handler for handling the server request

    // get query for url 
    q:= r.URL.Query().Get("q")
    
    // response map that is used to send the response
    var responseMap = make(map[string]string)

    responseMap["query"] = q
    
    // if q is not empty then procces the query 
    // for different service
    if q != ""{
        // initializing google service
        google_service := GoogleService{q}

        // initializing duckduckgo service
        duckduckgo_service := DuckDuckGoService{q}

        // create channel for syncronization
        channel := make(chan string)

        // call goroutine to get the search result for go service
        go google_service.get_search_result(channel)
        
        // call goroutine to get the search result for go service
        go duckduckgo_service.get_search_result(channel)
        
        // get search result for servies
        google_result, duckduckgo_result := <- channel, <-channel
        
        result_str := google_result + duckduckgo_result
        responseMap["result"] = result_str
    }
    //generate json response
    jsonResponse, _ := json.Marshal(responseMap)
     
    fmt.Fprintf(w, "%s", string(jsonResponse))
}

func main(){
    // main from where server is initialized

    // registering hander for http server
    http.HandleFunc("/", handler)

    // updating server port 
    http.ListenAndServe(":8000", nil)
}
