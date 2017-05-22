/*
Package httpserver runs the local HTTP server for a client
this is the presentation layer in the three layered structure
of the Supply Chain Visibility Project.
*/
package main

import(
    "fmt"
    "net/http"
    "html/template"
    //"strings"
    customer "./customer"
    //customer"./struct.go"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/mgo.v2/bson"
    "os"
)

//Date is composed of Year, Month and Day
//not much I can say about that.
var databaseAddr string
var database string
var collection string
var indexServerAddr string

type Date struct{
    Year int
    Month int
    Day int
}

//Request contains the information necessary
//to query a shipment from the database.
type Request struct{
    //ID to be searched.
    ID string
    //Suppliers filter.
    Suppliers []string
    //Carriers filter.
    Carriers []string
    //Shipment's origin filter.
    FromState []string
    //Shipment's destination filter.
    ToState []string
    //Lower bound for the shipment date placement filter.
    StartDate Date
    //Upper bound for the shipment date placement filter.
    EndDate Date
}//end Request

//loadPageInfo contains information necessary to load
//the page for thr first time.
type loadPageInfo struct{
    //Carriers will contain a list of
    //all the carriers present in any of the unfinished
    //shipments stored in the user's database.
    Carriers []string
    //Suppliers will contain a list of
    //all the suppliers present in any of the unfinished
    //shipments stored in the user's database.
    Suppliers []string
    //List of states of origin of all the unfinished shipments
    //stored in the user's database.
    Origins []string
    //List of destination states of all the unfinished shipments
    //stored in the user's database.
    Destinations []string
    //List of all the unfinished shipments in which the user
    //plays the role of Customer.
    Shipments []customer.Order
    //List of all the unfinished shipments in which the user
    //plays the role of Carrier.
    CarrierShipments []customer.Order
    //List of all the unfinished shipments in which the user
    //plays the role of Supplier.
    SupplierShipments []customer.Order
}//end loadPage

//ordersHandler will load and display the page that will list all of
//the shipments stored in the user's database and handle query requests
//issued by the browser.
func ordersHandler(w http.ResponseWriter, r *http.Request, ){
    fmt.Println(r.Method)
    if r.Method == "GET" {
        fmt.Println("loading page...")
        var orders loadPageInfo

        //loading all the data to be displayed on the page.
        orders.Shipments = customer.GetUnfinishedOrder(databaseAddr, database, "Customer")
        orders.CarrierShipments = customer.GetUnfinishedOrder(databaseAddr, database, "Carrier")
        orders.SupplierShipments = customer.GetUnfinishedOrder(databaseAddr, database, "Supplier")
        orders.Suppliers = customer.GetSupplierList([]string{"Customer","Supplier","Carrier"},databaseAddr,database)
        orders.Carriers = customer.GetCarrierList([]string{"Customer","Supplier","Carrier"},databaseAddr,database)
        orders.Origins = customer.GetOrigine([]string{"Customer","Supplier","Carrier"},databaseAddr,database)
        orders.Destinations = customer.GetDest([]string{"Customer","Supplier","Carrier"},databaseAddr,database)

        //parsing the HTML template from the local resources.
        t := template.New("page")
        t = t.Funcs(template.FuncMap{"makeHex":func (v bson.ObjectId) string {return v.Hex()}})
        t.ParseFiles("templates/page.gohtml")
        t.ExecuteTemplate(w, "page", orders)
        fmt.Println("page loaded!")
        //Writing the HTML to the browser?? -- not so sure about this comment

    }else if r.Method == "POST" {
        fmt.Println("Processing request...")
        var msg Request

        //Reading the request Body.
        //The Body should be a JSON string.
        fmt.Println("reading request body...")
        message, err := ioutil.ReadAll(r.Body)
        if err != nil {
            panic(err)
        }

        //Converting the JSON string to an object.
        // The object should be of type Request.
        json.Unmarshal(message, &msg)
        fmt.Println("getting response...")

        response := customer.GetConditionalOrder(msg.Suppliers, msg.Carriers, msg.FromState, msg.ToState, msg.StartDate.Year, msg.StartDate.Month, msg.StartDate.Day, msg.EndDate.Year, msg.EndDate.Month, msg.EndDate.Day, "mycc.cit.iupui.edu", "fedex",  "Customer")
        tp := customer.GetConditionalOrder(msg.Suppliers, msg.Carriers, msg.FromState, msg.ToState, msg.StartDate.Year, msg.StartDate.Month, msg.StartDate.Day, msg.EndDate.Year, msg.EndDate.Month, msg.EndDate.Day, "mycc.cit.iupui.edu", "fedex",  "Supplier")
        response = append(response,tp...)
        tp = customer.GetConditionalOrder(msg.Suppliers, msg.Carriers, msg.FromState, msg.ToState, msg.StartDate.Year, msg.StartDate.Month, msg.StartDate.Day, msg.EndDate.Year, msg.EndDate.Month, msg.EndDate.Day, "mycc.cit.iupui.edu", "fedex",  "Carrier")
        response = append(response,tp...)
        test := response[0].ID.Hex()
        fmt.Println(test)

        b, err := json.Marshal(response)
        if err != nil {
            fmt.Println("failed to get response!")
            panic(err)
        }
        fmt.Println("response sent!")
        fmt.Fprintf(w, "%s", b)
    }//end elseif
}//end ordersHandler

//orderHandler will load and display a page with detailed information about
//a specific unfinished shipment stored in the user's local database.
func orderHandler(w http.ResponseWriter, r *http.Request){
    //text contains the order number and the client's role in the shipment.
    //This information is extracted from the URL path.
    orderID := r.URL.Path[len("/orders/order/"):]
    //textarr will contain the split text string.

    //textarr := strings.Split(text, "-")   ------
    //orderNum := textarr[0]
    //collection := textarr[1]                          ----commented on 9/14/2016

    //querying the database for the selected shipment
    var order [3]customer.Order
    var temp customer.Order
    order[0] = customer.Get(orderID, databaseAddr, database, "Customer")
    order[1] = customer.Get(orderID, databaseAddr, database, "Supplier")
    order[2] = customer.Get(orderID, databaseAddr, database, "Carrier")

    for i:= 0; i < 3; i++{
        if len(order[i].CustomerName) != 0{
                temp = order[i]
            }
    }

    t := template.New("order")

    t = t.Funcs(template.FuncMap{"makeHex":func (v bson.ObjectId) string {return v.Hex()}})
    t.ParseFiles("templates/order.gohtml")
    t.ExecuteTemplate(w, "order.gohtml", temp)
    if(r.Method == "POST"){

    }
}//end orderHandler


func main(){
    
   // databaseAddr = "mycc.cit.iupui.edu"
    //database = "fedex"
    //collection = "Customer"
    //listening on the port for requests
     if(len(os.Args)<5){
        fmt.Println("not enough arguments ")
        return ;
     }
     if(len(os.Args)>5){
        fmt.Println("too many arguments  ")
        return ;
     }

      databaseAddr = os.Args[1]
     database = os.Args[2]
     collection = os.Args[3]
     indexServerAddr = os.Args[4]

    go customer.FindInfo(databaseAddr,database,indexServerAddr)
    go customer.ClientListen(":9998",databaseAddr,database)
    go customer.ListenGPS(":4444",databaseAddr,database,indexServerAddr)

    //creating static file servers
    http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("stylesheets"))))
    http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))

    //defining the handlers for the different routes
    http.HandleFunc("/orders/", ordersHandler)
    http.HandleFunc("/orders/order/", orderHandler)

    //Running the HTTP server
    http.ListenAndServe(":8889", nil)
    //somthing(c,d)

}
//func somthing(a string, b string){}