/*************************************************
function.go contains all function that are needed
in both http server and middleware.
 *
* @author  zhijie li
* @author  haoyan wu
* @data 09-02-2016
**************************************************/
package customer
import (
    "fmt"
    "net"
    "encoding/json"
    "log"
    "gopkg.in/mgo.v2"
    "time"
    "gopkg.in/mgo.v2/bson"
    "bufio"
    "strings"
    "strconv"
)
/**
 * @description the function handles an mongodb error, user input an err
  or the function prints the error if the error is not nil
 * @param  err the error user wants to handle
 * @return void
*/
func errhandler(err error) {
  if err != nil {
       fmt.Println(err)
        return
    }
}
/*
 * @description the function update a certain record in a particular collection of a mongodb
 * @param string databaseAddr the IP address of the server that contains the database is installed
 * @param Order old the old order record user wants to modify
 * @param Order new the new order record user wanto to replace with
 * @param string database the name of the database that contains the collection
 * @param string collection the name of the collection that contains the order record
 * @return void
*/

func update(databaseAddr string,old Order,new Order,database string,collection string ){
  fmt.Println("update collection",collection)
     session, err := mgo.Dial(databaseAddr)
     errhandler(err)
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)
     d := session.DB(database).C(collection)
     err = d.Update(old,new)
}

/*
 * @description the function searches a particular order record from given database
 * @param string databaseAddr the IP address of the server that contains the database is installed
 * @param string database the name of the database that contains the collection
 * @param string collection the name of the collection that contains the order record
 * @param string hexID the ID of the record
 * @return Order the order that is found in database, an empty order is returned if nothing found in database
*/
func findOldOrder(databaseAddr string,database string,collection string,hexID string) Order{
  var result Order
  fmt.Println(" database :" , database," collection :", collection, "hexID :   ",hexID)
  session, err := mgo.Dial(databaseAddr)
  if err != nil {
    panic(err)
  }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)
  c := session.DB(database).C(collection)
  err = c.FindId(bson.ObjectIdHex(hexID)).One(&result) //Works
  if err != nil {
    fmt.Println("find old order here errr ", err)
    }
  return result
}
/*
 * @description the function searches all unfinished orders in a given database
 * @param string databaseAddr the IP address of the server that contains the database is installed
 * @param string database the name of the database that contains the collection
 * @param string collection the name of the collection that contains the order record
 * @return []Order an arrary of order structs that has non-finished status
 */

func GetUnfinishedOrder(databaseAddr string,database string,collection string) []Order{
  var result []Order
  session, err := mgo.Dial(databaseAddr)
  if err != nil {
    panic(err)
  }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)
  c := session.DB(database).C(collection)
//  err= c.Find(nil).All(&result)
  err = c.Find(bson.M{"ordersts.status":bson.M{"$ne":"finished"}}).All(&result)
  if err != nil {
      log.Fatal(err)
  }
  return result
}
/*
 * @description the function searches a particular order record from given database
 * @param string hexID the ID of the record
 * @param string databaseAddr the IP address of the server that contains the database is installed
 * @param string database the name of the database that contains the collection
 * @param string collection the name of the collection that contains the order record
 * @return Order the order that is found in database, an empty order is returned if nothing found in database
 */

func Get(hexID string,databaseAddr string,database string,collection string) Order{
    session, err := mgo.Dial(databaseAddr)
    errhandler(err)
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)

     var result Order
     c := session.DB(database).C(collection)
     err = c.FindId(bson.ObjectIdHex(hexID)).One(&result)
     errhandler(err)
     return result
}

/*
 * @description the function connects to index server and queries about an ip address of one particular host with given ID and name and return the ip address
 * @param string the ip address of index server
 * @param string the port that index server is listening to
 * @param string ID the ID of the host user needs to query
 * @param string name the name of the host user needs to query
 * @return string ans.ipAddr the ip address of the host user needs to query. if the ID does not exist return 0.0.0.0
 */
func GettingIPAddr( IpAddr string, port string,ID string,Name string)string{
  var newReq indexServer
  newReq.ReqType = 0
  newReq.Name = Name
  newReq.ID=bson.ObjectIdHex(ID)
  c, err := net.Dial("tcp", IpAddr+":"+port)
  if err != nil {
        fmt.Println(err)
        return "0.0.0.0"
    }

  b,e := json.Marshal(newReq)
  if e != nil {
    fmt.Println(e)
    c.Close()
    return "0.0.0.0"
  }
  e1 := json.NewEncoder(c).Encode(b)
  if e1 != nil {
    fmt.Println(e1)
    c.Close()
    return "0.0.0.0"
  }
  //
  var ans indexServer
  var msg []byte
  time.Sleep(time.Second*3)
  err = json.NewDecoder(c).Decode(&msg)
  e = json.Unmarshal(msg,&ans)
  if e != nil {
    fmt.Println(e)
    c.Close()
    return "0.0.0.0"
  }
  if err != nil {
    fmt.Println(err)
    c.Close()
    return "0.0.0.0"
  }
  c.Close()
  return ans.IpAddr
}



/*
 * @description the function gets the status number of a particula
 * @param Order input the order needs to be checked
 * @return int count the status numbers
 */

func GetOrderStatusNum(input Order )int32 {
  var count int32
  count = 0;
  for input.OrderSts.Status != status[count]{
    count++
  }
  return count
}

/*
 * @description the function the function sends an order struct to an other party
 * @param string the ip address of the host user wants to send
 * @param Order mOrder the order needs to be sent out
 * @param string port the port the other host is listening to
 * @return void
 */

func sendToOther(ipaddr string,mOrder Order,Port string ){

  c, err := net.Dial("tcp", ipaddr+":"+Port)
  if err != nil {
        fmt.Println(err)
        return
  }
  b,e := json.Marshal(mOrder)
  if e != nil {
      fmt.Println(e)
      c.Close()
      return
    }
    e2 := json.NewEncoder(c).Encode(b)
    if e2 != nil {
      fmt.Println(e2)
      c.Close()
      return
    }
    c.Close()
}

/*
 * @description the function gets the roll of local client in a particular shipment(roll can be customer carrier supplier)
 * @param Order mOrder the order that the user needs to check for its roll
 * @return int the number that represents the roll (0:customer 1:carrier 2:supplier) if error happens return -1
 */
func GetRoll(mOrder Order, databaseAddr string,database string) int32{
  var result indexServer
  session, err := mgo.Dial(databaseAddr)
  if err != nil {
    fmt.Println("connection error to database ", databaseAddr)
    panic(err)
  }
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  c := session.DB(database).C("Info")
  err = c.Find(bson.M{"name":database}).One(&result) //Works
    if err != nil {
      fmt.Println("find roll err  ",err)
    }
   switch {
    case mOrder.CustomerCode == result.ID.Hex():
      return 0
    case mOrder.CarrierCode == result.ID.Hex():
      return 1
    case mOrder.SupplierCode == result.ID.Hex():
      return 2
  }
  return -1
}
/*
 * @description the function filters the all order that the client is a carrier in and the client need to update
 * @param []Order input the order list in which the client plays carrier roll
 * @return []Order the order list that the client plays carrier roll in and the client needs to update order status for.
 * if non is found then return empty array
 */
func GetCarrOrd(input []Order)[]Order{

  var CarrOrdr []Order
  //reusult:=GetUnfinishedOrder("mycc.cit.iupui.edu",Mdatabase,"Carrier")
  var statusNum int32
  for i := 0; i < len(input); i++ {
      statusNum =  GetOrderStatusNum(input[i])
      if(statusNum>=7 && statusNum<9){
        CarrOrdr= append (CarrOrdr, input[i])
      }
  }
  return CarrOrdr
}
/*
 * @description the function filters the all order that the client is a supplier in and the client need to update
 * @param []Order input the order list in which the client plays supplier roll
 * @return []Order the order list that the client plays supplier roll in and the client needs to update order status for.
 * if non is found then return empty array
 */
func GetSuppOrd(input []Order)[]Order{
  var suppOrdr   = make([]Order,len(input),len(input))
  var statusNum int32
  count :=0
  for i := 0; i < len(input); i++ {
      statusNum =  GetOrderStatusNum(input[i])
      if(statusNum<7){
          suppOrdr[count] = input[i]
          count++
      }

  }
  var res = make([]Order,count,count)
  for i := 0; i < count; i++ {
      res[i] = suppOrdr[i]
  }
  return res

}


  /*
 * @description the function listens to a certain port and wait for order updates, it will call recivemsg function once a order update is arrived
 * @param string port the port that the function needs to listen to
 * @return void
 */
func ClientListen(port string,databaseAddr string, database string) {
    ln, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("error\n")
        fmt.Println(err)
      return
    }
    for {
        nc, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        //fmt.println("123456789012345678901234567890")
      go recivemsg(nc,database,databaseAddr)
    }
}
/*
 * @description the function check an order update and update the the corresponding record in the database
 * @param net.conn the connection that contains the update message
 * @param string database the name of the databse which contains all order info
 * @return void
 */

func recivemsg(nc net.Conn,database string,databaseAddr string ){
    var msg []byte
    var nOrder Order
    fmt.Println("123456789012345678901234567890")
    err := json.NewDecoder(nc).Decode(&msg)
    errhandler(err)
    e := json.Unmarshal(msg,&nOrder)
    errhandler(e)
    fmt.Println(nOrder)
    nc.Close()
    //databaseAddr string,old Order,new Order,database string,collection string
    var collection string
    //func GetRoll(mOrder Order, databaseAddr string,database string) int32{
    rollNum := GetRoll(nOrder,databaseAddr,database)
    fmt.Println("1231231231231231231231232131232132132131231312312")

    switch rollNum{
      case  0:
        collection = "Customer"
        break;
      case 1:
        collection = "Carrier"
        break;
      case 2 :
        collection = "Supplier"
        break;
      case -1:
        fmt.Println("wrong")
        return
    }

    fmt.Println("490821219321809218439274893275895908324802")
    oldOrder :=findOldOrder(databaseAddr,database,collection,nOrder.ID.Hex())

    update(databaseAddr,oldOrder,nOrder,database,collection)

}

/*
 * @description the function connects to index server and queries about an ip address of one particular host with given ID and name and return the ip address
 * @param string the ip address of index server
 * @param string the port that index server is listening to
 * @param string ID the ID of the host user needs to query
 * @param string name the name of the host user needs to query
 * @return string ans.ipAddr the ip address of the host user needs to query. if the ID does not exist return 0.0.0.0
 */
func GetSupplierList(collection []string, databaseAddr string , database string ) []string{
  var mOrder []Order
  for i := 0; i < len(collection); i++ {
    mOrder = append(mOrder, GetUnfinishedOrder(databaseAddr,database,collection[i])...)
  }
  var SuppList []string
  index := -1
  for i := 0; i < len(mOrder); i++ {
    for j := 0; j < len(SuppList); j++ {
      if(mOrder[i].SupplierName==SuppList[j]){
        index = j;

        break;
      }
    }
    if(index == -1){
      SuppList = append(SuppList,mOrder[i].SupplierName)
    }
    index = -1
  }
  fmt.Println(SuppList,"  ",len(SuppList))
  return SuppList
}
/*
 * @description the function connects to index server and queries about an ip address of one particular host with given ID and name and return the ip address
 * @param string the ip address of index server
 * @param string the port that index server is listening to
 * @param string ID the ID of the host user needs to query
 * @param string name the name of the host user needs to query
 * @return string ans.ipAddr the ip address of the host user needs to query. if the ID does not exist return 0.0.0.0
 */
func GetCarrierList(collection []string,databaseAddr string , database string ) []string{
  var mOrder []Order
  for i := 0; i < len(collection); i++ {
    mOrder = append(mOrder, GetUnfinishedOrder(databaseAddr,database,collection[i])...)

  }
  var CarrierLst []string
  index := -1
  for i := 0; i < len(mOrder); i++ {

    for j := 0; j < len(CarrierLst); j++ {
      if(mOrder[i].CarrierName==CarrierLst[j]){
        index = j;
        break;
      }
    }
    if(index == -1){
      CarrierLst = append(CarrierLst,mOrder[i].CarrierName)
    }
    index = -1
  }
  fmt.Println(CarrierLst,"   ",len(CarrierLst))
  return CarrierLst
}
/*
 * @description the function connects to index server and queries about an ip address of one particular host with given ID and name and return the ip address
 * @param string the ip address of index server
 * @param string the port that index server is listening to
 * @param string ID the ID of the host user needs to query
 * @param string name the name of the host user needs to query
 * @return string ans.ipAddr the ip address of the host user needs to query. if the ID does not exist return 0.0.0.0
 */

func GetDest(collection []string,databaseAddr string , database string )[]string{
  var allOrder []Order//:= GetUnfinishedOrder()
  fmt.Println(len(allOrder))
   for i := 0; i < len(collection); i++ {
    allOrder = append(allOrder, GetUnfinishedOrder(databaseAddr,database,collection[i])...)

  }
  fmt.Println(len(allOrder))
  var Dest []string
  var ifExist bool
  ifExist = false
  for i := 0; i < len(allOrder); i++ {
  	fmt.Println(allOrder[i].Dest)
    for j := 0; j < len(Dest); j++ {
      if(Dest[j] == allOrder[i].Dest ){
        ifExist = true
        break;
      }
    }
    if(ifExist == false){
      Dest = append(Dest,allOrder[i].Dest)
    }
    ifExist = false

  }
  fmt.Println("DDDDDDDDDDDDDDDDDDDDDDDDDDDDD ",Dest,"   ",len(Dest))
  return Dest
}
/*
 * @description the function get
 * @param string the ip address of index server
 * @param string the port that index server is listening to
 * @param string ID the ID of the host user needs to query
 * @param string name the name of the host user needs to query
 * @return string ans.ipAddr the ip address of the host user needs to query. if the ID does not exist return 0.0.0.0
 */
func GetOrigine(collection []string,databaseAddr string , database string)[]string{
 var allOrder []Order//:= GetUnfinishedOrder()
   for i := 0; i < len(collection); i++ {
    allOrder = append(allOrder, GetUnfinishedOrder(databaseAddr,database,collection[i])...)

  }
  fmt.Println(len(allOrder))
    fmt.Println(allOrder)
  var Origin []string
  var ifExist bool
  ifExist = false
  for i := 0; i < len(allOrder); i++ {
    for j := 0; j < len(Origin); j++ {
      if(Origin[j] == allOrder[i].Origin ){
        ifExist = true
        break;
      }
    }
    if(ifExist == false){
      Origin = append(Origin,allOrder[i].Origin)
    }
    ifExist = false
  }
  fmt.Println("oooooooooooo  ",Origin,"   ",len(Origin))

  return Origin
}

func GetConditionalOrder(supplier []string, carrier []string,origine []string, dest []string, startYear int, startMonth int,startDay int, endYear int, endMonth int, endDay int,databaseAddr string,database string,collection string ) ([]Order){
  result := GetUnfinishedOrder(databaseAddr ,database ,collection)
  var res1 []Order
  var res2 []Order
  var res3 []Order
  var res4 []Order
  var res5 []Order
  fmt.Println(supplier,"aaaaaaaaaaaaa", carrier," ccccccc  ",origine ," bbbbbbbb   ", dest )
  if(supplier[0] == "any"){
    for i := 0; i < len(result); i++ {
      res1  = append(res1,result[i])
    }
  }else{

          for i:=0;i<len(result);i++{

            for j:=0;j<len(supplier);j++{
              if(result[i].SupplierName == supplier[j]){
                res1 = append(res1,result[i])
                break;
                }
            }
          }
  }
  fmt.Println("res111111111111111111111111111111111")
  fmt.Println(res1)
  if(carrier[0] == "any"){
    for i := 0; i < len(res1); i++ {
      res2  = append(res2,result[i])
    }
  }else{
  for i:=0;i<len(res1);i++{
      for j:=0;j<len(carrier);j++{

          if(res1[i].CarrierName == carrier[j]){
            res2 = append(res2,res1[i])
            break;
        }
      }
    }
  }
  fmt.Println("res2222222222222222222222222222222222222   ")
  fmt.Println(res2)
  if(origine[0] =="any"){
      for i := 0; i < len(res2); i++ {
        res3  = append(res3,res2[i])
      }
    }else{
      for i := 0; i< len(res2); i++ {
        for j:=0; j< len(origine); j++ {
            if(res2[i].Origin == origine[j]){
              res3 = append(res3,res2[i])
              break;
          }
        }
      }
    }
    fmt.Println("res3333333333333333333333333   ")
    fmt.Println(res3)
  if(dest[0] =="any"){
      for i := 0; i < len(res3); i++ {
        res4  = append(res4,res3[i])
      }
    }else{
      for i := 0; i< len(res3); i++ {
        for j:=0; j< len(dest); j++ {
            if(res3[i].Dest == dest[j]){
              res4 = append(res4,res3[i])
              break;
          }
        }
      }
    }
    fmt.Println("res444444444444444444444444444444   ")
    fmt.Println(res4)
  startdate := parseTimeStart(startYear, startMonth,startDay)
  enddate := parseTimeEnd(endYear, endMonth, endDay)
  for i := 0; i < len(res4); i++ {
    fmt.Println(res4[i].OrderDate)
    checkDate := ParseTime(res4[i].OrderDate)
    if Compare(startdate,enddate,checkDate)==true {
      fmt.Println("truetruetruetruetruetruetruetruetruetruetruetruetruetruetruetruetruetruetruetrue")
      res5= append(res5,res4[i])
    }
  }
  fmt.Println("res55555555555555555555555555555555 ")

  fmt.Println(res5)
  return res5
}

func parseTimeStart(year int,month int,day int)time.Time{
  locationTime :=time.Now()
  dateParse := time.Date(year , time.Month(month), day, 0, 0, 0, 0, locationTime.Location())
  fmt.Println("time start ",dateParse)
  return dateParse
}

func parseTimeEnd(year int,month int,day int) time.Time{
  locationTime :=time.Now()
  dateParse := time.Date(year , time.Month(month),day,23,59,59,999999999, locationTime.Location())
  fmt.Println("time end ",dateParse)
  return dateParse
}
func ParseTime(date string)time.Time {
  timeFormat := "2006-01-02 15:04 MST"
  testdata := date
  splitstr :=strings.Split(testdata, "-")
  month, err := strconv.Atoi(splitstr[1])
  fmt.Println(splitstr)
  if(month<10){
    splitstr[1]= "0"+splitstr[1]
  }
  fmt.Println(splitstr)
  date = splitstr[0];
  for i:=1;i<len(splitstr);i++{
    date =date +"-"+splitstr[i]
  }
  fmt.Println(date)
  then, err := time.Parse(timeFormat, date)
  if(err!=nil){
    fmt.Println(err)
  }
  fmt.Println("shippment time ",then)
  return then
}
func Compare(start time.Time,end time.Time,Mdate time.Time)bool{
  if(Mdate.Before(end)==true && Mdate.After(start)==true){
    return true
  }
  return false
}

func ListenGPS (port string, databaseAddr string, database string, indexServer string ) {
    ln, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("listening error: ", err)
        return
    }
    go gpsAccept(ln,databaseAddr,database,indexServer)
    
}
func gpsAccept(ln net.Listener,databaseAddr string, database string, indexServer string ) {

    i:=0
    nc, err := ln.Accept()
    if err != nil {
        fmt.Println("accepting error: ", err)
        nc.Close()
        //return
    }else{
        fmt.Println("============================================================")
        fmt.Println("=                 Connected with Client                    =")
        fmt.Println("============================================================")
        reciveGps(nc, i, ln,databaseAddr,database,indexServer)
    }

}
func reciveGps(nc net.Conn, i int, ln net.Listener,databaseAddr string, database string, indexServer string ){
    //var msg []byte  
    for{ 
      message, err := bufio.NewReader(nc).ReadString('\n')
      if(err!=nil){
            //fmt.Println("receiving error: ", err)
            fmt.Println("============================================================")
            fmt.Println("=  Client has disconnected. Waiting for new connection...  =")
            fmt.Println("============================================================\n\n")
            nc.Close()
            ln.Close()
            go ListenGPS(":4444",databaseAddr,database,indexServer)

            break
      }else{

        updateLoc(message,i,databaseAddr,database,indexServer)
        i++
      }
    }
}
//status,truckid,time,x,y,
func findOrderByTruckID(truckid string,collection string,database string, databaseAddr string) []Order{
     var result []Order
     fmt.Println(database,"    ",collection,"     ", truckid)
    session, err := mgo.Dial(databaseAddr)
    errhandler(err)
    defer session.Close()
     //var result Order
     c := session.DB(database).C(collection)

     err = c.Find(bson.M{"ordersts.trucks":truckid}).All(&result)
     //if err!=nil {
      fmt.Println("find by truck id err  ",err)
     //}
     return result
}
func findOrderByID(databaseAddr string,database string,collection string,hexID string)Order{
    var result Order
    session, err := mgo.Dial(databaseAddr)
    if err != nil {
    panic(err)
    }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)
  c := session.DB(database).C(collection)
  err = c.FindId(bson.ObjectIdHex(hexID)).One(&result) //Works
  if err != nil {
    fmt.Println(err)
  }
  return result

}


//company,objID, status, truckID,time,x,y
func updateLoc(message string, i int,databaseAddr string , database string , indexServer string ){
  var upOrder Order
  if(message != ""){
    fmt.Println(message)
    GPS:=strings.SplitN(message,",",6)
    z,_ := strconv.Atoi(GPS[1])
    if(z==7){
        x,_ := strconv.ParseFloat(GPS[4],64)
        y,_ := strconv.ParseFloat(strings.Replace(GPS[5],"\n","",1),64)
        old:=findOrderByTruckID(GPS[2],"Carrier",database,databaseAddr)
        fmt.Println("here here here here here here ", len(old))
        for i:=0;i<len(old);i++{
            upOrder=old[i]
            upOrder.OrderSts.TimeStamp= append(old[i].OrderSts.TimeStamp,GPS[3])
            upOrder.OrderSts.GPSCordsX= append(old[i].OrderSts.GPSCordsX,x)
            upOrder.OrderSts.GPSCordsY= append(old[i].OrderSts.GPSCordsY,y)
            upOrder.OrderSts.Status = status[z]
            fmt.Println("before update ", upOrder)
            update(databaseAddr,old[i],upOrder,database,"Carrier")
            addCu := GettingIPAddr(indexServer,"9999",upOrder.CustomerCode,upOrder.CustomerName)
            addother := GettingIPAddr(indexServer,"9999",upOrder.SupplierCode,upOrder.SupplierName)
            sendToOther(addCu, upOrder, "9998")
            sendToOther(addother,upOrder , "9998")

        }
    }else{
            fmt.Println(GPS[0])
     
            singleOrder := findOrderByID(databaseAddr,database,"Carrier",GPS[0])
            fmt.Println(singleOrder.ID.Hex(),"         Carrier")
            if singleOrder.ID.Hex()=="" {
                singleOrder = findOrderByID(databaseAddr,database,"Supplier",GPS[0])
            }
            fmt.Println(singleOrder.ID.Hex(),"          Supplier")
            if singleOrder.ID.Hex()=="" {
              fmt.Println("nothing found here")
                return
            }
            updateOrder := singleOrder
            var addother string
            addother = "0.0.0.0"
            roll := GetRoll(singleOrder,databaseAddr,database)
            	
            /*
				  case mOrder.CustomerCode == result.ID.Hex():
      				return 0
    			case mOrder.CarrierCode == result.ID.Hex():
      				return 1
    			case mOrder.SupplierCode == result.ID.Hex():
					return 2
            */
            if(roll==1 && (z==6 ||z==8)){
               updateOrder.OrderSts.Status = status[z]
               updateOrder.OrderSts.Trucks = GPS[2]
               fmt.Println("before update")
               update(databaseAddr,singleOrder,updateOrder,database,"Carrier")

               addother = GettingIPAddr(indexServer,"9999",updateOrder.SupplierCode,updateOrder.SupplierName)
              // fmt.Println(addother)
            }else if (roll ==2 && z<=7&&z>=0){
               fmt.Println("before update")
                updateOrder.OrderSts.Status = status[z]
                fmt.Println(singleOrder,"     " , updateOrder)
                //func update(databaseAddr string,old Order,new Order,database string,collection string )
                 update(databaseAddr,singleOrder,updateOrder,database,"Supplier")
                 fmt.Println("after update")
                 test := findOrderByID(databaseAddr,database,"Supplier",GPS[0])
                 fmt.Println(test)
                  addother = GettingIPAddr(indexServer,"9999",updateOrder.CarrierCode,updateOrder.CarrierName)
                // fmt.Println(addother)
            }  else  if(roll==-1){
                   fmt.Println("wrong roll")
                   return
                }
            addCu := GettingIPAddr(indexServer,"9999",updateOrder.CustomerCode,updateOrder.CustomerName)
            fmt.Println(addCu)
            sendToOther(addCu, updateOrder, "9998")
            if(addother != "0.0.0.0"){
              fmt.Println("sending..................")
              sendToOther(addother,updateOrder , "9998")
          }
        }
    }
}

func FindInfo(databaseAddr string,database string , indexServerAddr string  ){
     session, err := mgo.Dial(databaseAddr)
     errhandler(err)
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)
    
     d := session.DB(database).C("Info")
     var result indexServer
     err = d.Find(bson.M{"name": database}).One(&result)
     errhandler(err)
     result.ReqType = 1
     sendToIS(indexServerAddr, result, "9999",databaseAddr,database)
}

func sendToIS(ipaddr string,mIndex indexServer,Port string, databaseAddr string, database string ){

  c, err := net.Dial("tcp", ipaddr+":"+Port)
  if err != nil {
        fmt.Println(err)
        return
  }
  b,e := json.Marshal(mIndex)
  if e != nil {
      fmt.Println(e)
      c.Close()
      return
    }
    e2 := json.NewEncoder(c).Encode(b)
    if e2 != nil {
      fmt.Println(e2)
      c.Close()
      return
    }
    recivemsgInitial(c, mIndex,databaseAddr,database)
    c.Close()
}

func recivemsgInitial(c net.Conn, mIndex indexServer, databaseAddr string, database string ){

        var msg []byte
        var ans indexServer
        err := json.NewDecoder(c).Decode(&msg)
        e := json.Unmarshal(msg,&ans)
        if e != nil {
                fmt.Println(e)
            return
        }  
        if err != nil {
                fmt.Println(err)
        } else{
        c.Close()

        infoUpdate(ans, mIndex,databaseAddr,database)
    }
}

func infoUpdate(record indexServer, mIndex indexServer, databaseAddr string, database string ){
    session, err := mgo.Dial(databaseAddr)
    errhandler(err)
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    f := session.DB(database).C("Info")
    err =f.Remove(bson.M{"name":mIndex.Name})
    if err != nil {
        panic(err)
    }
    err = f.Insert(record)
    if err != nil {
        panic(err)
    }
}




