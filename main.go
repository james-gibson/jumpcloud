package main

import "fmt"
import "time"
import "crypto/sha512"
import "encoding/base64"
import "net/http"
import "log"
import "strconv"
import "strings"

type hashJob struct {
    id string
    password string
}

type completedHashJob struct {
    id string
    sha string
    timing time.Duration 
}

type statsObject struct {
    total int
    average string 
}

func hashEncode(password string) string {
    hash := sha512.New()
    hash.Write([]byte(password))
    byteSlice := hash.Sum(nil)

    result := base64.StdEncoding.EncodeToString([]byte(byteSlice))
    return result
}

func shutdownResponse(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("503 - Server is shutting down"))
}

func statsHandler(lookupTable map[string]completedHashJob) func (http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, r *http.Request) {
        lookupTableCopy := lookupTable
        hashCount := len(lookupTableCopy)

        sum := int64(0)
        for _, value := range lookupTableCopy {
            sum += int64(value.timing / time.Millisecond)
        }

        average := "0"
        if(hashCount != 0) {
            average = strconv.FormatInt(sum/int64(hashCount),10)
        }
        obj := statsObject{hashCount, average}

        result := strings.Join([]string{"{total:",strconv.Itoa(obj.total),",average:",obj.average,"}"},"")
        fmt.Fprintf(w,result)
    }
}

func hashHandler(addJob func (string) string) func (http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, r *http.Request) {
        fmt.Println("Hash requested")
        r.ParseForm()
        formValues := r.Form
        password := formValues.Get("password")

        jobId := addJob(password)
        fmt.Println("Serving up job id ", jobId)
        if(jobId != "-1") {
            fmt.Fprintf(w, jobId)
            fmt.Println("Hash id ",jobId, " reported to client") 
        } else {
            shutdownResponse(w,r)
        }
    }
}

func shutdownHandler(done chan bool,hashJobs chan hashJob) func (http.ResponseWriter, *http.Request){
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Server shutdown requested")
        close(hashJobs)
        done <- true
    }
}

func shutdownServer() {
    fmt.Println("Server shutting down")

    // Allow time for existing processes to finish
    delay := time.NewTimer(10 * time.Second)
    <-delay.C
    fmt.Println("Goodbye")
}

func addJobToQueue(jobs chan hashJob) func (string) string {
    fmt.Println("Queue initialized")
    i := 0
    return func (password string) string {

        job := hashJob{strconv.Itoa(i),password}
        i++
        id := "-1"
        select {
            case _,more := <-jobs:
                if(!more) {
                   jobs = nil
                }
           default:
        }

        if(jobs != nil) {
            fmt.Println("Adding job ", i, " to queue")
            jobs <- job
            id =  job.id
        } else {
            fmt.Println("Hash request canceled due to shutdown")
        }
        return id
    }
}

func processJobInQueue(id int, jobs chan hashJob, completedJobs chan completedHashJob) {
    go func () {
        for j := range jobs {
            fmt.Println("Worker ", id," picked up job ", j.id)
            password := j.password
            start := time.Now()
            // Delaying to simulate external work taking time
            artificialDelay := time.NewTimer(5 * time.Second)
            <-artificialDelay.C
            sha := hashEncode(password)

            timing := time.Now().Sub(start)
            completedJobs <- completedHashJob{j.id, sha, timing}
        }
    }()
}

func processCompletedJobs(id int, completedJobs chan completedHashJob, lookupTable map[string]completedHashJob) {
   for j := range completedJobs {
        fmt.Println("Cleanup worker ", id, " is adding ", j.id, " to lookup table")
        lookupTable[j.id] = j
   }
}

// leaned on example code from stack overflow to figure this out
func getHashIdFromPath(r *http.Request) string {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) == 1 {
        return "-1"
    } else if (len(parts) > 1) {
        id := parts[2]
        if id != "" {
            return id
        } else {
            return "-1"
        }
    } else {
        return "-1"
    }
}

func lookupId(lookupTable map[string]completedHashJob) func (http.ResponseWriter, *http.Request){
    fmt.Println("Launching lookup service")
    return func (w http.ResponseWriter, r *http.Request) {
        fmt.Println("someone has requested a lookup")
        inputId := getHashIdFromPath(r)
        id := inputId
        j, exists := lookupTable[id]

        if(exists) {
            fmt.Fprintf(w,j.sha)
        } else {
            fmt.Fprintf(w, "error")
        }
    }
}

func main() {
    fmt.Println("Sample golang api ~ James Gibson")
    fmt.Println("Server launching")
    done := make(chan bool,0)
    lookupTable := make(map[string]completedHashJob)
    jobs := make(chan hashJob)
    completedJobs := make(chan completedHashJob)


    //spin up job cleanup
    for w := 1; w <= 1000; w++ {
        go processCompletedJobs(w, completedJobs, lookupTable)
    }
    //spin up workers
    for w := 1; w <= 1000; w++ {
        go processJobInQueue(w, jobs,completedJobs)
    }
    go func(done chan bool) {
        fmt.Println("Spawned server thread")
        http.HandleFunc("/stats", statsHandler(lookupTable))
        http.HandleFunc("/hash/", lookupId(lookupTable))
        http.HandleFunc("/hash", hashHandler(addJobToQueue(jobs)))

        http.HandleFunc("/shutdown", shutdownHandler(done,jobs))

        fmt.Println("Server listening on port 8000")
        err := http.ListenAndServe(":8000", nil)
        if err != nil {
            log.Fatal("ListenAndServe: ", err)
            done <- true
        }
    }(done)

    //Wait for shutdown signal
    <-done
    shutdownServer()
}
