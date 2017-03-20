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
    "strconv"
    "os"
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
        var i faceRecogInput
        err := json.NewDecoder(r.Body).Decode(&i)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        uid := r.Header.Get("Authorization")
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

        _, err = strconv.Atoi(result[0])
        if err != nil {
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

        resp, err := http.Get("http://localhost/api/v1/users/loved-one?id=" + result[0])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        defer resp.Body.Close()
        _, err = io.Copy(w, resp.Body)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Add("Person-Type", "loved-one")
        w.WriteHeader(http.StatusOK)
    }
}
