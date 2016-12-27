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
    query string
    
}

func (google_service GoogleService)get_rest_url() string{
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
    var errorMsg string 
    restUrl := google_service.get_rest_url()
    timeout := time.Duration(1* time.Second)
    client := &http.Client{Timeout: timeout}
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
        htmlData, err := ioutil.ReadAll(resp.Body)
        if err !=nil {
            errorMsg = "some error ocurred while reading data"
        }
        
        serviceResultMap["text"] = string(htmlData)
    }
    serviceResultMap["error"] = errorMsg
    jsonServiceResult, _ := json.Marshal(serviceResultMap)
    googleServiceResult["google"] = string(jsonServiceResult)
    jsonGoogleServiceResult, _ := json.Marshal(googleServiceResult)
    channel <- string(jsonGoogleServiceResult)
}


type DuckDuckGoService struct{
    query string   
}

func (duckduckgo_service DuckDuckGoService)get_rest_url() string{
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
    var errorMsg string 
    errorMsg = ""
    restUrl := duckduckgo_service.get_rest_url()
    timeout := time.Duration(1* time.Second)
    client := &http.Client{Timeout: timeout}
    duckduckgoServiceResult := map[string]string{
            "duckduckgo": "",
            }
    serviceResultMap := map[string]string{
            "error": errorMsg,
            "url": restUrl,
        }
 
    resp, getErr := client.Get(restUrl)
    if getErr !=nil {
         errorMsg = "Some error occured while getting the data"

    }else{

        htmlData, err := ioutil.ReadAll(resp.Body)
        if err !=nil {
            errorMsg = "some error ocurred while reading data"
        }
        serviceResultMap["text"] = string(htmlData)
            
    }

    serviceResultMap["error"] = errorMsg
    jsonServiceResult, _ := json.Marshal(serviceResultMap)
    duckduckgoServiceResult["duckduckgo"] = string(jsonServiceResult)
    jsonDuckduckgoServiceResult, _ := json.Marshal(duckduckgoServiceResult) 
    channel <- string(jsonDuckduckgoServiceResult)

}

func handler(w http.ResponseWriter, r *http.Request){

    q:= r.URL.Query().Get("q")
    var responseMap = make(map[string]string)
    responseMap["query"] = q
    if q != ""{
        google_service := GoogleService{q}
        duckduckgo_service := DuckDuckGoService{q}
        channel := make(chan string)
        go google_service.get_search_result(channel)
        go duckduckgo_service.get_search_result(channel)
        google_result, duckduckgo_result := <- channel, <-channel
        result_str := google_result + duckduckgo_result
        responseMap["result"] = result_str
    }
    jsonResponse, _ := json.Marshal(responseMap)
    fmt.Fprintf(w, "%s", string(jsonResponse))
}

func main(){
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8000", nil)
}
