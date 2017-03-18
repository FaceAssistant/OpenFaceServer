package web

import (
    "net/http"
    "fa/model"
    "fa/openface"
    "fa/s3util"
    "encoding/json"
    "os"
    "fmt"
)

func TrainingHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var l model.LovedOne
        err := json.NewDecoder(r.Body).Decode(&l)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        dir := fmt.Sprintf("/tmp/%d", l.UserId)
        err = os.Mkdir(dir, os.ModeTemporary)
        defer  os.RemoveAll(dir)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        imgDir := dir + "images"
        alignDir := dir + "align"
        featureDir := dir + "feature"

        err = l.WriteImagesToFile(imgDir)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = openface.Train(imgDir, alignDir, featureDir, l.UserId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        id, err := l.InsertIntoDB()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = s3util.UploadFile(fmt.Sprintf("%s/labels.csv", featureDir), fmt.Sprintf("%s/%s/labels.csv", l.UserId, id))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = s3util.UploadFile(fmt.Sprintf("%s/reps.csv", featureDir), fmt.Sprintf("%s/%s/reps.csv", l.UserId, id))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
}
