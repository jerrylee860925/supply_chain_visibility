package main

import (
 "log"
 "gopkg.in/mgo.v2"
 "gopkg.in/mgo.v2/bson"
 "strconv"
 "fmt"
 "time"
 "math"
 "math/rand"
 "os"
)
type indexServer struct{
  ReqType int
  Name string
  IpAddr string
  ID bson.ObjectId `bson:"_id,omitempty"`
}
type product struct {
  ProductName,ProductCode,ProductState,UnitMeasure string 
  UnitPrice,Quantity float64
}
type ProductList struct{
  ListOfProduct []product
}
type Order struct {
  ID bson.ObjectId `bson:"_id,omitempty"`
  OrderList ProductList
  CustomerName,CustomerCode,SupplierName,SupplierCode,CarrierName,CarrierCode string

  Origin, Dest string
  OrderSts OrderStatus
  PickUptime string
  ETA int
  OrderDate string
}
type OrderStatus struct {
  Status string
  Trucks string
  TimeStamp []string
  GPSCordsX []float64
  GPSCordsY []float64
}
type col struct{
     Id string
     Password string
     Uid int
     Other string
}
var m bson.M
var status = [...]string {
  "on the ship",
  "off the ship",
  "on the dock",
  "in the storage",
  "getting ready",
  "ready for pickup",
  "picked up",
  "in transit",
  "delivered",
  "finished",
}

func CalcDay(date string)int {

  timeFormat := "2006-01-02 15:04 MST"
  then, err := time.Parse(timeFormat, date)
  if err != nil {
    fmt.Println(err)
    return (-1)
    }
    duration := time.Since(then)
    day := math.Ceil(duration.Hours()/24)
    return int(day)
}

func errhandler(err error,something string) {
  if err != nil {
        log.Fatal(err)
        fmt.Println(something)
    }
}
func INSERT ( role1 string, role2 string, role3 string,databseAddr string,database string, collection string ,flag int){
      year := "2016"
      var tmp [3]Order
      var nilorder Order;
      //nilorder = nil
      var tmp4 ProductList
      var tmp5 ProductList
      tmp2 := make([]product,10,10) 
      tmp3 := make([]product,10,10)
      for i:=0;i<10;i++{
        tmp2[i].ProductCode = strconv.Itoa((i+1)*rand.Intn(1700))
        tmp2[i].ProductName = strconv.Itoa((i+1)*rand.Intn(200))
        tmp2[i].ProductState = "liquid"
        tmp2[i].UnitMeasure = "liter"
        tmp2[i].UnitPrice = float64(i+1)*(1.7)
        tmp2[i].Quantity  = float64(i+1)*(2)
      }
       for i:=0;i<10;i++{
        tmp3[i].ProductCode = strconv.Itoa((i+1)*rand.Intn(3000))
        tmp3[i].ProductName = strconv.Itoa((i+1)*rand.Intn(400))
        tmp3[i].ProductState = "gas"
        tmp3[i].UnitMeasure = "gallon"
        tmp3[i].UnitPrice = float64(i+1)*(1.7)
        tmp3[i].Quantity  = float64(i+1)*(2)
      }
      tmp4.ListOfProduct =tmp2 
      tmp5.ListOfProduct =tmp3
      tmp[0].OrderList =  tmp4
      tmp[1].OrderList =  tmp5
      tmp[2].OrderList =  tmp5
     
      for i:=0;i<3;i++{


        tmp[i].CustomerName = role1
        tmp[i].CustomerCode = "57dafad59b8f13c51c80ad45"
        tmp[i].SupplierName = role2
        tmp[i].SupplierCode = "57dafad59b8f13c51c80ad46"
        tmp[i].CarrierName = role3
        tmp[i].CarrierCode =  "57dafad59b8f13c51c80ad47"

      if i%2==0{
        tmp[i].OrderSts.Status = status[0]

      }else{
        tmp[i].OrderSts.Status  = status[0]
      }

        month := strconv.Itoa(i%5+1)
        day := strconv.Itoa((i*13)%30+1)
        total := year+"-"+month+"-"+day+" 15:04 MST"
        tmp[i].OrderDate=total


    }
    if flag==1{
          tmp[0].Origin = "IL"
          tmp[1].Origin = "IN"
          tmp[2].Origin = "CA"
          tmp[0].Dest = "ID"
          tmp[1].Dest = "IA"
          tmp[2].Dest = "KS"
    }else if flag==2{
          tmp[0].Origin = "AL"
          tmp[1].Origin = "AZ"
          tmp[2].Origin = "CO"
          tmp[0].Dest = "KS"
          tmp[1].Dest = "KY"
          tmp[2].Dest = "LA"
    }else if flag==3{
          tmp[0].Origin = "DE"
          tmp[1].Origin = "FL"
          tmp[2].Origin = "GA"
          tmp[0].Dest = "MI"
          tmp[1].Dest = "MN"
          tmp[2].Dest = "NY"
    }
    tmp[0].ID = bson.ObjectIdHex("57dafad69b8f13c51c80ad48");
    tmp[1].ID = bson.ObjectIdHex("57dafad69b8f13c51c80ad4b");
    tmp[2].ID = bson.ObjectIdHex("57dafad69b8f13c51c80ad49");
    insertClient(nilorder,"Supplier",database,databseAddr)
    insertClient(nilorder,"Carrier",database,databseAddr)
    insertClient(nilorder,"Customer",database,databseAddr)

    for i := 0; i < 3; i++ {
         // clientPut(tmp[i],flag,"1")
        insertClient(tmp[i],collection,database,databseAddr)

    }
}


func insertClient(in Order,collection string,database string,databseAddr string){
     session, err := mgo.Dial(databseAddr)
     errhandler(err,"connection")
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)
     d := session.DB(database).C(collection)
     err = d.Insert(&in)
     errhandler(err,"db")     
}

func dropCol(collectionName string,dbName string,databseAddr string  ){
     session, err := mgo.Dial("localhost")
     errhandler(err,"connection")
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)
    
     d := session.DB(dbName).C(collectionName)
     err = d.DropCollection()
     errhandler(err,"db")
}
func insertInfo(databseAddr string, database string, IpAddr string){
    var in indexServer;
    in.ReqType = 1;
    in.Name = database;
    in.IpAddr = IpAddr;
    session, err := mgo.Dial(databseAddr)
    errhandler(err,"connection")
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    d := session.DB(database).C("Info")
    err = d.Insert(&in)
    errhandler(err,"db")
}
func collectionExist(databaseAddr string,database string) []string{
  //var result Order
  session, err := mgo.Dial(databaseAddr)
  if err != nil {
    panic(err)
  }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)
  d,_ := session.DB(database).CollectionNames()
 
  return d;
}

func main(){

  roll1 := os.Args[1]
  roll2 := os.Args[2]
  roll3 := os.Args[3]
  databseAddr := os.Args[4]
  database := os.Args[5]
  collection := os.Args[6]
  IpAddr := os.Args[7]//local IP address

  collections := collectionExist(databseAddr,database) 
    for i := 0; i < len(collections); i++ {
       // fmt.Println(i,"  " ,len())
        dropCol(collections[i],os.Args[5],os.Args[4])
    }
  insertInfo(databseAddr , database, IpAddr)
	//dropCol(os.Args[6],os.Args[5],os.Args[4])
  //role1 string, role2 string, role3 string,databseAddr string,database string, collection string ,flag int
  INSERT(roll1,roll2,roll3,databseAddr,database,collection,1 )

 

}
