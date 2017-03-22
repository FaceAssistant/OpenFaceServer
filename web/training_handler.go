package web

import (
    "net/http"
    "fa/model"
    "fa/openface"
    "fa/s3util"
    "io/ioutil"
    "encoding/json"
    "os"
    "fmt"
    "strconv"
)

func TrainingHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var l model.LovedOne
        err := json.NewDecoder(r.Body).Decode(&l)
        if err != nil {
            http.Error(w, "Failed to decode JSON: " + err.Error(), http.StatusBadRequest)
            return
        }

        uid := r.Header.Get("Authorization")
        userId, err := strconv.Atoi(uid)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        //dir := fmt.Sprintf("/tmp/%d", userId)
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

        err = l.WriteImagesToFile(imgDir + "/" + uid)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = openface.Train(imgDir, alignDir, featureDir, userId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        id, err := l.InsertIntoDB(uid)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = s3util.UploadFile(fmt.Sprintf("%s/labels", featureDir), fmt.Sprintf("features/%d/%d/labels.csv", userId, id))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = s3util.UploadFile(fmt.Sprintf("%s/reps", featureDir), fmt.Sprintf("features/%d/%d/reps.csv", userId, id))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = s3util.UploadFile(fmt.Sprintf("%s/classifier.pkl", featureDir), fmt.Sprintf("features/%d/classifier.pkl", userId))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
}
