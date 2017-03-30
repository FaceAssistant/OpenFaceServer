package web

import (
    "net/http"
    "encoding/json"
    "encoding/base64"
    "io"
    "io/ioutil"
    "fa/s3util"
    "fa/model"
    "fa/openface"
    "os"
    "fmt"
)

type faceRecogInput struct {
    Image string `json:"image"`
}

type celebOutput struct {
    Name string `json:"name"`
    Confidence string `json:"confidence"`
}

func FaceRecogHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if uid, ok := r.Context().Value("uid").(string); ok {
            var i faceRecogInput
            err := json.NewDecoder(r.Body).Decode(&i)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            classifier, err := ioutil.TempFile("/tmp/", "classifier")
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            defer classifier.Close()
            defer os.Remove(classifier.Name())

            err = s3util.GetClassifier(uid,classifier)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            data, err := base64.StdEncoding.DecodeString(i.Image)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            img, err := model.WriteBytesToFile(data, "/tmp/")
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            defer os.Remove(img)

            result, err := openface.Infer(classifier.Name(), img)
            if err != nil  {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            resp, err := model.GetLovedOneById(result[0], r.Header.Get("Authorization"))
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if resp.StatusCode != http.StatusOK {
                fmt.Println("get celeb")
                o := &celebOutput{
                    Name: result[0],
                    Confidence: result[1],
                }
                w.Header().Add("Person-Type", "celeb")
                w.Header().Set("Content-Type", "applicaton/json")
                err = json.NewEncoder(w).Encode(&o)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                return
            }
            defer resp.Body.Close()

            w.Header().Add("Person-Type", "loved-one")
            _, err = io.Copy(w, resp.Body)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            fmt.Println("Copying")

        } else {
            http.Error(w, "Error getting user id", http.StatusInternalServerError)
            return
        }
    }
}
