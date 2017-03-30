package web

import (
    "net/http"
    "fa/model"
    "fa/openface"
    "fa/s3util"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
    "os"
    "fmt"
    "crypto/rand"
)

func generateId() (string, error) {
    b := make([]byte, 12)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }

    return base64.URLEncoding.EncodeToString(b), nil
}

func TrainingHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if userId, ok := r.Context().Value("uid").(string); ok {
            var l model.LovedOne
            err := json.NewDecoder(r.Body).Decode(&l)
            if err != nil {
                http.Error(w, "Failed to decode JSON: " + err.Error(), http.StatusBadRequest)
                return
            }

            dir, err := ioutil.TempDir("/tmp/", userId)
            if err != nil {
                http.Error(w, "Failed to make dir", http.StatusInternalServerError)
                return
            }
            imgDir := dir + "/images"
            alignDir := dir + "/align"
            featureDir := dir + "/feature"

            err = os.Mkdir(dir, 0777)
            defer  os.RemoveAll(dir)
            err = os.Mkdir(imgDir, 0777)
            err = os.Mkdir(alignDir, 0777)
            err = os.Mkdir(featureDir, 0777)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            id, err := generateId()
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            l.Prof.Id = id

            err = l.WriteImagesToFile(fmt.Sprintf("%s/%s",imgDir,id))
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            err = openface.Train(imgDir, alignDir, featureDir, userId, r.Header.Get("Authorization"))
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            _, err = l.InsertIntoDB(r.Header.Get("Authorization"))
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            err = s3util.UploadFile(fmt.Sprintf("%s/labels", featureDir), fmt.Sprintf("features/%s/%s/labels.csv", userId, id))
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            err = s3util.UploadFile(fmt.Sprintf("%s/reps", featureDir), fmt.Sprintf("features/%s/%s/reps.csv", userId, id))
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            err = s3util.UploadFile(fmt.Sprintf("%s/classifier.pkl", featureDir), fmt.Sprintf("features/%s/classifier.pkl", userId))
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        } else {
            http.Error(w, "Error getting ID", http.StatusInternalServerError)
            return
        }
    }
}
