package main
import (
        "fmt"
        "github.com/ant0ine/go-json-rest"
        "net/http"
        "launchpad.net/goamz/aws"
        "launchpad.net/goamz/s3"
)

type Signup struct {
        Id string
        Email string
}

func GetSignup(w *rest.ResponseWriter, req *rest.Request) {
        signup := Signup{
                Id:   req.PathParam("id"),
                Email: "marton@sequenceiq.com",
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.WriteJson(&signup)
}

func PostSignup(w *rest.ResponseWriter, req *rest.Request) {
        signup := Signup{}
	      err := req.DecodeJsonPayload(&signup)
      	if err != nil {
      		rest.Error(w, err.Error(), http.StatusInternalServerError)
      		return
      	}
      	if signup.Email == "" {
      		rest.Error(w, "email is required", 400)
      		return
      	}

        auth, err := aws.EnvAuth()
      	if err != nil {
      		panic(err.Error())
      	}

        s := s3.New(auth, aws.EUWest)
        bucket := s.Bucket("seq-signup")
        data := []byte(signup.Email)
        err = bucket.Put(signup.Email, data, "text/plain", s3.BucketOwnerFull)
      	if err != nil {
      		panic(err.Error())
      	}

        fmt.Println(signup)
        w.Header().Set("Access-Control-Allow-Origin", "*")
	      w.WriteJson(&signup)
}

func OptionsSignup(w *rest.ResponseWriter, req *rest.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
        handler := rest.ResourceHandler{}
        handler.SetRoutes(
                rest.Route{"GET", "/signups/:id", GetSignup},
                rest.Route{"POST", "/signup", PostSignup},
                rest.Route{"OPTIONS", "/signup", OptionsSignup},
        )
        http.ListenAndServe(":8288", &handler)
}
